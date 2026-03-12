# Basic notification template
resource "flashduty_template" "basic" {
  template_name = "Default Notification"
  team_id       = 6153171266454
  description   = "Default notification template for all channels"

  email        = "Incident: {{.Title}}"
  sms          = "Incident {{.Title}} triggered"
  dingtalk     = "**Incident**: {{.Title}}"
  wecom        = "**Incident**: {{.Title}}"
  feishu       = "**Incident**: {{.Title}}"
  feishu_app   = "**Incident**: {{.Title}}"
  dingtalk_app = "**Incident**: {{.Title}}"
  wecom_app    = "**Incident**: {{.Title}}"
  teams_app    = "**Incident**: {{.Title}}"
  slack_app    = "*Incident*: {{.Title}}"
  zoom         = "Incident: {{.Title}}"
  telegram     = "*Incident*: {{.Title}}"
}

# Team-scoped template
resource "flashduty_template" "team_scoped" {
  template_name = "SRE Team Notification"
  description   = "Custom template for SRE team"
  team_id       = 6153171266454

  email        = "[SRE] Incident: {{.Title}}"
  sms          = "[SRE] Incident {{.Title}} triggered"
  dingtalk     = "**[SRE] Incident**: {{.Title}}"
  wecom        = "**[SRE] Incident**: {{.Title}}"
  feishu       = "**[SRE] Incident**: {{.Title}}"
  feishu_app   = "**[SRE] Incident**: {{.Title}}"
  dingtalk_app = "**[SRE] Incident**: {{.Title}}"
  wecom_app    = "**[SRE] Incident**: {{.Title}}"
  teams_app    = "**[SRE] Incident**: {{.Title}}"
  slack_app    = "*[SRE] Incident*: {{.Title}}"
  zoom         = "[SRE] Incident: {{.Title}}"
  telegram     = "*[SRE] Incident*: {{.Title}}"
}
