package writer

import (
	"clash-rules-cn/internal/types"
	"os"
	"path/filepath"
	"testing"
)

func TestLoonSyntax(t *testing.T) {
	valid := []string{"DOMAIN,api.example.com", "DOMAIN-SUFFIX,example.com", "DOMAIN-KEYWORD,ai", "IP-CIDR,10.0.0.0/8", "IP-CIDR6,2001:db8::/32", "IP-CIDR,10.0.0.0/8,no-resolve", "IP-ASN,62041,no-resolve", "USER-AGENT,TikTok*", "PROCESS-NAME,aria2", "PROCESS-NAME,Surge 2", "USER-AGENT,Some Agent*"}
	for _, line := range valid {
		if err := ValidateLoonLine(line); err != nil {
			t.Errorf("%s: %v", line, err)
		}
	}
	invalid := []string{"payload:", "  - example.com", "DOMAIN,", "DOMAIN-SUFFIX,http:", "DOMAIN,example.com,DIRECT", "UNKNOWN,example.com", "IP-CIDR,2001:db8::/32", "IP-CIDR6,10.0.0.0/8", "IP-CIDR,not-an-ip", "IP-ASN,0", "IP-ASN,abc", "IP-CIDR,10.0.0.0/8,no-resolve,no-resolve"}
	for _, line := range invalid {
		if err := ValidateLoonLine(line); err == nil {
			t.Errorf("意外接受: %s", line)
		}
	}
}

func TestWriteLoonAndErrors(t *testing.T) {
	dir := t.TempDir()
	w := New(dir)
	if err := w.WriteLoon("test.lsr", []types.Entry{{Type: "DOMAIN", Payload: "api.example.com"}, {Type: "IP-CIDR6", Payload: "2001:db8::/32,no-resolve"}}); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "test.lsr"))
	if string(got) != "DOMAIN,api.example.com\nIP-CIDR6,2001:db8::/32,no-resolve\n" {
		t.Fatalf("格式错误: %q", got)
	}
	if err := w.WriteLoon("empty.lsr", nil); err == nil {
		t.Fatal("空文件必须失败")
	}
	if err := os.Mkdir(filepath.Join(dir, "blocked"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := w.Write("blocked", "", "domain", nil); err == nil {
		t.Fatal("创建文件错误未传播")
	}
	// Unix /dev/full 在 flush 时才报 ENOSPC，验证缓冲错误不会被忽略。
	if _, err := os.Stat("/dev/full"); err == nil {
		for _, loon := range []bool{false, true} {
			full := New("/dev")
			var err error
			if loon {
				err = full.WriteLoon("full", []types.Entry{{Type: "DOMAIN", Payload: "example.com"}})
			} else {
				err = full.Write("full", "", "domain", []types.Entry{{Type: "DOMAIN", Payload: "example.com"}})
			}
			if err == nil {
				t.Fatal("未传播 flush 错误")
			}
		}
	}
}
