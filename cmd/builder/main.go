package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"clash-rules-cn/internal/config"
	"clash-rules-cn/internal/fetcher"
	"clash-rules-cn/internal/merger"
	"clash-rules-cn/internal/parser"
	"clash-rules-cn/internal/types"
	"clash-rules-cn/internal/writer"
)

func main() {
	log.Println("开始构建 Clash 规则...")

	// 加载配置
	cfg, err := config.Load("sources.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建 fetcher
	f := fetcher.New(3)

	// 收集需要拉取的 URL
	urls := make(map[string]string)
	for name, source := range cfg.Upstream {
		if source.Parser != "static" {
			urls[name] = source.URL
		}
	}

	// 并发拉取所有数据
	log.Printf("正在拉取 %d 个数据源...", len(urls))
	results := f.FetchAll(urls)

	// 收集拉取结果
	fetchResults := make(map[string][]byte)
	for result := range results {
		if result.Error != nil {
			log.Printf("拉取 %s 失败: %v", result.Name, result.Error)
			continue
		}
		fetchResults[result.Name] = result.Content
		log.Printf("拉取 %s 成功", result.Name)
	}

	// 解析所有数据
	log.Println("正在解析数据...")
	parsedData := make(map[string][]types.Entry)

	for name, source := range cfg.Upstream {
		var entries []types.Entry
		var err error

		if source.Parser == "static" {
			// 静态数据直接解析
			entries = parser.ParseStatic(source.Entries)
		} else {
			// 从拉取结果中获取数据
			data, exists := fetchResults[name]
			if !exists {
				log.Printf("跳过 %s: 无数据", name)
				continue
			}

			// 获取对应的解析器
			p := parser.GetParser(source.Parser)
			if p == nil {
				log.Printf("跳过 %s: 未知的解析器类型 %s", name, source.Parser)
				continue
			}

			entries, err = p.Parse(data)
			if err != nil {
				log.Printf("解析 %s 失败: %v", name, err)
				continue
			}
		}

		parsedData[name] = entries
		log.Printf("解析 %s 完成: %d 条规则", name, len(entries))
	}

	// 创建 writer
	w := writer.New("rules")

	// 处理每个输出规则
	log.Println("正在生成规则文件...")
	for ruleName, rule := range cfg.Output {
		// 收集该规则的所有条目
		var allEntries []types.Entry
		for _, sourceName := range rule.Sources {
			if entries, exists := parsedData[sourceName]; exists {
				allEntries = append(allEntries, entries...)
			}
		}

		// 合并去重排序
		merged := merger.Merge(allEntries, rule.Behavior)

		// 写入文件
		if err := w.Write(rule.Filename, rule.Description, rule.Behavior, merged); err != nil {
			log.Printf("写入 %s 失败: %v", ruleName, err)
			continue
		}

		log.Printf("生成 %s 完成: %d 条规则", rule.Filename, len(merged))
	}

	// 下载 GeoIP 文件
	log.Println("正在下载 GeoIP 文件...")
	if err := downloadGeoIP(cfg.GeoIP); err != nil {
		log.Printf("下载 GeoIP 文件失败: %v", err)
	}

	log.Println("构建完成!")
}

// downloadGeoIP 下载 GeoIP 文件
func downloadGeoIP(geoipConfig config.GeoIPConfig) error {
	// 确保 geoip 目录存在
	if err := os.MkdirAll("geoip", 0755); err != nil {
		return fmt.Errorf("创建 geoip 目录失败: %w", err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	for _, file := range geoipConfig.Files {
		filePath := filepath.Join("geoip", file.Filename)

		// 检查文件是否已存在
		if _, err := os.Stat(filePath); err == nil {
			log.Printf("GeoIP 文件 %s 已存在，跳过下载", file.Filename)
			continue
		}

		log.Printf("正在下载 %s...", file.Filename)

		// 下载文件
		resp, err := client.Get(file.URL)
		if err != nil {
			return fmt.Errorf("下载 %s 失败: %w", file.Filename, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("下载 %s 失败: HTTP 状态码 %d", file.Filename, resp.StatusCode)
		}

		// 写入文件
		outFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("创建文件 %s 失败: %w", file.Filename, err)
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, resp.Body); err != nil {
			return fmt.Errorf("写入文件 %s 失败: %w", file.Filename, err)
		}

		log.Printf("下载 %s 完成", file.Filename)
	}

	return nil
}
