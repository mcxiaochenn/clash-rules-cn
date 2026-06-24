# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 铁律

1. **必须使用中文** — 所有交流、回复、解释一律使用中文，代码注释也优先中文
2. **默认只 commit，绝对不要 push** — 完成修改后执行 `git commit` 即可，用户审查后决定是否推送。只有用户明确说了「推送」或「push」时才执行 `git push`
3. **不确定就问，不要猜** — 任何不确定的事情都要先问用户，不要自作主张

## Git 规范

- 使用 **Conventional Commits** 格式：`feat: 新增 xxx`、`fix: 修复 xxx`、`docs: 更新 xxx`
- 常用类型：`feat`（新功能）、`fix`（修复）、`docs`（文档）、`refactor`（重构）、`chore`（构建/工具变动）、`ci`（CI 配置）
- 构建产物（`rules/`、`geoip/`）不提交到 main 分支，通过 CI 推送到 `rules` 分支
- `.gitignore` 中 `/builder` 使用前缀 `/` 避免误匹配 `cmd/builder/` 目录

## Project Overview

Clash proxy rule aggregation system. Fetches routing rules from multiple upstream sources, parses/merges/deduplicates them, and outputs YAML rule-provider files for mihomo. A GitHub Actions CI workflow runs daily to build and release.

## Commands

```bash
go build -o builder.exe ./cmd/builder   # Build the binary (Windows)
go build -o builder ./cmd/builder       # Build the binary (Linux/macOS)
go run ./cmd/builder                     # Run directly
go mod tidy                              # Generate go.sum after dependency changes
go test ./...                            # Run all tests
go test ./internal/parser/...            # Run tests for a single package
```

The compiled binary is `builder` (or `builder.exe` on Windows), output to project root.

## Architecture

Data flows through a linear pipeline:

```
sources.yaml config
    → Fetcher (concurrent HTTP, 3x retry)
    → Parser (one of 8 format-specific parsers)
    → Merger (aggregate by output rule, filter by behavior, deduplicate, sort)
    → Writer (format as YAML rule-provider, bufio streaming)
```

**Entry point:** `cmd/builder/main.go` — orchestrates the pipeline.

**Shared type:** `internal/types/types.go` defines `Entry{Type, Payload}` where Type is one of: `DOMAIN`, `DOMAIN-SUFFIX`, `DOMAIN-KEYWORD`, `IP-CIDR`, `IP-CIDR6`.

### Internal packages

| Package | Responsibility |
|---------|---------------|
| `config` | Load and validate `sources.yaml` |
| `fetcher` | Concurrent HTTP downloads with retry logic |
| `parser` | 8 parsers (see below) |
| `merger` | Aggregate entries by output rule, filter by behavior (`domain`/`ipcidr`/`classical`), deduplicate, sort |
| `writer` | Write YAML rule-provider files with `bufio.Writer` streaming |
| `types` | Shared `Entry` struct |

### Parser types and their output

| Parser | Input format | Entry.Type produced | Notes |
|--------|-------------|-------------------|-------|
| `cidr` | One CIDR per line | `IP-CIDR` / `IP-CIDR6` | |
| `dnsmasq` | `server=/domain/ip` | `DOMAIN-SUFFIX` | |
| `domain-list` | One domain per line | `DOMAIN-SUFFIX` | |
| `gfwlist` | Base64-encoded ABP rules | `DOMAIN-SUFFIX` | |
| `v2fly-dlc` | `domain:/full:/keyword:` + direct domains | `DOMAIN-SUFFIX` / `DOMAIN` / `DOMAIN-KEYWORD` | Supports `include:` recursive parsing with concurrent downloads and global cache |
| `classical` | `TYPE,payload` format | preserved as-is | |
| `loyalsoldier` | YAML `payload:` list | auto-detected: `DOMAIN-SUFFIX` / `IP-CIDR` / `IP-CIDR6` | Handles `+.domain` format |
| `static` | Inline entries from config | auto-detected | |

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

12 rule files in `rules/` (9 domain-type, 3 IP-CIDR-type) plus GeoIP MMDB/ASN files in `geoip/`. Rules are published to the `rules` branch and accessible via CDN:
- `https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/rules/{filename}`
- `https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/rules/{filename}`

## Design Notes

- `v2fly/domain-list-community` `include:` directives ARE recursively resolved (concurrent downloads, global cache across parsers)
- `LoyalsoldierParser` auto-detects CIDR vs domain entries; handles `+.domain` format
- DOMAIN-KEYWORD entries are dropped when writing `behavior: domain` output
- GeoIP files use MMDB+ASN format only (no legacy .dat)
- License: AGPL v3
- Primary design document is `记录.md` (Chinese) — consult it for upstream source details
