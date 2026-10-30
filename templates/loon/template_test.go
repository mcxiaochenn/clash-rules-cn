package loon

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

func TestOnlySixRulesAndTwoDatabaseURLsChange(t *testing.T) {
	if got := fmt.Sprintf("%x", sha256.Sum256([]byte(Upstream))); got != "3a07e4dc097243a437679dffd2b0fae6e6382038688624cfddb284c8d1580629" {
		t.Fatalf("上游快照变动: %s", got)
	}
	generated, err := Generate()
	if err != nil {
		t.Fatal(err)
	}
	before, after := strings.SplitAfter(Upstream, "\n"), strings.SplitAfter(generated, "\n")
	if len(before) != len(after) {
		t.Fatal("行数变化")
	}
	section, count, databaseCount := "", 0, 0
	for i, line := range before {
		if strings.HasPrefix(strings.TrimSpace(line), "[") {
			section = strings.TrimSpace(line)
		}
		if line == after[i] {
			continue
		}
		if section == "[General]" {
			expected := map[string]string{
				"geoip-url": "https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/Country.mmdb",
				"ipasn-url": "https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/GeoLite2-ASN.mmdb",
			}
			key, value, ok := strings.Cut(strings.TrimSpace(after[i]), "=")
			oldKey, _, _ := strings.Cut(strings.TrimSpace(line), "=")
			if !ok || key != oldKey || expected[key] == "" || expected[key] != value {
				t.Fatalf("数据库地址替换错误: %s", after[i])
			}
			databaseCount++
			continue
		}
		if section != "[Remote Rule]" {
			t.Fatalf("修改了其他配置: %s", section)
		}
		oldURL, oldOptions, _ := strings.Cut(line, ",")
		newURL, newOptions, _ := strings.Cut(after[i], ",")
		if oldOptions != newOptions || newURL != BaseURL+replacements[oldURL] {
			t.Fatalf("不允许的修改: %s", after[i])
		}
		count++
	}
	if count != 6 || databaseCount != 2 {
		t.Fatalf("修改分流 %d 项、数据库 %d 项", count, databaseCount)
	}
	// 检查六项订阅和 FINAL 的策略都在原配置中定义。
	groups := map[string]bool{"DIRECT": true}
	section = ""
	for _, line := range strings.Split(generated, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") {
			section = line
		}
		if section == "[Proxy Group]" && !strings.HasPrefix(line, "#") {
			if name, _, ok := strings.Cut(line, "="); ok {
				groups[strings.TrimSpace(name)] = true
			}
		}
	}
	section = ""
	for _, line := range strings.Split(generated, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") {
			section = line
		}
		if section == "[Remote Rule]" && strings.HasPrefix(line, "https://") {
			for _, field := range strings.Split(line, ",") {
				if name, ok := strings.CutPrefix(strings.TrimSpace(field), "policy="); ok && !groups[name] {
					t.Fatalf("无效策略: %s", name)
				}
			}
		}
		if strings.HasPrefix(line, "FINAL,") && !groups[strings.TrimSpace(strings.TrimPrefix(line, "FINAL,"))] {
			t.Fatal("无效 FINAL 策略")
		}
	}
}
