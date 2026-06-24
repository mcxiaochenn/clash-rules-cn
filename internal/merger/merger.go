package merger

import (
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
