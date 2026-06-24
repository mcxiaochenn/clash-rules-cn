# Clash Rules CN

Clash 分流规则聚合系统，自动从多个上游源获取规则，合并去重后发布为 YAML rule-provider 文件。

## 功能特性

- 每日自动构建（GitHub Actions）
- 支持 12 种规则分类（8 域名类 + 4 IP 类）
- 从 7 个上游仓库 + 静态数据聚合，支持 v2fly `include:` 递归解析
- 并发拉取，3 次重试，全局缓存加速
- 自动去重排序，约 8 秒完成全流程
- 包含 GeoIP MMDB/ASN 文件

## 规则分类

| 分类 | 文件名 | behavior |
|------|--------|----------|
| 中国大陆直连域名 | direct-domain.yaml | domain |
| 需要代理的域名 | proxy-domain.yaml | domain |
| 广告及恶意域名 | reject-domain.yaml | domain |
| 私有网络专用域名 | private-domain.yaml | domain |
| Apple 直连域名 | apple-direct.yaml | domain |
| iCloud 域名 | icloud-domain.yaml | domain |
| GFWList 域名 | gfwlist-domain.yaml | domain |
| 非中国 TLD | non-china-tld.yaml | domain |
| 常见软件直连 | common-software.yaml | domain |
| Telegram IP | telegram-ip.yaml | ipcidr |
| 局域网/保留 IP | lan-reserved-ip.yaml | ipcidr |
| 中国大陆 IP | china-ip.yaml | ipcidr |

## 使用方法

### 直接使用

从 [Releases](../../releases) 下载最新的规则文件，或直接引用 `rules` 分支的文件。

### Clash 配置示例

```yaml
rule-providers:
  direct-domain:
    type: http
    behavior: domain
    url: "https://raw.githubusercontent.com/mcxiaochenn/clash-rules-cn/rules/rules/direct-domain.yaml"
    path: ./ruleset/direct-domain.yaml
    interval: 86400

rules:
  - RULE-SET,direct-domain,DIRECT
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

## 项目结构

```
clash-rules-cn/
├── cmd/builder/main.go        # 入口
├── internal/
│   ├── config/                # 配置加载
│   ├── fetcher/               # HTTP 并发拉取
│   ├── parser/                # 8 种解析器
│   ├── merger/                # 合并去重
│   ├── writer/                # YAML 输出
│   └── types/                 # 共享类型
├── sources.yaml               # 上游源配置
├── rules/                     # 生成的规则文件
└── geoip/                     # GeoIP 文件
```

## 许可证

[AGPL-3.0](LICENSE)

## 致谢

本项目的规则数据来自以下上游仓库，感谢所有作者的贡献：

| 仓库 | 作者 | 用途 |
|------|------|------|
| [17mon/china_ip_list](https://github.com/17mon/china_ip_list) | 17mon | 中国大陆 IP 段 |
| [felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list) | felixonmars | 国内直连域名（dnsmasq 格式） |
| [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community) | v2fly | 各分类域名列表 |
| [Loyalsoldier/clash-rules](https://github.com/Loyalsoldier/clash-rules) | Loyalsoldier | Clash 分流规则（域名 + IP CIDR） |
| [Loyalsoldier/geoip](https://github.com/Loyalsoldier/geoip) | Loyalsoldier | GeoIP MMDB/ASN 数据 |
| [blackmatrix7/ios_rule_script](https://github.com/blackmatrix7/ios_rule_script) | blackmatrix7 | iOS 分流规则脚本 |
| [gfwlist/gfwlist](https://github.com/gfwlist/gfwlist) | gfwlist | GFW 域名列表 |
| [privacy-protection-tools/anti-AD](https://github.com/privacy-protection-tools/anti-AD) | privacy-protection-tools | 广告域名拦截列表 |
