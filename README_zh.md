# Terraform Provider for Flashduty

[![Tests](https://github.com/flashcatcloud/terraform-provider-flashduty/actions/workflows/test.yml/badge.svg)](https://github.com/flashcatcloud/terraform-provider-flashduty/actions/workflows/test.yml)
[![Terraform Registry](https://img.shields.io/badge/terraform-registry-blueviolet.svg)](https://registry.terraform.io/providers/flashcatcloud/flashduty/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/flashcatcloud/terraform-provider-flashduty)](https://goreportcard.com/report/github.com/flashcatcloud/terraform-provider-flashduty)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

[English](README.md) | 中文

Flashduty Terraform Provider 允许您使用 [Terraform](https://www.terraform.io) 以基础设施即代码的方式管理 [Flashduty](https://flashcat.cloud/product/flashduty) 资源。

## 支持的资源

### Resources

| 资源 | 描述 |
|------|------|
| `flashduty_team` | 管理团队 |
| `flashduty_member_invite` | 邀请成员加入账户 |
| `flashduty_channel` | 管理协作空间，支持告警聚合和抖动检测 |
| `flashduty_schedule` | 配置值班计划和轮换规则 |
| `flashduty_incident` | 以编程方式创建和管理故障 |
| `flashduty_escalate_rule` | 定义分层通知的告警升级规则 |
| `flashduty_silence_rule` | 配置维护期间的告警静默规则 |
| `flashduty_inhibit_rule` | 设置基于条件的告警抑制规则 |
| `flashduty_field` | 定义故障的自定义元数据字段 |
| `flashduty_route` | 配置共享集成的告警路由 |
| `flashduty_template` | 管理跨渠道的通知模板 |
| `flashduty_alert_pipeline` | 定义告警处理流水线规则（转换、丢弃、抑制） |

### Data Sources

| 数据源 | 描述 |
|--------|------|
| `flashduty_team` / `flashduty_teams` | 查询团队 |
| `flashduty_channel` / `flashduty_channels` | 查询协作空间 |
| `flashduty_member` / `flashduty_members` | 查询成员 |
| `flashduty_field` / `flashduty_fields` | 查询自定义字段 |
| `flashduty_route` / `flashduty_route_history` | 查询路由规则及变更历史 |
| `flashduty_template` / `flashduty_templates` | 查询通知模板 |
| `flashduty_alert_pipeline` | 查询告警处理流水线规则 |

## 环境要求

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25.5（仅构建 Provider 时需要）

## 快速开始

通过 Flashduty APP Key 配置 Provider：

```hcl
terraform {
  required_providers {
    flashduty = {
      source  = "flashcatcloud/flashduty"
      version = "~> 0.1"
    }
  }
}

provider "flashduty" {
  app_key = "your-app-key"
}
```

或者设置环境变量 `FLASHDUTY_APP_KEY`：

```shell
export FLASHDUTY_APP_KEY="your-app-key"
```

### 示例：创建团队并邀请成员

```hcl
resource "flashduty_team" "engineering" {
  team_name   = "Engineering"
  description = "Engineering team"
}

resource "flashduty_member_invite" "example" {
  email       = "user@example.com"
  member_name = "Example User"
}
```

## 文档

完整文档可在 [Terraform Registry](https://registry.terraform.io/providers/flashcatcloud/flashduty/latest/docs) 上查看。

您也可以在本地浏览文档：
- Provider 配置：[`docs/index.md`](docs/index.md)
- 资源：[`docs/resources/`](docs/resources/)
- 数据源：[`docs/data-sources/`](docs/data-sources/)
- 示例：[`examples/`](examples/)

## 开发 Provider

请参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 了解完整的开发指南。

```shell
# 构建
go build -o terraform-provider-flashduty

# 运行测试
export FLASHDUTY_APP_KEY="your-app-key"
make testacc

# 生成/更新文档
make generate
```

## 贡献

欢迎贡献！提交 Pull Request 前请阅读 [贡献指南](CONTRIBUTING.md)。

## 安全

如需报告安全漏洞，请参阅 [SECURITY.md](SECURITY.md)。**请勿通过公开 Issue 报告。**

## 许可证

本项目基于 [Mozilla Public License 2.0](LICENSE) 许可。
