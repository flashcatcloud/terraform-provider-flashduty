# Terraform Provider for Flashduty

[![Tests](https://github.com/flashcatcloud/terraform-provider-flashduty/actions/workflows/test.yml/badge.svg)](https://github.com/flashcatcloud/terraform-provider-flashduty/actions/workflows/test.yml)
[![Terraform Registry](https://img.shields.io/badge/terraform-registry-blueviolet.svg)](https://registry.terraform.io/providers/flashcatcloud/flashduty/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/flashcatcloud/terraform-provider-flashduty)](https://goreportcard.com/report/github.com/flashcatcloud/terraform-provider-flashduty)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

English | [中文](README_zh.md)

The Flashduty Terraform provider allows you to manage [Flashduty](https://flashcat.cloud/product/flashduty) resources as infrastructure-as-code using [Terraform](https://www.terraform.io/).

## Supported Resources

### Resources

| Resource | Description |
|----------|-------------|
| `flashduty_team` | Manage teams |
| `flashduty_member_invite` | Invite members to the account |
| `flashduty_channel` | Manage collaboration spaces with alert grouping and flap detection |
| `flashduty_schedule` | Configure on-call schedules with rotation rules |
| `flashduty_incident` | Create and manage incidents programmatically |
| `flashduty_escalate_rule` | Define alert escalation rules with layered notification |
| `flashduty_silence_rule` | Configure alert silence rules during maintenance |
| `flashduty_inhibit_rule` | Set up alert inhibition based on conditions |
| `flashduty_field` | Define custom metadata fields for incidents |
| `flashduty_route` | Configure alert routing for shared integrations |
| `flashduty_template` | Manage notification templates across channels |
| `flashduty_alert_pipeline` | Define alert processing pipeline rules (transform, drop, inhibit) |

### Data Sources

| Data Source | Description |
|-------------|-------------|
| `flashduty_team` / `flashduty_teams` | Look up teams |
| `flashduty_channel` / `flashduty_channels` | Look up channels |
| `flashduty_member` / `flashduty_members` | Look up members |
| `flashduty_field` / `flashduty_fields` | Look up custom fields |
| `flashduty_route` / `flashduty_route_history` | Look up routing rules and change history |
| `flashduty_template` / `flashduty_templates` | Look up notification templates |
| `flashduty_alert_pipeline` | Look up alert pipeline rules |

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25.5 (only for building the provider)

## Quick Start

Configure the provider with your Flashduty APP key:

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

Or set the `FLASHDUTY_APP_KEY` environment variable:

```shell
export FLASHDUTY_APP_KEY="your-app-key"
```

### Example: Create a Team and Invite a Member

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

## Documentation

Full documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/flashcatcloud/flashduty/latest/docs).

You can also browse the docs locally:
- Provider configuration: [`docs/index.md`](docs/index.md)
- Resources: [`docs/resources/`](docs/resources/)
- Data Sources: [`docs/data-sources/`](docs/data-sources/)
- Examples: [`examples/`](examples/)

## Developing the Provider

See [CONTRIBUTING.md](CONTRIBUTING.md) for full development guidelines.

```shell
# Build
go build -o terraform-provider-flashduty

# Run tests
export FLASHDUTY_APP_KEY="your-app-key"
make testacc

# Generate/update documentation
make generate
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) before submitting a pull request.

## Security

To report a security vulnerability, please see [SECURITY.md](SECURITY.md). **Do not open a public issue.**

## License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).
