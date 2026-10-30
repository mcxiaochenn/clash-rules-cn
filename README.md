# Clash Rules CN  ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/mcxiaochenn/clash-rules-cn/total?logo=github) ![GitHub Downloads (all assets, latest release)](https://img.shields.io/github/downloads/mcxiaochenn/clash-rules-cn/latest/total?logo=github) [![jsdelivr stats](https://data.jsdelivr.com/v1/package/gh/mcxiaochenn/clash-rules-cn/badge?style=rounded)](https://www.jsdelivr.com/package/gh/mcxiaochenn/clash-rules-cn)

Clash 分流规则聚合系统，自动从多个上游源获取规则，合并去重后发布为 YAML rule-provider 文件。使用 GitHub Actions 自动构建（main 分支 push 忽略文档文件 + 每日 UTC 12:00 定时 + 手动触发），构建产物推送到 `rules` 分支并通过 Release 发布，保证规则最新。

## 说明

本项目聚合以下上游数据源，生成适用于 [**mihomo 内核**](https://github.com/MetaCubeX/mihomo) 的规则集（RULE-SET），同时适用于所有使用 mihomo 内核的图形用户界面（GUI）客户端，包括但不限于 [clash-verge-rev](https://github.com/clash-verge-rev/clash-verge-rev)、[Clash Meta for Android](https://github.com/MetaCubeX/ClashMetaForAndroid)、[FlClash](https://github.com/chen08209/FlClash)、[clash-party](https://github.com/mihomo-party-org/clash-party)、[clashmi](https://github.com/KaringX/clashmi)、[OpenClash](https://github.com/vernesong/OpenClash)。

### 上游数据源

本项目的规则数据全部来自以下上游仓库，没有这些项目作者的辛勤付出，就没有本项目。在此向所有上游作者致以诚挚的感谢 🙏

| 仓库 | 作者 | 用途 |
|------|------|------|
| [17mon/china_ip_list](https://github.com/17mon/china_ip_list) | 17mon | 中国大陆 IP 段 |
| [felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list) | felixonmars | 国内直连域名（dnsmasq 格式） |
| [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community) | v2fly | 各分类域名列表（支持 `include:` 递归解析） |
| [Loyalsoldier/clash-rules](https://github.com/Loyalsoldier/clash-rules) | Loyalsoldier | Clash 分流规则（域名 + IP CIDR） |
| [Loyalsoldier/geoip](https://github.com/Loyalsoldier/geoip) | Loyalsoldier | GeoIP 国家库 |
| [P3TERX/GeoLite.mmdb](https://github.com/P3TERX/GeoLite.mmdb) | P3TERX / MaxMind | IP→ASN 数据库 |
| [blackmatrix7/ios_rule_script](https://github.com/blackmatrix7/ios_rule_script) | blackmatrix7 | iOS 分流规则脚本 |
| [gfwlist/gfwlist](https://github.com/gfwlist/gfwlist) | gfwlist | GFW 域名列表 |
| [privacy-protection-tools/anti-AD](https://github.com/privacy-protection-tools/anti-AD) | privacy-protection-tools | 广告域名拦截列表 |
| [dler-io/Rules](https://github.com/dler-io/Rules) | dler-io | AI 服务域名 |
| [SukkaW/Surge](https://github.com/SukkaW/Surge) | SukkaW | AI 服务域名 |
| [ConnersHua/RuleGo](https://github.com/ConnersHua/RuleGo) | ConnersHua | AI 服务域名 |
| [blackmatrix7/ios_rule_script](https://github.com/blackmatrix7/ios_rule_script) (AI) | blackmatrix7 | OpenAI/Gemini/Claude/Copilot 域名 |

## 规则文件地址及使用方式

### 在线地址（URL）

> 如果无法访问域名 `raw.githubusercontent.com`，可以使用第二个地址（`cdn.jsdelivr.net`），但是内容更新会有 12 小时的延迟。以下地址填写在 Clash 配置文件里的 `rule-providers` 里的 `url` 配置项中。

#### 域名类规则（behavior: domain）

- **中国大陆直连域名 direct-domain.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/direct-domain.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/direct-domain.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/direct-domain.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/direct-domain.yaml)
- **需要代理的域名 proxy-domain.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/proxy-domain.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/proxy-domain.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/proxy-domain.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/proxy-domain.yaml)
- **广告及恶意域名 reject-domain.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/reject-domain.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/reject-domain.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/reject-domain.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/reject-domain.yaml)
- **私有网络专用域名 private-domain.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/private-domain.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/private-domain.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/private-domain.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/private-domain.yaml)
- **Apple 在中国大陆可直连的域名 apple-direct.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/apple-direct.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/apple-direct.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/apple-direct.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/apple-direct.yaml)
- **iCloud 域名列表 icloud-domain.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/icloud-domain.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/icloud-domain.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/icloud-domain.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/icloud-domain.yaml)
- **GFWList 域名列表 gfwlist-domain.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/gfwlist-domain.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/gfwlist-domain.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/gfwlist-domain.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/gfwlist-domain.yaml)
- **非中国大陆使用的顶级域名 non-china-tld.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/non-china-tld.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/non-china-tld.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/non-china-tld.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/non-china-tld.yaml)
- **需要直连的常见软件列表 common-software.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/common-software.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/common-software.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/common-software.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/common-software.yaml)
- **AI 服务相关域名 ai-domain.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/ai-domain.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/ai-domain.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/ai-domain.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/ai-domain.yaml)

#### IP 类规则（behavior: ipcidr）

- **Telegram 使用的 IP 地址 telegram-ip.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/telegram-ip.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/telegram-ip.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/telegram-ip.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/telegram-ip.yaml)
- **局域网 IP 及保留 IP 地址 lan-reserved-ip.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/lan-reserved-ip.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/lan-reserved-ip.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/lan-reserved-ip.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/lan-reserved-ip.yaml)
- **中国大陆 IP 地址 china-ip.yaml**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/china-ip.yaml](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/china-ip.yaml)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/china-ip.yaml](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/china-ip.yaml)

#### GeoIP 文件

Loon 的 `geoip-url` 使用 `Country.mmdb`，`ipasn-url` 使用 `GeoLite2-ASN.mmdb`。历史 `Country-asn.mmdb` 的类型仍是 `GeoLite2-Country`，保留现有地址，但不能拿它替代 IP→ASN 查询库。

- **Country.mmdb**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/Country.mmdb](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/Country.mmdb)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/Country.mmdb](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/Country.mmdb)
- **Country-asn.mmdb**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/Country-asn.mmdb](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/Country-asn.mmdb)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/Country-asn.mmdb](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/Country-asn.mmdb)

- **GeoLite2-ASN.mmdb（IP→ASN）**：
  - [Raw 地址](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/GeoLite2-ASN.mmdb)
  - [CDN 地址](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/GeoLite2-ASN.mmdb)
  - 上游：[P3TERX/GeoLite.mmdb](https://github.com/P3TERX/GeoLite.mmdb)，数据归 MaxMind，遵循该上游列明的 GeoLite2 许可条款。

### 使用方式

要想使用本项目的规则集，只需要在 Clash 配置文件中添加如下 `rule-providers` 和 `rules`。

#### Rule Providers 配置方式

```yaml
rule-providers:
  direct-domain:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/direct-domain.yaml"
    path: ./ruleset/direct-domain.yaml
    interval: 86400

  proxy-domain:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/proxy-domain.yaml"
    path: ./ruleset/proxy-domain.yaml
    interval: 86400

  reject-domain:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/reject-domain.yaml"
    path: ./ruleset/reject-domain.yaml
    interval: 86400

  private-domain:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/private-domain.yaml"
    path: ./ruleset/private-domain.yaml
    interval: 86400

  apple-direct:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/apple-direct.yaml"
    path: ./ruleset/apple-direct.yaml
    interval: 86400

  icloud-domain:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/icloud-domain.yaml"
    path: ./ruleset/icloud-domain.yaml
    interval: 86400

  gfwlist-domain:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/gfwlist-domain.yaml"
    path: ./ruleset/gfwlist-domain.yaml
    interval: 86400

  non-china-tld:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/non-china-tld.yaml"
    path: ./ruleset/non-china-tld.yaml
    interval: 86400

  common-software:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/common-software.yaml"
    path: ./ruleset/common-software.yaml
    interval: 86400

  ai-domain:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/ai-domain.yaml"
    path: ./ruleset/ai-domain.yaml
    interval: 86400

  telegram-ip:
    type: http
    behavior: ipcidr
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/telegram-ip.yaml"
    path: ./ruleset/telegram-ip.yaml
    interval: 86400

  lan-reserved-ip:
    type: http
    behavior: ipcidr
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/lan-reserved-ip.yaml"
    path: ./ruleset/lan-reserved-ip.yaml
    interval: 86400

  china-ip:
    type: http
    behavior: ipcidr
    url: "https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/china-ip.yaml"
    path: ./ruleset/china-ip.yaml
    interval: 86400
```

#### 白名单模式 Rules 配置方式（推荐）

- 白名单模式，意为「**没有命中规则的网络流量，统统使用代理**」，适用于服务器线路网络质量稳定、快速，不缺服务器流量的用户。
- 以下配置中，除了 `DIRECT` 和 `REJECT` 是默认存在于 Clash 中的 policy（路由策略/流量处理策略），其余均为自定义 policy，对应配置文件中 `proxies` 或 `proxy-groups` 中的 `name`。如你直接使用下面的 `rules` 规则，则需要在 `proxies` 或 `proxy-groups` 中手动配置一个 `name` 为 `PROXY` 的 policy。
- 如你希望 Apple、iCloud 列表中的域名也走代理，则把对应 `RULE-SET` 的 policy 由 `DIRECT` 改为 `PROXY`。Google 域名已包含在 `proxy-domain` 中，默认为 `PROXY`，无需额外修改。
- 如你不希望进行 DNS 解析，可在 `GEOIP` 规则的最后加上 `,no-resolve`，如 `GEOIP,CN,DIRECT,no-resolve`。

```yaml
rules:
  - RULE-SET,common-software,DIRECT
  - DOMAIN,clash.razord.top,DIRECT
  - DOMAIN,yacd.haishan.me,DIRECT
  - RULE-SET,private-domain,DIRECT
  - RULE-SET,reject-domain,REJECT
  - RULE-SET,icloud-domain,DIRECT
  - RULE-SET,apple-direct,DIRECT
  - RULE-SET,proxy-domain,PROXY
  - RULE-SET,ai-domain,PROXY
  - RULE-SET,direct-domain,DIRECT
  - RULE-SET,lan-reserved-ip,DIRECT
  - RULE-SET,china-ip,DIRECT
  - RULE-SET,telegram-ip,PROXY
  - GEOIP,LAN,DIRECT
  - GEOIP,CN,DIRECT
  - MATCH,PROXY
```

#### 黑名单模式 Rules 配置方式

- 黑名单模式，意为「**只有命中规则的网络流量，才使用代理**」，适用于服务器线路网络质量不稳定或不够快，或服务器流量紧缺的用户。通常也是软路由用户、家庭网关用户的常用模式。
- 以下配置中，除了 `DIRECT` 和 `REJECT` 是默认存在于 Clash 中的 policy（路由策略/流量处理策略），其余均为自定义 policy，对应配置文件中 `proxies` 或 `proxy-groups` 中的 `name`。如你直接使用下面的 `rules` 规则，则需要在 `proxies` 或 `proxy-groups` 中手动配置一个 `name` 为 `PROXY` 的 policy。

```yaml
rules:
  - RULE-SET,common-software,DIRECT
  - DOMAIN,clash.razord.top,DIRECT
  - DOMAIN,yacd.haishan.me,DIRECT
  - RULE-SET,private-domain,DIRECT
  - RULE-SET,reject-domain,REJECT
  - RULE-SET,non-china-tld,PROXY
  - RULE-SET,gfwlist-domain,PROXY
  - RULE-SET,ai-domain,PROXY
  - RULE-SET,telegram-ip,PROXY
  - MATCH,DIRECT
```

## 本地构建

```bash
# 克隆仓库
git clone https://github.com/mcxiaochenn/clash-rules-cn.git
cd clash-rules-cn

# 安装依赖（Go 1.22+）
go mod tidy

# 构建
go build -o builder ./cmd/builder

# 运行
./builder
```

构建产物输出到 `rules/` 目录（域名和 IP 规则文件）和 `geoip/` 目录（GeoIP MMDB/ASN 文件）。

## 项目结构

```
clash-rules-cn/
├── cmd/builder/main.go        # 入口
├── internal/
│   ├── config/                # 配置加载
│   ├── fetcher/               # HTTP 并发拉取（3 次重试）
│   ├── parser/                # 8 种解析器（cidr/dnsmasq/domain-list/gfwlist/v2fly-dlc/classical/loyalsoldier/static）
│   ├── merger/                # 合并去重排序
│   ├── writer/                # YAML 流式输出
│   └── types/                 # 共享类型
├── sources.yaml               # 上游源 + 输出规则配置
├── rules/                     # 生成的规则文件（.gitignore）
└── geoip/                     # GeoIP 文件（.gitignore）
```

## 许可证

[AGPL-3.0](LICENSE)


## Loon 原生规则与自动配置

从现有统一规则条目直接生成 Loon `.lsr`，不经过在线转换服务。现有 13 个 Clash YAML 地址、格式和行为保持兼容，GeoIP 数据库仍发布到 `rules` 分支根目录。

### 自动配置

配置下载地址：

```text
https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/loon/Loon_Auto_Select_Config.lcf
```

先备份 Loon 现有配置，再通过「从 URL 下载」导入，并添加自己的节点订阅。配置来自可莉自动配置固定快照，构建仅替换 `[Remote Rule]` 六个 URL 和 `[General]` 的 GeoIP / ASN 两个数据库 URL，配置中的 GitHub 资源（含解析器和图标）统一使用 `cdn.jsdelivr.net`，其余正文、策略及插件保持原样。上游版本与许可见 [固定快照说明](templates/loon/README.md)。配置模板不会随定时任务自动追随上游修改。

**插件地址仍指向原服务。分流和数据库换源不代表原插件恢复下载，已有缓存的用户请勿因换源而清除插件缓存或重装 Loon。** 本仓库不分发可莉插件。

### 规则订阅地址

所有 Loon 文件位于 `rules` 分支的 `loon/` 目录，配置拉取入口与内部 GitHub 依赖均使用 jsDelivr：

```text
https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/loon/<文件名>
```

| 文件 | 用途 | 自动配置策略 |
|------|------|------|
| `telegram.lsr` | Telegram 域名、IP 与 ASN | 兜底后备策略 |
| `tiktok.lsr` | TikTok 等上游收录服务 | 新国手动策略 |
| `ai-domain.lsr` | AI 服务 | 美国手动策略 |
| `speedtest-international.lsr` | Speedtest 测速服务 | 全球手动策略 |
| `lan.lsr` | 私有域名与局域网、保留 IP | DIRECT |
| `cn-region.lsr` | 国内直连域名与中国大陆 IP | DIRECT |

另外为现有 13 个 YAML 各提供一个同名 `.lsr`：`direct-domain`、`proxy-domain`、`reject-domain`、`private-domain`、`apple-direct`、`icloud-domain`、`gfwlist-domain`、`non-china-tld`、`common-software`、`ai-domain`、`telegram-ip`、`lan-reserved-ip`、`china-ip`。合计 18 个 Loon 规则文件，规则中不写策略名，通过 `[Remote Rule]` 的 `policy` 指定。

例如：

```ini
[Remote Rule]
https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/loon/ai-domain.lsr, policy=美国手动策略, tag=AI, enabled=true
```

### 覆盖范围与兼容说明

- Loon 原生保留域名精确匹配、后缀、关键词、IPv4/IPv6 CIDR；补充上游中的 USER-AGENT、IP-ASN 与 no-resolve 参数也保留。
- `common-software.lsr` 将 applications 来源的进程条目恢复为 PROCESS-NAME，不把进程名当作域名；进程规则的实际适用性由 Loon 平台能力决定。
- Loon 路径清理域名行内注释，并修正 IPv6 类型；Clash 路径保留既有字节输出，包括旧解析的排序和类型兼容行为。旧 Clash 数据语义问题不在本次迁移范围。
- TikTok、Telegram 补充规则及 Speedtest 来源于 blackmatrix7；Speedtest 文件名对应原配置的 Speedtest Intl，但覆盖范围遵循该上游，**并非可莉原规则的逐条等价副本**。其余组合规则也采用本仓库聚合范围。
- `cn-region.lsr` 使用明确的国内域名和 CIDR，而不是可莉原区域分流文件；模板里的原始 GEOIP 注释保留，实际内容以本仓库规则为准。
- Surge 域名规则的 `extended-matching` 在 Loon 中没有对应参数；转换时显式记录日志并保留基本域名匹配，不承诺 CNAME 扩展匹配等价。其他未知类型或参数会阻止发布。
- 旧 GFWList 解析生成的 `http:` / `https:` 伪域名在 Loon 路径显式记录并排除；为满足字节兼容约定，Clash 历史条目仍保留，产物校验只对这两个已知旧条目作兼容放行。
- 本次不自动启用 `reject-domain.lsr`，完整配置仍仅包含原来的六项分流订阅。

CDN 可能有缓存延迟，发布后会验证配置及引用资源能够下载。

### 构建、验证和发布

```bash
go test ./...
go vet ./...
go run ./cmd/builder
go run ./cmd/check
```

`rules/` 下保留 13 个 YAML，并生成 `rules/loon/` 下的 18 个 LSR、一个 LCF 及署名/许可文件。写入、刷新、必需组合源缺失或产物校验失败会终止构建。发布前检查规则类型、协议族、重复条目、空文件、配置引用及 MMDB 标记。

CI 先将产物完整暂存，再发布到 `rules` 分支；Release 同步提供全部格式。保留原有更新触发方式。验证分为固定输入下的 Clash 字节兼容测试、Loon 格式测试、真实上游构建和线上产物检查；真机导入与实际网络体验需单独验证。

### 许可

仓库代码沿用 AGPL-3.0。可莉模板原件及仅换源修改版单独遵循 CC BY-NC-SA 4.0，保留作者署名、上游提交和修改说明，见 [模板许可](templates/loon/LICENSE)。上游规则保留来源说明并遵循各自许可。
