package merger

import (
	"log"
	"sort"
	"strings"

	"clash-rules-cn/internal/types"
)

// Merge 合并多个源的条目，按 behavior 过滤并去重排序
func Merge(entries []types.Entry, behavior string) []types.Entry {
	// 按 behavior 过滤
	filtered := filterByBehavior(entries, behavior)

	// 去重
	deduped := deduplicate(filtered)

	// 排序
	sortEntries(deduped)

	return deduped
}

// filterByBehavior 根据 behavior 类型过滤条目
func filterByBehavior(entries []types.Entry, behavior string) []types.Entry {
	// 预分配容量
	result := make([]types.Entry, 0, len(entries))

	for _, entry := range entries {
		switch behavior {
		case "domain":
			// domain behavior 只接受 DOMAIN, DOMAIN-SUFFIX, DOMAIN-KEYWORD
			if entry.Type == "DOMAIN" || entry.Type == "DOMAIN-SUFFIX" || entry.Type == "DOMAIN-KEYWORD" {
				result = append(result, entry)
			}
		case "ipcidr":
			// ipcidr behavior 只接受 IP-CIDR, IP-CIDR6
			if entry.Type == "IP-CIDR" || entry.Type == "IP-CIDR6" {
				result = append(result, entry)
			}
		case "classical":
			// classical behavior 接受所有类型
			result = append(result, entry)
		}
	}

	return result
}

// deduplicate 去除重复条目
func deduplicate(entries []types.Entry) []types.Entry {
	seen := make(map[string]struct{}, len(entries))
	result := make([]types.Entry, 0, len(entries))

	for _, entry := range entries {
		// 使用 Type+Payload 作为唯一标识
		key := entry.Type + ":" + entry.Payload
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, entry)
		}
	}

	return result
}

// sortEntries 对条目进行排序
func sortEntries(entries []types.Entry) {
	sort.Slice(entries, func(i, j int) bool {
		// 先按 Type 排序
		if entries[i].Type != entries[j].Type {
			return entries[i].Type < entries[j].Type
		}
		// 再按 Payload 排序（不区分大小写）
		return strings.ToLower(entries[i].Payload) < strings.ToLower(entries[j].Payload)
	})
}

// MergeLoon 将 applications 上游被旧解析器包装的进程规则恢复为原生类型。
// 该转换只用于 Loon，Clash 沿用现有解析及输出以保持兼容。
func MergeLoon(entries []types.Entry, behavior string) []types.Entry {
	var normalized, processes []types.Entry
	for _, entry := range entries {
		if entry.Type == "DOMAIN" || entry.Type == "DOMAIN-SUFFIX" || entry.Type == "DOMAIN-KEYWORD" {
			// Surge extended-matching 没有 Loon 对应参数，保留域名匹配并显式报告降级。
			if strings.HasSuffix(entry.Payload, ",extended-matching") {
				log.Printf("Loon 格式转换: %s,%s 移除 Surge 专用 extended-matching（不承诺 CNAME 扩展匹配）", entry.Type, entry.Payload)
				entry.Payload = strings.TrimSuffix(entry.Payload, ",extended-matching")
			}
			// v2fly 原始数据可能在域名后保留行内注释，不能写进 Loon 匹配值。
			payload, _, _ := strings.Cut(entry.Payload, " #")
			entry.Payload = strings.TrimSpace(payload)
			if entry.Payload == "http:" || entry.Payload == "https:" {
				log.Printf("Loon 格式转换: 排除 GFWList 旧解析产生的协议伪域名 %q", entry.Payload)
				continue
			}
		}
		if entry.Type == "DOMAIN-SUFFIX" && strings.HasPrefix(entry.Payload, "PROCESS-NAME,") {
			entry = types.Entry{Type: "PROCESS-NAME", Payload: strings.TrimPrefix(entry.Payload, "PROCESS-NAME,")}
		}
		if behavior == "domain" && entry.Type == "PROCESS-NAME" {
			processes = append(processes, entry)
		} else {
			normalized = append(normalized, entry)
		}
	}
	filtered := filterByBehavior(normalized, behavior)
	return Merge(append(filtered, processes...), "classical")
}

// ApplyExclude 剔除上游误分类的条目。
// exclude 采用 DOMAIN-SUFFIX 语义：条目与 payload 相等，或 payload 以 ".条目" 结尾。
// 条目允许写成 ".example.com" 或 "example.com"，实现时统一去掉前导点。
// 只做减法，不改分类逻辑；exclude 为空时原样返回，不产生额外分配。
func ApplyExclude(entries []types.Entry, exclude []string) []types.Entry {
	patterns := normalizePatterns(exclude)
	if len(patterns) == 0 {
		return entries
	}

	result := make([]types.Entry, 0, len(entries))
	removed := 0
	for _, entry := range entries {
		if matchExclude(entry, patterns) {
			removed++
			continue
		}
		result = append(result, entry)
	}
	if removed > 0 {
		log.Printf("exclude: 剔除 %d 条上游误分类条目（清单 %d 项）", removed, len(patterns))
	}
	return result
}

// normalizePatterns 归一化排除清单：去空白、去前导点、转小写，丢弃空条目
func normalizePatterns(exclude []string) []string {
	patterns := make([]string, 0, len(exclude))
	for _, raw := range exclude {
		p := strings.ToLower(strings.TrimSpace(raw))
		p = strings.TrimPrefix(p, ".")
		if p != "" {
			patterns = append(patterns, p)
		}
	}
	return patterns
}

// matchExclude 判断单条条目是否命中排除清单。
// 只对域名类条目做匹配，IP-CIDR / PROCESS-NAME 等类型保留原样，
// 避免把 CIDR 或进程名误当域名比对。
func matchExclude(entry types.Entry, patterns []string) bool {
	switch entry.Type {
	case "DOMAIN", "DOMAIN-SUFFIX", "DOMAIN-KEYWORD":
	default:
		return false
	}

	payload := strings.ToLower(strings.TrimSpace(entry.Payload))
	// classical 源可能在 payload 里带参数（如 no-resolve），只取逗号前的匹配值
	if i := strings.IndexByte(payload, ','); i >= 0 {
		payload = payload[:i]
	}
	payload = strings.TrimPrefix(payload, ".")

	for _, p := range patterns {
		if payload == p || strings.HasSuffix(payload, "."+p) {
			return true
		}
	}
	return false
}
