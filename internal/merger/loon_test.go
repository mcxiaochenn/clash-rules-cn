package merger

import (
	"reflect"
	"testing"

	"clash-rules-cn/internal/types"
)

func TestLoonMergePreservesSemantics(t *testing.T) {
	input := []types.Entry{{Type: "DOMAIN-SUFFIX", Payload: "z.test # 行内注释"}, {Type: "DOMAIN-KEYWORD", Payload: "ai"}, {Type: "DOMAIN-SUFFIX", Payload: "PROCESS-NAME,aria2"}, {Type: "DOMAIN-SUFFIX", Payload: "z.test,extended-matching"}, {Type: "IP-CIDR", Payload: "10.0.0.0/8"}, {Type: "DOMAIN-SUFFIX", Payload: "https:"}}
	want := []types.Entry{{Type: "DOMAIN-KEYWORD", Payload: "ai"}, {Type: "DOMAIN-SUFFIX", Payload: "z.test"}, {Type: "PROCESS-NAME", Payload: "aria2"}}
	got := MergeLoon(input, "domain")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("合并结果错误: %v", got)
	}
	reversed := append([]types.Entry(nil), input...)
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	if !reflect.DeepEqual(got, MergeLoon(reversed, "domain")) {
		t.Fatal("输出不稳定")
	}
	if input[2].Type != "DOMAIN-SUFFIX" {
		t.Fatal("污染 Clash 输入")
	}
}
