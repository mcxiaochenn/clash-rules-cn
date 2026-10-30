// Package loon 保存固定上游配置，仅替换远程分流及 GeoIP / ASN URL。
package loon

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed upstream.lcf
var Upstream string

//go:embed LICENSE
var License []byte

const ConfigFilename = "Loon_Auto_Select_Config.lcf"
const BaseURL = "https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/loon/"

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
		lines[i] = BaseURL + filename + line[len(url):]
		seen[url] = true
		count++
	}
	if count != len(replacements) {
		return "", fmt.Errorf("预期替换六项分流订阅，实际 %d 项", count)
	}
	if databaseCount != len(databaseReplacements) {
		return "", fmt.Errorf("预期替换两项数据库地址，实际 %d 项", databaseCount)
	}
	return strings.Join(lines, ""), nil
}
