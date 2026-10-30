package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"clash-rules-cn/internal/config"
	"clash-rules-cn/internal/merger"
	"clash-rules-cn/internal/types"
	"clash-rules-cn/internal/writer"
)

func TestAllClashOutputsRemainByteIdentical(t *testing.T) {
	cfg, err := config.Load("../../sources.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Output) != 13 || len(cfg.LoonOutput) != 5 {
		t.Fatal("输出分类数量变化")
	}
	before, after := t.TempDir(), t.TempDir()
	for _, rule := range cfg.Output {
		var oldEntries, newEntries []types.Entry
		for _, source := range rule.Sources {
			parserType := cfg.Upstream[source].Parser
			// 同一数据覆盖域名类型、关键词、重复项及 IPv4/IPv6 排序。
			corrected := []types.Entry{{Type: "DOMAIN", Payload: "api.test"}, {Type: "DOMAIN-SUFFIX", Payload: "z.test"}, {Type: "DOMAIN-KEYWORD", Payload: "ai"}, {Type: "IP-CIDR", Payload: "192.168.0.0/16"}, {Type: "IP-CIDR6", Payload: "2001:db8::/32"}, {Type: "IP-CIDR", Payload: "240.0.0.0/4"}}
			legacy := append([]types.Entry(nil), corrected...)
			if parserType == "static" || parserType == "loyalsoldier" {
				legacy[4].Type = "IP-CIDR"
			}
			oldEntries = append(oldEntries, legacy...)
			newEntries = append(newEntries, clashCompatibleEntries(parserType, corrected)...)
		}
		if err := newLegacyWriter(before).Write(rule.Filename, rule.Description, rule.Behavior, merger.Merge(oldEntries, rule.Behavior)); err != nil {
			t.Fatal(err)
		}
		if err := writer.New(after).Write(rule.Filename, rule.Description, rule.Behavior, merger.Merge(newEntries, rule.Behavior)); err != nil {
			t.Fatal(err)
		}
		oldBytes, _ := os.ReadFile(filepath.Join(before, rule.Filename))
		newBytes, _ := os.ReadFile(filepath.Join(after, rule.Filename))
		if !bytes.Equal(oldBytes, newBytes) {
			t.Fatalf("Clash 文件内容变化: %s", rule.Filename)
		}
	}
}
