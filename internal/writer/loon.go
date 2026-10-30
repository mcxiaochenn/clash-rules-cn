package writer

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"clash-rules-cn/internal/types"
)

// ValidateLoonLine 检查不含策略的 Loon 规则，不允许未知规则静默通过。
func ValidateLoonLine(line string) error {
	fields := strings.Split(line, ",")
	if len(fields) < 2 || len(fields) > 3 || strings.TrimSpace(fields[1]) == "" {
		return fmt.Errorf("无效 Loon 规则: %q", line)
	}
	kind, payload := fields[0], fields[1]
	ipRule := kind == "IP-CIDR" || kind == "IP-CIDR6" || kind == "IP-ASN" || kind == "GEOIP"
	if len(fields) == 3 && (!ipRule || fields[2] != "no-resolve") {
		return fmt.Errorf("无效规则参数: %q", line)
	}
	switch kind {
	case "USER-AGENT", "PROCESS-NAME":
		if strings.ContainsAny(payload, "\r\n") {
			return fmt.Errorf("无效匹配模式: %q", line)
		}
	case "DOMAIN", "DOMAIN-SUFFIX", "DOMAIN-KEYWORD":
		if strings.ContainsAny(payload, " \t\r\n:") {
			return fmt.Errorf("无效域名或匹配模式: %q", line)
		}
	case "IP-CIDR", "IP-CIDR6":
		ip, _, err := net.ParseCIDR(payload)
		if err != nil || (ip.To4() == nil) != (kind == "IP-CIDR6") {
			return fmt.Errorf("CIDR 协议族不匹配: %q", line)
		}
	case "IP-ASN":
		asn, err := strconv.ParseUint(payload, 10, 32)
		if err != nil || asn == 0 {
			return fmt.Errorf("无效 ASN: %q", line)
		}
	case "GEOIP":
		if len(payload) != 2 {
			return fmt.Errorf("无效 GEOIP: %q", line)
		}
	default:
		return fmt.Errorf("不支持的 Loon 类型: %q", kind)
	}
	return nil
}

// WriteLoon 保留原生规则类型及参数，由订阅配置指定策略。
func (w *Writer) WriteLoon(filename string, entries []types.Entry) error {
	if len(entries) == 0 {
		return fmt.Errorf("Loon 规则 %s 为空", filename)
	}
	return w.writeFile(filename, func(out *bufio.Writer) error {
		for _, entry := range entries {
			line := entry.Type + "," + entry.Payload
			if err := ValidateLoonLine(line); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(out, line); err != nil {
				return err
			}
		}
		return nil
	})
}
