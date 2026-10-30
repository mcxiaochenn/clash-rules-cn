// Package loon 保存固定上游配置，生成分流、数据库及 GitHub 资源的 jsDelivr 配置。
package loon

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
)

//go:embed upstream.lcf
var Upstream string

//go:embed LICENSE
var License []byte

const ConfigFilename = "Loon_Auto_Select_Config.lcf"
const rawBaseURL = "https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/loon/"
const BaseURL = "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/loon/"

var rawGitHubURL = regexp.MustCompile(`https://raw\.githubusercontent\.com/([^/\s]+)/([^/\s]+)/([^/\s]+)/([^\s,]+)`)
var githubRawURL = regexp.MustCompile(`https://github\.com/([^/\s]+)/([^/\s]+)/raw/([^/\s]+)/([^\s,]+)`)

// cdnGitHubURLs 只转换配置中使用的两种 GitHub 文件地址，保留分支与文件路径。
func cdnGitHubURLs(text string) string {
	text = rawGitHubURL.ReplaceAllString(text, "https://cdn.jsdelivr.net/gh/$1/$2@$3/$4")
	return githubRawURL.ReplaceAllString(text, "https://cdn.jsdelivr.net/gh/$1/$2@$3/$4")
}

var replacements = map[string]string{
	"https://rule.kelee.one/Loon/Telegram.lsr":                   "telegram.lsr",
	"https://kelee.one/Tool/Loon/Lsr/TikTok.lsr":                 "tiktok.lsr",
	"https://kelee.one/Tool/Loon/Lsr/AI.lsr":                     "ai-domain.lsr",
	"https://kelee.one/Tool/Loon/Lsr/SpeedtestInternational.lsr": "speedtest-international.lsr",
	"https://kelee.one/Tool/Loon/Lsr/LAN_SPLITTER.lsr":           "lan.lsr",
	"https://kelee.one/Tool/Loon/Lsr/REGION_SPLITTER.lsr":        "cn-region.lsr",
}

// 数据库使用 rules 分支根目录的真实国家库和 IP→ASN 库。
var databaseReplacements = map[string]string{
	"geoip-url=https://geodata.kelee.one/Country-Masaiki.mmdb":     "geoip-url=https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/Country.mmdb",
	"ipasn-url=https://geodata.kelee.one/GeoLite2-ASN-P3TERX.mmdb": "ipasn-url=https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/GeoLite2-ASN.mmdb",
}

// Generate 按行处理并保留换行符，模板结构变动时拒绝生成。
func Generate() (string, error) {
	lines := strings.SplitAfter(Upstream, "\n")
	section, count := "", 0
	seen := make(map[string]bool)
	databaseCount := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") {
			section = trimmed
		}
		if section == "[General]" {
			if replacement, ok := databaseReplacements[trimmed]; ok {
				lines[i] = strings.Replace(line, trimmed, replacement, 1)
				databaseCount++
			}
		}
		if section != "[Remote Rule]" || !strings.HasPrefix(trimmed, "https://") {
			continue
		}
		url, _, ok := strings.Cut(line, ",")
		filename, exists := replacements[url]
		if !ok || !exists || seen[url] {
			return "", fmt.Errorf("无法映射分流订阅: %s", trimmed)
		}
		lines[i] = rawBaseURL + filename + line[len(url):]
		seen[url] = true
		count++
	}
	if count != len(replacements) {
		return "", fmt.Errorf("预期替换六项分流订阅，实际 %d 项", count)
	}
	if databaseCount != len(databaseReplacements) {
		return "", fmt.Errorf("预期替换两项数据库地址，实际 %d 项", databaseCount)
	}
	return cdnGitHubURLs(strings.Join(lines, "")), nil
}
