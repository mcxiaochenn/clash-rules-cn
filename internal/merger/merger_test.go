package merger

import (
	"reflect"
	"testing"

	"clash-rules-cn/internal/types"
)

// entriesOf 把字符串列表转成 DOMAIN-SUFFIX 条目，方便断言
func entriesOf(payloads ...string) []types.Entry {
	result := make([]types.Entry, 0, len(payloads))
	for _, p := range payloads {
		result = append(result, types.Entry{Type: "DOMAIN-SUFFIX", Payload: p})
	}
	return result
}

func TestApplyExclude命中精确与子域(t *testing.T) {
	// 精确命中 mtalk.google.com、子域命中 a.mtalk.google.com；
	// alt1-mtalk.google.com 不是 mtalk.google.com 的子域（差个 alt1- 前缀），
	// 按 DOMAIN-SUFFIX 语义不应被剔除 —— 这正是 sources.yaml 里要逐条写全 alt1~alt8 的原因。
	input := entriesOf("mtalk.google.com", "a.mtalk.google.com", "alt1-mtalk.google.com", "google.com")
	got := ApplyExclude(input, []string{"mtalk.google.com"})
	want := entriesOf("alt1-mtalk.google.com", "google.com")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("命中判定不符\n got: %v\nwant: %v", got, want)
	}

	// 要把 alt 系一起剔除，必须显式列出
	got = ApplyExclude(input, []string{"mtalk.google.com", "alt1-mtalk.google.com"})
	want = entriesOf("google.com")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("显式列出 alt1 后应全部剔除\n got: %v\nwant: %v", got, want)
	}
}

func TestApplyExclude不误伤相似域名(t *testing.T) {
	// "notmtalk.google.com" 与 "mtalk.google.com.evil.com" 都不应被剔除
	input := entriesOf("notmtalk.google.com", "mtalk.google.com.evil.com", "mtalk.google.com", "mtalkxgoogle.com")
	got := ApplyExclude(input, []string{"mtalk.google.com"})
	want := entriesOf("notmtalk.google.com", "mtalk.google.com.evil.com", "mtalkxgoogle.com")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("误伤判定不符\n got: %v\nwant: %v", got, want)
	}
}

func TestApplyExclude前导点写法等价(t *testing.T) {
	input := entriesOf("gvt1-cn.com", "r1.gvt1-cn.com", "gvt1-cn.com.cn")
	for _, pattern := range []string{"gvt1-cn.com", ".gvt1-cn.com", "  .gvt1-cn.com  "} {
		got := ApplyExclude(input, []string{pattern})
		want := entriesOf("gvt1-cn.com.cn")
		if !reflect.DeepEqual(got, want) {
			t.Errorf("pattern %q 结果不符\n got: %v\nwant: %v", pattern, got, want)
		}
	}
}

func TestApplyExclude大小写不敏感(t *testing.T) {
	input := entriesOf("MTALK.Google.COM", "Example.com")
	got := ApplyExclude(input, []string{"mtalk.google.com"})
	want := entriesOf("Example.com")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("大小写处理不符\n got: %v\nwant: %v", got, want)
	}
}

func TestApplyExclude空清单原样返回(t *testing.T) {
	input := entriesOf("example.com")
	for _, exclude := range [][]string{nil, {}, {"", "  ", "."}} {
		got := ApplyExclude(input, exclude)
		if !reflect.DeepEqual(got, input) {
			t.Errorf("exclude=%v 应原样返回，实际: %v", exclude, got)
		}
	}
}

func TestApplyExclude只作用于域名类条目(t *testing.T) {
	input := []types.Entry{
		{Type: "IP-CIDR", Payload: "1.2.3.0/24"},
		{Type: "IP-CIDR6", Payload: "2001:db8::/32"},
		{Type: "PROCESS-NAME", Payload: "example.com"}, // 进程名撞了域名，不应被剔除
		{Type: "DOMAIN", Payload: "example.com"},
		{Type: "DOMAIN-KEYWORD", Payload: "example.com"},
	}
	got := ApplyExclude(input, []string{"example.com"})
	want := []types.Entry{
		{Type: "IP-CIDR", Payload: "1.2.3.0/24"},
		{Type: "IP-CIDR6", Payload: "2001:db8::/32"},
		{Type: "PROCESS-NAME", Payload: "example.com"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("类型过滤不符\n got: %v\nwant: %v", got, want)
	}
}

func TestApplyExcludeClassical参数截断(t *testing.T) {
	// classical 源的 payload 可能带 no-resolve 之类参数，只应比对逗号前的匹配值
	input := []types.Entry{
		{Type: "DOMAIN", Payload: "mtalk.google.com,no-resolve"},
		{Type: "DOMAIN", Payload: "keep.example.com,no-resolve"},
	}
	got := ApplyExclude(input, []string{"mtalk.google.com"})
	want := []types.Entry{{Type: "DOMAIN", Payload: "keep.example.com,no-resolve"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("参数截断不符\n got: %v\nwant: %v", got, want)
	}
}

func TestApplyExclude保留顺序且不改动剩余条目(t *testing.T) {
	// ApplyExclude 只做减法，剩余条目的顺序与内容必须原样保留
	input := entriesOf("a.example.com", "b.example.com", "drop.me", "c.example.com")
	got := ApplyExclude(input, []string{"drop.me"})
	want := entriesOf("a.example.com", "b.example.com", "c.example.com")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("顺序或内容被改动\n got: %v\nwant: %v", got, want)
	}
}

func TestApplyExclude与Merge配合(t *testing.T) {
	// 端到端：合并后再排除，验证不影响去重与排序
	raw := entriesOf("mtalk.google.com", "mtalk.google.com", "zzz.example.com", "aaa.example.com")
	got := ApplyExclude(Merge(raw, "domain"), []string{"mtalk.google.com"})
	want := Merge(entriesOf("aaa.example.com", "zzz.example.com"), "domain")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("与 Merge 配合不符\n got: %v\nwant: %v", got, want)
	}
}

func TestApplyExcludeLoon路径同样生效(t *testing.T) {
	input := entriesOf("mtalk.google.com", "keep.example.com")
	got := ApplyExclude(MergeLoon(input, "domain"), []string{"mtalk.google.com"})
	if len(got) != 1 || got[0].Payload != "keep.example.com" {
		t.Errorf("Loon 路径排除未生效: %v", got)
	}
}
