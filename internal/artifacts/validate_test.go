package artifacts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"clash-rules-cn/internal/config"
	loontemplate "clash-rules-cn/templates/loon"
)

func TestPublishGate(t *testing.T) {
	cfg, err := config.Load("../../sources.yaml")
	if err != nil {
		t.Fatal(err)
	}
	cases := []string{"valid", "missing-rule", "empty-rule", "invalid-rule", "duplicate-rule", "changed-config", "missing-mmdb", "invalid-mmdb", "empty-yaml", "extra-rule", "malformed-yaml", "legacy-gfw", "wrong-asn-type"}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			rulesDir := filepath.Join(root, "rules")
			loonDir := filepath.Join(rulesDir, "loon")
			geoDir := filepath.Join(root, "geoip")
			if err := os.MkdirAll(loonDir, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(geoDir, 0755); err != nil {
				t.Fatal(err)
			}
			write := func(path, body string) {
				t.Helper()
				if err := os.WriteFile(path, []byte(body), 0644); err != nil {
					t.Fatal(err)
				}
			}
			for _, rule := range cfg.Output {
				write(filepath.Join(rulesDir, rule.Filename), "payload:\n  - example.com\n")
				write(filepath.Join(loonDir, strings.TrimSuffix(rule.Filename, ".yaml")+".lsr"), "DOMAIN,example.com\n")
			}
			for _, rule := range cfg.LoonOutput {
				write(filepath.Join(loonDir, rule.Filename), "DOMAIN,example.com\n")
			}
			generated, err := loontemplate.Generate()
			if err != nil {
				t.Fatal(err)
			}
			write(filepath.Join(loonDir, loontemplate.ConfigFilename), generated)
			write(filepath.Join(loonDir, "LICENSE"), string(loontemplate.License))
			write(filepath.Join(loonDir, "README.md"), "署名")
			for _, file := range cfg.GeoIP.Files {
				dbType := "GeoLite2-Country"
				if file.Filename == "GeoLite2-ASN.mmdb" {
					dbType = "GeoLite2-ASN"
				}
				write(filepath.Join(geoDir, file.Filename), "test\xab\xcd\xefMaxMind.com"+dbType)
			}
			target := filepath.Join(loonDir, "telegram.lsr")
			switch name {
			case "missing-rule":
				if err := os.Remove(target); err != nil {
					t.Fatal(err)
				}
			case "empty-rule":
				write(target, "")
			case "invalid-rule":
				write(target, "payload:\n  - example.com\n")
			case "duplicate-rule":
				write(target, "DOMAIN,example.com\nDOMAIN,example.com\n")
			case "changed-config":
				write(filepath.Join(loonDir, loontemplate.ConfigFilename), strings.Replace(generated, "disable-stun=true", "disable-stun=false", 1))
			case "missing-mmdb":
				if err := os.Remove(filepath.Join(geoDir, cfg.GeoIP.Files[0].Filename)); err != nil {
					t.Fatal(err)
				}
			case "wrong-asn-type":
				write(filepath.Join(geoDir, "GeoLite2-ASN.mmdb"), "test\xab\xcd\xefMaxMind.comGeoLite2-Country")
			case "invalid-mmdb":
				write(filepath.Join(geoDir, cfg.GeoIP.Files[0].Filename), "not mmdb")
			case "empty-yaml":
				write(filepath.Join(rulesDir, "ai-domain.yaml"), "payload: []\n")
			case "legacy-gfw":
				write(filepath.Join(rulesDir, "gfwlist-domain.yaml"), "payload:\n  - .http:\n  - .https:\n")
			case "malformed-yaml":
				write(filepath.Join(rulesDir, "ai-domain.yaml"), "payload:\n  - unexpected: value\n")
			case "extra-rule":
				write(filepath.Join(loonDir, "stale.lsr"), "DOMAIN,example.com\n")
			}
			err = Validate(cfg, rulesDir, geoDir)
			if (name == "valid" || name == "legacy-gfw") != (err == nil) {
				t.Fatalf("校验结果不符: %v", err)
			}
		})
	}
}
