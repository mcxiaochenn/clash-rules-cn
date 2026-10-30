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
	return w.writeFile(filename, func(out *bufio.Writer) error {
		if _, err := out.WriteString("payload:\n"); err != nil {
			return err
		}
		for _, entry := range entries {
			if formatted := w.formatEntry(entry, behavior); formatted != "" {
				if _, err := fmt.Fprintf(out, "  - %s\n", formatted); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// writeFile 检查写入、刷新和关闭错误，避免成功发布截断文件。
func (w *Writer) writeFile(filename string, write func(*bufio.Writer) error) (err error) {
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(w.outputDir, filename))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	out := bufio.NewWriter(f)
	if err := write(out); err != nil {
		return err
	}
	return out.Flush()
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
