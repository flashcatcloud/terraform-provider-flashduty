# Manage alert processing pipeline for an integration
resource "flashduty_alert_pipeline" "example" {
  integration_id = 5592304204454

  rules = [
    # Override title for alerts from Prometheus
    {
      kind = "title_reset"
      if = [
        {
          key  = "labels.source"
          oper = "IN"
          vals = ["prometheus"]
        }
      ]
      settings = jsonencode({
        title = "[Prometheus] $${title}"
      })
    },

    # Escalate Info-level production alerts to Warning
    {
      kind = "severity_reset"
      if = [
        {
          key  = "severity"
          oper = "IN"
          vals = ["Info"]
        },
        {
          key  = "labels.env"
          oper = "IN"
          vals = ["production"]
        }
      ]
      settings = jsonencode({
        severity = "Warning"
      })
    },

    # Drop known noisy test alerts
    {
      kind = "alert_drop"
      if = [
        {
          key  = "title"
          oper = "IN"
          vals = ["test alert", "debug alert"]
        }
      ]
      settings = jsonencode({})
    },

    # Inhibit lower-severity alerts when a Critical alert exists for the same host
    {
      kind = "alert_inhibit"
      settings = jsonencode({
        source_filters = [
          {
            key  = "severity"
            oper = "IN"
            vals = ["Critical"]
          }
        ]
        equals = ["labels.host", "labels.service"]
      })
    }
  ]
}
