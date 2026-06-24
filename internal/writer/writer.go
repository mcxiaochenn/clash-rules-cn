package writer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"clash-rules-cn/internal/types"
)

// Writer 负责将规则写入 YAML 文件
type Writer struct {
	outputDir string
}

// New 创建一个新的 Writer
func New(outputDir string) *Writer {
	return &Writer{
		outputDir: outputDir,
	}
}

// Write 将条目写入指定的文件（使用 bufio.Writer 流式写入，避免内存中构建大字符串）
func (w *Writer) Write(filename string, description string, behavior string, entries []types.Entry) error {
	// 确保输出目录存在
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	filePath := filepath.Join(w.outputDir, filename)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	// 写入头部
	writer.WriteString("payload:\n")

	// 根据 behavior 格式化并写入条目
	for _, entry := range entries {
		formatted := w.formatEntry(entry, behavior)
		if formatted != "" {
			writer.WriteString("  - ")
			writer.WriteString(formatted)
			writer.WriteString("\n")
		}
	}

	return nil
}

// formatEntry 根据 behavior 格式化单个条目
func (w *Writer) formatEntry(entry types.Entry, behavior string) string {
	switch behavior {
	case "domain":
		return w.formatDomain(entry)
	case "ipcidr":
		return w.formatIPCIDR(entry)
	case "classical":
		return w.formatClassical(entry)
	default:
		return ""
	}
}

// formatDomain 格式化域名类型的条目
func (w *Writer) formatDomain(entry types.Entry) string {
	switch entry.Type {
	case "DOMAIN":
		return entry.Payload
	case "DOMAIN-SUFFIX":
		return "." + entry.Payload
	case "DOMAIN-KEYWORD":
		// DOMAIN-KEYWORD 在 domain behavior 中不兼容，跳过
		return ""
	default:
		return ""
	}
}

// formatIPCIDR 格式化 IP-CIDR 类型的条目
func (w *Writer) formatIPCIDR(entry types.Entry) string {
	if entry.Type == "IP-CIDR" || entry.Type == "IP-CIDR6" {
		return entry.Payload
	}
	return ""
}

// formatClassical 格式化 classical 类型的条目
func (w *Writer) formatClassical(entry types.Entry) string {
	return entry.Type + "," + entry.Payload
}
