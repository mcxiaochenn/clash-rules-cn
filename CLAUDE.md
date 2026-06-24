# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 铁律

1. **必须使用中文** — 所有交流、回复、解释一律使用中文，代码注释也优先中文
2. **默认只 commit，绝对不要 push** — 完成修改后执行 `git commit` 即可，用户审查后决定是否推送。只有用户明确说了「推送」或「push」时才执行 `git push`
3. **不确定就问，不要猜** — 任何不确定的事情都要先问用户，不要自作主张

## Git 规范

- 使用 **Conventional Commits** 格式：`feat: 新增 xxx`、`fix: 修复 xxx`、`docs: 更新 xxx`
- 常用类型：`feat`（新功能）、`fix`（修复）、`docs`（文档）、`refactor`（重构）、`chore`（构建/工具变动）、`ci`（CI 配置）

## Project Overview

Clash proxy rule aggregation system. Fetches routing rules from multiple upstream sources, parses/merges/deduplicates them, and outputs YAML rule-provider files for Clash Premium and Clash.Meta/mihomo. A GitHub Actions CI workflow runs daily to build and release.

## Commands

```bash
go build ./cmd/builder          # Build the binary
go run ./cmd/builder             # Run directly
go mod tidy                      # Generate go.sum after dependency changes
go test ./...                    # Run all tests
go test ./internal/parser/...    # Run tests for a single package
```

The compiled binary is `builder` (or `builder.exe` on Windows), output to project root.

## Architecture

Data flows through a linear pipeline:

```
sources.yaml config
    → Fetcher (concurrent HTTP, 3x retry)
    → Parser (one of 7 format-specific parsers)
    → Merger (aggregate by output rule, filter by behavior, deduplicate, sort)
    → Writer (format as YAML rule-provider)
```

**Entry point:** `cmd/builder/main.go` — orchestrates the pipeline.

**Shared type:** `internal/types/types.go` defines `Entry{Type, Payload}` where Type is one of: `DOMAIN`, `DOMAIN-SUFFIX`, `DOMAIN-KEYWORD`, `IP-CIDR`, `IP-CIDR6`.

### Internal packages

| Package | Responsibility |
|---------|---------------|
| `config` | Load and validate `sources.yaml` |
| `fetcher` | Concurrent HTTP downloads with retry logic |
| `parser` | 7 parsers: `cidr`, `dnsmasq`, `domain-list`, `gfwlist`, `v2fly-dlc`, `classical`, `static` |
| `merger` | Aggregate entries by output rule, filter by behavior (`domain`/`ipcidr`/`classical`), deduplicate, sort |
| `writer` | Write YAML rule-provider files with correct formatting per behavior type |
| `types` | Shared `Entry` struct |

### Parser types and their output

| Parser | Input format | Entry.Type produced |
|--------|-------------|-------------------|
| `cidr` | One CIDR per line | `IP-CIDR` / `IP-CIDR6` |
| `dnsmasq` | `server=/domain/ip` | `DOMAIN-SUFFIX` |
| `domain-list` | One domain per line | `DOMAIN-SUFFIX` |
| `gfwlist` | Base64-encoded ABP rules | `DOMAIN-SUFFIX` |
| `v2fly-dlc` | `domain:/full:/keyword:` format | `DOMAIN-SUFFIX` / `DOMAIN` / `DOMAIN-KEYWORD` |
| `classical` | `TYPE,payload` format | preserved as-is |
| `static` | Inline entries from config | auto-detected |

### Writer formatting rules

- `behavior: domain` → DOMAIN-SUFFIX outputs `.domain`, DOMAIN outputs `domain`, DOMAIN-KEYWORD is **skipped** (incompatible)
- `behavior: ipcidr` → outputs CIDR string directly
- `behavior: classical` → outputs `TYPE,payload`

## Configuration

`sources.yaml` defines three sections:
- **`upstream`**: source name → URL, parser type, optional inline entries
- **`output`**: rule name → filename, description, behavior type, list of source names
- **`geoip`**: MMDB/ASN download URLs from `Loyalsoldier/geoip`

## Output

12 rule files in `rules/` (8 domain-type, 4 IP-CIDR-type) plus GeoIP MMDB/ASN files in `geoip/`. See `记录.md` for the full category list.

## Design Notes

- `v2fly/domain-list-community` `include:` directives are not recursively resolved — only direct entries are parsed
- DOMAIN-KEYWORD entries are dropped when writing `behavior: domain` output
- GeoIP files use MMDB+ASN format only (no legacy .dat)
- License: AGPL v3
- Primary design document is `记录.md` (Chinese) — consult it for upstream source details and remaining work items
