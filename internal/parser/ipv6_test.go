package parser

import "testing"

func TestIPv6CIDRClassification(t *testing.T) {
	input := []string{"10.0.0.0/8", "2001:db8::/32"}
	parsed := ParseStatic(input)
	remote, err := (&LoyalsoldierParser{}).Parse([]byte("payload:\n  - 10.0.0.0/8\n  - 2001:db8::/32\n"))
	if err != nil {
		t.Fatal(err)
	}
	if parsed[0].Type != "IP-CIDR" || parsed[1].Type != "IP-CIDR6" || remote[0].Type != "IP-CIDR" || remote[1].Type != "IP-CIDR6" {
		t.Fatalf("IPv6 分类错误: %v %v", parsed, remote)
	}
}
