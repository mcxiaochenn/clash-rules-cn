# Clash Rules CN  ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/mcxiaochenn/clash-rules-cn/total?logo=github) ![GitHub Downloads (all assets, latest release)](https://img.shields.io/github/downloads/mcxiaochenn/clash-rules-cn/latest/total?logo=github) [![jsdelivr stats](https://data.jsdelivr.com/v1/package/gh/mcxiaochenn/clash-rules-cn/badge?style=rounded)](https://www.jsdelivr.com/package/gh/mcxiaochenn/clash-rules-cn)

Clash 分流规则聚合系统，自动从多个上游源获取规则，合并去重后发布为 YAML rule-provider 文件。使用 GitHub Actions 每日自动构建，保证规则最新。

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
| [Loyalsoldier/geoip](https://github.com/Loyalsoldier/geoip) | Loyalsoldier | GeoIP MMDB/ASN 数据 |
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

- **Country.mmdb**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/Country.mmdb](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/Country.mmdb)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/Country.mmdb](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/Country.mmdb)
- **Country-asn.mmdb**：
  - [https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/Country-asn.mmdb](https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/Country-asn.mmdb)
  - [https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/Country-asn.mmdb](https://cdn.jsdelivr.net/gh/mcxiaochenn/clash-rules-cn@rules/Country-asn.mmdb)

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
- 如你希望 Apple、iCloud 和 Google 列表中的域名使用代理，则把 policy 由 `DIRECT` 改为 `PROXY`，以此类推，举一反三。
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
│   ├── parser/                # 8 种解析器
│   ├── merger/                # 合并去重排序
│   ├── writer/                # YAML 流式输出
│   └── types/                 # 共享类型
├── sources.yaml               # 上游源 + 输出规则配置
├── rules/                     # 生成的规则文件（.gitignore）
└── geoip/                     # GeoIP 文件（.gitignore）
```

## 许可证

[AGPL-3.0](LICENSE)
