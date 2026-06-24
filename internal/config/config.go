package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 表示整个配置文件的结构
type Config struct {
	Upstream map[string]UpstreamSource `yaml:"upstream"`
	Output   map[string]OutputRule     `yaml:"output"`
	GeoIP    GeoIPConfig               `yaml:"geoip"`
}

// UpstreamSource 表示一个上游数据源
type UpstreamSource struct {
	URL     string   `yaml:"url,omitempty"`
	Parser  string   `yaml:"parser"`
	Entries []string `yaml:"entries,omitempty"` // 仅 static 类型使用
}

// OutputRule 表示一个输出规则
type OutputRule struct {
	Filename    string   `yaml:"filename"`
	Description string   `yaml:"description"`
	Behavior    string   `yaml:"behavior"` // domain, ipcidr, classical
	Sources     []string `yaml:"sources"`
}

// GeoIPConfig 表示 GeoIP 配置
type GeoIPConfig struct {
	Files []GeoIPFile `yaml:"files"`
}

// GeoIPFile 表示一个 GeoIP 文件
type GeoIPFile struct {
	URL      string `yaml:"url"`
	Filename string `yaml:"filename"`
}

// Load 从指定路径加载配置文件
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// Validate 验证配置的有效性
func (c *Config) Validate() error {
	// 验证上游源配置
	for name, source := range c.Upstream {
		if source.Parser == "" {
			return fmt.Errorf("上游源 %s 缺少 parser 配置", name)
		}
		if source.Parser != "static" && source.URL == "" {
			return fmt.Errorf("上游源 %s 不是 static 类型但缺少 url", name)
		}
		if source.Parser == "static" && len(source.Entries) == 0 {
			return fmt.Errorf("上游源 %s 是 static 类型但 entries 为空", name)
		}
	}

	// 验证输出规则配置
	for name, rule := range c.Output {
		if rule.Filename == "" {
			return fmt.Errorf("输出规则 %s 缺少 filename", name)
		}
		if rule.Behavior == "" {
			return fmt.Errorf("输出规则 %s 缺少 behavior", name)
		}
		if rule.Behavior != "domain" && rule.Behavior != "ipcidr" && rule.Behavior != "classical" {
			return fmt.Errorf("输出规则 %s 的 behavior 无效: %s", name, rule.Behavior)
		}
		// 验证引用的上游源是否存在
		for _, sourceName := range rule.Sources {
			if _, exists := c.Upstream[sourceName]; !exists {
				return fmt.Errorf("输出规则 %s 引用了不存在的上游源: %s", name, sourceName)
			}
		}
	}

	return nil
}
