package parser

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"clash-rules-cn/internal/types"
)

// Parser 定义解析器接口
type Parser interface {
	Parse(data []byte) ([]types.Entry, error)
}

// GetParser 根据类型返回对应的解析器
func GetParser(parserType string) Parser {
	switch parserType {
	case "cidr":
		return &CIDRParser{}
	case "dnsmasq":
		return &DnsmasqParser{}
	case "domain-list":
		return &DomainListParser{}
	case "gfwlist":
		return &GFWListParser{}
	case "v2fly-dlc":
		return NewV2flyDLCParser()
	case "classical":
		return &ClassicalParser{}
	case "loyalsoldier":
		return &LoyalsoldierParser{}
	case "static":
		return &StaticParser{}
	default:
		return nil
	}
}

// CIDRParser 解析 CIDR 格式
type CIDRParser struct{}

func (p *CIDRParser) Parse(data []byte) ([]types.Entry, error) {
	var entries []types.Entry
	scanner := scannerFromBytes(data)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 验证是否为有效的 CIDR
		_, ipNet, err := net.ParseCIDR(line)
		if err != nil {
			continue
		}

		// 判断是 IPv4 还是 IPv6
		entryType := "IP-CIDR"
		if ipNet.IP.To4() == nil {
			entryType = "IP-CIDR6"
		}

		entries = append(entries, types.Entry{
			Type:    entryType,
			Payload: line,
		})
	}

	return entries, nil
}

// DnsmasqParser 解析 dnsmasq 格式
type DnsmasqParser struct{}

func (p *DnsmasqParser) Parse(data []byte) ([]types.Entry, error) {
	var entries []types.Entry
	scanner := scannerFromBytes(data)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 格式: server=/domain/ip
		if !strings.HasPrefix(line, "server=/") {
			continue
		}

		// 提取域名部分
		parts := strings.Split(line, "/")
		if len(parts) < 3 {
			continue
		}

		domain := parts[1]
		if domain == "" {
			continue
		}

		entries = append(entries, types.Entry{
			Type:    "DOMAIN-SUFFIX",
			Payload: domain,
		})
	}

	return entries, nil
}

// DomainListParser 解析域名列表格式
type DomainListParser struct{}

func (p *DomainListParser) Parse(data []byte) ([]types.Entry, error) {
	var entries []types.Entry
	scanner := scannerFromBytes(data)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		entries = append(entries, types.Entry{
			Type:    "DOMAIN-SUFFIX",
			Payload: line,
		})
	}

	return entries, nil
}

// GFWListParser 解析 GFWList（Base64 编码的 ABP 规则）
type GFWListParser struct{}

func (p *GFWListParser) Parse(data []byte) ([]types.Entry, error) {
	// Base64 解码
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	var entries []types.Entry
	scanner := scannerFromBytes(decoded)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "!") || strings.HasPrefix(line, "[") {
			continue
		}

		// 简化的 ABP 解析，提取域名
		domain := extractDomainFromABP(line)
		if domain != "" {
			entries = append(entries, types.Entry{
				Type:    "DOMAIN-SUFFIX",
				Payload: domain,
			})
		}
	}

	return entries, nil
}

// V2flyDLCParser 解析 v2fly domain-list-community 格式
type V2flyDLCParser struct {
	client     *http.Client
	baseURL    string
	visited    map[string]bool
	visitedMu  sync.Mutex
	cache      map[string][]byte
	cacheMu    sync.RWMutex
}

// 全局缓存，所有 V2flyDLCParser 实例共享
var globalV2flyCache = struct {
	data map[string][]byte
	mu   sync.RWMutex
}{
	data: make(map[string][]byte),
}

// NewV2flyDLCParser 创建一个支持递归解析的 V2flyDLCParser
func NewV2flyDLCParser() *V2flyDLCParser {
	return &V2flyDLCParser{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/",
		visited: make(map[string]bool),
		cache:   make(map[string][]byte),
	}
}

func (p *V2flyDLCParser) Parse(data []byte) ([]types.Entry, error) {
	return p.parseRecursive(data)
}

// parseRecursive 递归解析（支持 include）
func (p *V2flyDLCParser) parseRecursive(data []byte) ([]types.Entry, error) {
	var entries []types.Entry
	var includeNames []string

	// 第一遍：收集所有 include 指令和直接条目
	scanner := scannerFromBytes(data)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 处理 include 指令
		if strings.HasPrefix(line, "include:") {
			includeName := strings.TrimPrefix(line, "include:")
			includeName = strings.TrimSpace(includeName)

			// 移除可能的 @ 后缀（如 @cn）
			if idx := strings.Index(includeName, "@"); idx != -1 {
				includeName = includeName[:idx]
			}

			// 检查是否已访问过（防止循环引用）
			p.visitedMu.Lock()
			if !p.visited[includeName] {
				p.visited[includeName] = true
				includeNames = append(includeNames, includeName)
			}
			p.visitedMu.Unlock()
			continue
		}

		entry := p.parseLine(line)
		if entry != nil {
			entries = append(entries, *entry)
		}
	}

	// 第二遍：并发下载并解析所有 include 文件
	if len(includeNames) > 0 {
		type includeResult struct {
			entries []types.Entry
		}

		results := make([]includeResult, len(includeNames))
		var wg sync.WaitGroup

		for i, name := range includeNames {
			wg.Add(1)
			go func(idx int, includeName string) {
				defer wg.Done()
				includeEntries, err := p.fetchAndParse(includeName)
				if err == nil {
					results[idx] = includeResult{entries: includeEntries}
				}
			}(i, name)
		}

		wg.Wait()

		// 合并结果
		for _, result := range results {
			entries = append(entries, result.entries...)
		}
	}

	return entries, nil
}

// parseLine 解析单行内容
func (p *V2flyDLCParser) parseLine(line string) *types.Entry {
	// 跳过 include 指令
	if strings.HasPrefix(line, "include:") {
		return nil
	}

	// 跳过 regexp: 指令
	if strings.HasPrefix(line, "regexp:") {
		return nil
	}

	// 格式: domain:example.com 或 full:example.com 或 keyword:example
	// 也支持直接的域名（如 google.com）
	if strings.Contains(line, ":") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil
		}

		entryType := parts[0]
		payload := strings.TrimSpace(parts[1])

		// 移除可能的 @cn 等后缀
		if idx := strings.Index(payload, " @"); idx != -1 {
			payload = payload[:idx]
		}

		switch entryType {
		case "domain":
			return &types.Entry{
				Type:    "DOMAIN-SUFFIX",
				Payload: payload,
			}
		case "full":
			return &types.Entry{
				Type:    "DOMAIN",
				Payload: payload,
			}
		case "keyword":
			return &types.Entry{
				Type:    "DOMAIN-KEYWORD",
				Payload: payload,
			}
		}
	} else {
		// 直接的域名（如 google.com）
		payload := strings.TrimSpace(line)

		// 移除可能的 @cn 等后缀
		if idx := strings.Index(payload, " @"); idx != -1 {
			payload = payload[:idx]
		}

		if payload != "" {
			return &types.Entry{
				Type:    "DOMAIN-SUFFIX",
				Payload: payload,
			}
		}
	}

	return nil
}

// fetchAndParse 下载并解析指定的文件（带缓存）
func (p *V2flyDLCParser) fetchAndParse(name string) ([]types.Entry, error) {
	// 先检查全局缓存
	globalV2flyCache.mu.RLock()
	if cached, ok := globalV2flyCache.data[name]; ok {
		globalV2flyCache.mu.RUnlock()
		return p.parseRecursive(cached)
	}
	globalV2flyCache.mu.RUnlock()

	// 下载文件
	url := p.baseURL + name
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("下载 %s 失败: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载 %s 失败: HTTP 状态码 %d", name, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 %s 失败: %w", name, err)
	}

	// 存入全局缓存
	globalV2flyCache.mu.Lock()
	globalV2flyCache.data[name] = body
	globalV2flyCache.mu.Unlock()

	// 递归解析
	return p.parseRecursive(body)
}

// ClassicalParser 解析 classical 格式（TYPE,payload）
type ClassicalParser struct{}

func (p *ClassicalParser) Parse(data []byte) ([]types.Entry, error) {
	var entries []types.Entry
	scanner := scannerFromBytes(data)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			continue
		}

		entryType := strings.TrimSpace(parts[0])
		payload := strings.TrimSpace(parts[1])

		entries = append(entries, types.Entry{
			Type:    entryType,
			Payload: payload,
		})
	}

	return entries, nil
}

// LoyalsoldierParser 解析 Loyalsoldier 格式（YAML payload 列表）
// 自动检测条目类型：CIDR → IP-CIDR/IP-CIDR6，其余 → DOMAIN-SUFFIX
type LoyalsoldierParser struct{}

func (p *LoyalsoldierParser) Parse(data []byte) ([]types.Entry, error) {
	var entries []types.Entry
	scanner := scannerFromBytes(data)
	inPayload := false

	for scanner.Scan() {
		line := scanner.Text()

		// 检测 payload: 行
		if strings.TrimSpace(line) == "payload:" {
			inPayload = true
			continue
		}

		if !inPayload {
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 移除 YAML 列表前缀 "- "
		if strings.HasPrefix(line, "- ") {
			line = strings.TrimPrefix(line, "- ")
		}

		// 移除引号
		if (strings.HasPrefix(line, "'") && strings.HasSuffix(line, "'")) ||
			(strings.HasPrefix(line, "\"") && strings.HasSuffix(line, "\"")) {
			line = line[1 : len(line)-1]
		}

		if line == "" {
			continue
		}

		// 自动检测类型
		if _, _, err := net.ParseCIDR(line); err == nil {
			entryType := "IP-CIDR"
			if ip := net.ParseIP(line); ip != nil && ip.To4() == nil {
				entryType = "IP-CIDR6"
			}
			entries = append(entries, types.Entry{
				Type:    entryType,
				Payload: line,
			})
		} else {
			// 处理 +.domain 格式（如 +.gov → DOMAIN-SUFFIX: gov）
			payload := strings.TrimPrefix(line, "+.")
			entries = append(entries, types.Entry{
				Type:    "DOMAIN-SUFFIX",
				Payload: payload,
			})
		}
	}

	return entries, nil
}

// StaticParser 解析静态配置的条目
type StaticParser struct{}

func (p *StaticParser) Parse(data []byte) ([]types.Entry, error) {
	// Static 解析器不使用 data，而是直接从配置中读取
	// 这个方法实际上不会被调用，StaticParser 的使用方式不同
	return nil, nil
}

// ParseStatic 解析静态配置的条目列表
func ParseStatic(entries []string) []types.Entry {
	var result []types.Entry

	for _, entry := range entries {
		// 尝试解析为 CIDR
		_, _, err := net.ParseCIDR(entry)
		if err == nil {
			entryType := "IP-CIDR"
			// 判断是否为 IPv6
			ip := net.ParseIP(entry)
			if ip != nil && ip.To4() == nil {
				entryType = "IP-CIDR6"
			}
			result = append(result, types.Entry{
				Type:    entryType,
				Payload: entry,
			})
			continue
		}

		// 否则作为域名处理
		result = append(result, types.Entry{
			Type:    "DOMAIN-SUFFIX",
			Payload: entry,
		})
	}

	return result
}

// scannerFromBytes 创建一个 bufio.Scanner
func scannerFromBytes(data []byte) *bufio.Scanner {
	return bufio.NewScanner(strings.NewReader(string(data)))
}

// extractDomainFromABP 从 ABP 规则中提取域名
func extractDomainFromABP(rule string) string {
	// 移除 ABP 特殊规则
	if strings.HasPrefix(rule, "@@") {
		return ""
	}

	// 移除选项部分
	if idx := strings.Index(rule, "$"); idx != -1 {
		rule = rule[:idx]
	}

	// 移除通配符
	rule = strings.ReplaceAll(rule, "*", "")

	// 移除协议前缀
	rule = strings.TrimPrefix(rule, "||")
	rule = strings.TrimPrefix(rule, "|")

	// 移除路径
	if idx := strings.Index(rule, "/"); idx != -1 {
		rule = rule[:idx]
	}

	// 验证是否为有效域名
	if rule == "" || strings.Contains(rule, " ") {
		return ""
	}

	return rule
}
