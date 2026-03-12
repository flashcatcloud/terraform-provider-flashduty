# Inhibit low-severity alerts when critical alert exists
resource "flashduty_inhibit_rule" "critical_inhibits_warning" {
  channel_id  = 6148622168454
  rule_name   = "Critical Inhibits Warning"
  description = "Suppress warning alerts when critical alert is active for same host"

  source_filters = [
    {
      conditions = [
        {
          key  = "severity"
          oper = "IN"
          vals = ["Critical"]
        }
      ]
    }
  ]

  target_filters = [
    {
      conditions = [
        {
          key  = "severity"
          oper = "IN"
          vals = ["Warning", "Info"]
        }
      ]
    }
  ]

  equals = ["host", "service"]

  is_directly_discard = true
}

# Inhibit downstream service alerts when upstream is down
resource "flashduty_inhibit_rule" "upstream_inhibits_downstream" {
  channel_id  = 6148622168454
  rule_name   = "Upstream Inhibits Downstream"
  description = "Suppress downstream alerts when upstream service is down"

  source_filters = [
    {
      conditions = [
        {
          key  = "labels.service"
          oper = "IN"
          vals = ["database", "cache"]
        }
      ]
    }
  ]

  target_filters = [
    {
      conditions = [
        {
          key  = "labels.service"
          oper = "IN"
          vals = ["api", "web"]
        }
      ]
    }
  ]

  equals = ["datacenter"]

  is_directly_discard = true
}
