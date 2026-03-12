# Look up alert pipeline rules for an integration
data "flashduty_alert_pipeline" "example" {
  integration_id = 5592304204454
}

output "pipeline_rules" {
  value = data.flashduty_alert_pipeline.example.rules
}

output "rule_count" {
  value = length(data.flashduty_alert_pipeline.example.rules)
}
