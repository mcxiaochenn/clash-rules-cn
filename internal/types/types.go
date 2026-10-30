package types

// Entry represents a single rule entry with its type and payload.
type Entry struct {
	Type    string // 域名、CIDR；Loon 专用源还包含 USER-AGENT / IP-ASN
	Payload string // classical 源保留逗号后的 no-resolve 参数
}
