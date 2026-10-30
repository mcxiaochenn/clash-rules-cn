package loon

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestOnlyResourceURLsChange(t *testing.T) {
	if got := fmt.Sprintf("%x", sha256.Sum256([]byte(Upstream))); got != "3a07e4dc097243a437679dffd2b0fae6e6382038688624cfddb284c8d1580629" {
		t.Fatalf("上游快照变动: %s", got)
	}
	generated, err := Generate()
	if err != nil {
		t.Fatal(err)
	}
	// 固定反向映射验证：去掉 CDN 地址变化后仍只能有原定六项分流和两项数据库变化。
	canonical := strings.NewReplacer(
		"https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/", "https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/",
		"https://cdn.jsdelivr.net/gh/Orz-3/mini@master/", "https://raw.githubusercontent.com/Orz-3/mini/master/",
		"https://cdn.jsdelivr.net/gh/sub-store-org/Sub-Store@release/", "https://raw.githubusercontent.com/sub-store-org/Sub-Store/release/",
		"https://cdn.jsdelivr.net/gh/shindgewongxj/WHATSINStash@main/icon/fallback.png", "https://github.com/shindgewongxj/WHATSINStash/raw/main/icon/fallback.png",
	).Replace(generated)
	before, after := strings.SplitAfter(Upstream, "\n"), strings.SplitAfter(canonical, "\n")
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
		if oldOptions != newOptions || newURL != rawBaseURL+replacements[oldURL] {
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

func TestCDNGitHubURLs(t *testing.T) {
	cases := map[string]string{
		"https://raw.githubusercontent.com/Orz-3/mini/master/Color/HK.png":                               "https://cdn.jsdelivr.net/gh/Orz-3/mini@master/Color/HK.png",
		"https://github.com/shindgewongxj/WHATSINStash/raw/main/icon/fallback.png":                       "https://cdn.jsdelivr.net/gh/shindgewongxj/WHATSINStash@main/icon/fallback.png",
		"https://raw.githubusercontent.com/sub-store-org/Sub-Store/release/sub-store-parser.loon.min.js": "https://cdn.jsdelivr.net/gh/sub-store-org/Sub-Store@release/sub-store-parser.loon.min.js",
		"https://cdn.jsdelivr.net/gh/Orz-3/mini@master/Color/HK.png":                                     "https://cdn.jsdelivr.net/gh/Orz-3/mini@master/Color/HK.png",
		"https://kelee.one/Tool/Loon/Lpx/BoxJs.lpx":                                                      "https://kelee.one/Tool/Loon/Lpx/BoxJs.lpx",
	}
	for input, want := range cases {
		if got := cdnGitHubURLs(input); got != want {
			t.Fatalf("%s 转换错误: %s", input, got)
		}
	}
}

func TestEveryGitHubResourceUsesJSDelivr(t *testing.T) {
	generated, err := Generate()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(generated, "gh-proxy") || strings.Contains(generated, "ghproxy") {
		t.Fatal("残留 ghproxy 地址")
	}
	unique := make(map[string]bool)
	for _, resource := range regexp.MustCompile(`https://[^\s,]+`).FindAllString(generated, -1) {
		parsed, err := url.Parse(resource)
		if err != nil {
			t.Fatal(err)
		}
		if parsed.Hostname() == "github.com" || parsed.Hostname() == "raw.githubusercontent.com" {
			t.Fatalf("残留 GitHub 直连资源: %s", resource)
		}
		if parsed.Hostname() == "cdn.jsdelivr.net" {
			if !strings.HasPrefix(parsed.Path, "/gh/") || !strings.Contains(parsed.Path, "@") {
				t.Fatalf("无效 jsDelivr 地址: %s", resource)
			}
			unique[resource] = true
		}
	}
	if len(unique) != 18 {
		t.Fatalf("CDN 资源数量不符: %d", len(unique))
	}
}
