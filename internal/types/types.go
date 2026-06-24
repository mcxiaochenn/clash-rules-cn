package types

// Entry represents a single rule entry with its type and payload.
type Entry struct {
	Type    string // DOMAIN, DOMAIN-SUFFIX, DOMAIN-KEYWORD, IP-CIDR, IP-CIDR6
	Payload string
}
