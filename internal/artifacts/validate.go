// Package artifacts 检查发布前的规则、配置和数据库完整性。
package artifacts

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"clash-rules-cn/internal/config"
	"clash-rules-cn/internal/writer"
	loontemplate "clash-rules-cn/templates/loon"
	"gopkg.in/yaml.v3"
)

func Validate(cfg *config.Config, rulesDir, geoDir string) error {
	loonDir := filepath.Join(rulesDir, "loon")
	expected := make(map[string]bool)
	for _, rule := range cfg.Output {
		data, err := os.ReadFile(filepath.Join(rulesDir, rule.Filename))
		if err != nil {
			return err
		}
		var parsed struct {
			Payload []yaml.Node `yaml:"payload"`
		}
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return fmt.Errorf("%s: %w", rule.Filename, err)
		}
		if len(parsed.Payload) == 0 {
			return fmt.Errorf("Clash 文件 %s 为空", rule.Filename)
		}
		// 原 Clash GFWList 解析把协议头写为 .http: / .https:，YAML 会解释成映射。
		// 仅兼容这两个历史条目；不改变旧文件字节，Loon 路径已显式排除。
		for _, node := range parsed.Payload {
			if node.Kind == yaml.ScalarNode {
				continue
			}
			if rule.Filename == "gfwlist-domain.yaml" && node.Kind == yaml.MappingNode && len(node.Content) == 2 && (node.Content[0].Value == ".http" || node.Content[0].Value == ".https") && node.Content[1].Tag == "!!null" {
				log.Printf("兼容保留 %s 的历史协议伪域名 %s，Loon 已排除", rule.Filename, node.Content[0].Value)
				continue
			}
			return fmt.Errorf("%s 包含非字符串规则", rule.Filename)
		}
		expected[strings.TrimSuffix(rule.Filename, filepath.Ext(rule.Filename))+".lsr"] = true
	}
	for _, rule := range cfg.LoonOutput {
		expected[rule.Filename] = true
	}
	files, err := filepath.Glob(filepath.Join(loonDir, "*.lsr"))
	if err != nil {
		return err
	}
	if len(files) != len(expected) {
		return fmt.Errorf("Loon 文件数不符: 实际 %d，预期 %d", len(files), len(expected))
	}
	for filename := range expected {
		if err := validateRuleFile(filepath.Join(loonDir, filename)); err != nil {
			return err
		}
	}
	want, err := loontemplate.Generate()
	if err != nil {
		return err
	}
	got, err := os.ReadFile(filepath.Join(loonDir, loontemplate.ConfigFilename))
	if err != nil {
		return err
	}
	if string(got) != want {
		return fmt.Errorf("配置正文不符合六项分流和两项数据库 URL 替换约定")
	}
	section, count := "", 0
	for _, line := range strings.Split(want, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") {
			section = trimmed
		}
		if section == "[Remote Rule]" && strings.HasPrefix(trimmed, "https://") {
			url, _, _ := strings.Cut(trimmed, ",")
			if !strings.HasPrefix(url, loontemplate.BaseURL) || !expected[strings.TrimPrefix(url, loontemplate.BaseURL)] {
				return fmt.Errorf("配置引用未发布规则: %s", url)
			}
			count++
		}
	}
	if count != 6 {
		return fmt.Errorf("配置分流订阅不是六项")
	}
	for _, filename := range []string{"LICENSE", "README.md"} {
		body, err := os.ReadFile(filepath.Join(loonDir, filename))
		if err != nil {
			return err
		}
		if len(body) == 0 {
			return fmt.Errorf("缺少署名或许可内容: %s", filename)
		}
	}
	for _, file := range cfg.GeoIP.Files {
		body, err := os.ReadFile(filepath.Join(geoDir, file.Filename))
		if err != nil {
			return err
		}
		marker := []byte("\xab\xcd\xefMaxMind.com")
		metadataStart := bytes.LastIndex(body, marker)
		if metadataStart < 0 {
			return fmt.Errorf("无效 MMDB: %s", file.Filename)
		}
		expectedType := "GeoLite2-Country"
		if file.Filename == "GeoLite2-ASN.mmdb" {
			expectedType = "GeoLite2-ASN"
		}
		if !bytes.Contains(body[metadataStart+len(marker):], []byte(expectedType)) {
			return fmt.Errorf("数据库 %s 类型错误，预期 %s", file.Filename, expectedType)
		}
	}
	return nil
}

func validateRuleFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	seen := make(map[string]bool)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if err := writer.ValidateLoonLine(line); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if seen[line] {
			return fmt.Errorf("%s 存在重复规则: %s", path, line)
		}
		seen[line] = true
		count++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Loon 文件 %s 为空", path)
	}
	return nil
}
