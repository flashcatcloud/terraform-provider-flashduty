# Create a simple escalation rule
resource "flashduty_escalate_rule" "default" {
  channel_id  = 6148622168454
  rule_name   = "Default Escalation"
  description = "Default escalation policy for all alerts"
  template_id = "6321aad26c12104586a88916"

  layers = [
    {
      escalate_window = 15
      target = {
        person_ids = [75340551232454]
        by = {
          follow_preference = true
          critical          = ["email", "sms", "voice"]
          warning           = ["email", "sms"]
          info              = ["email"]
        }
      }
    },
    {
      escalate_window = 30
      target = {
        schedule_ids = [4973305527804]
        by = {
          follow_preference = true
          critical          = ["email", "sms", "voice"]
          warning           = ["email", "sms"]
          info              = ["email"]
        }
      }
    }
  ]
}

# Create an escalation rule with full configuration
resource "flashduty_escalate_rule" "critical" {
  channel_id  = 6148622168454
  rule_name   = "Critical Alerts Escalation"
  description = "Escalation policy for critical alerts"
  template_id = "6321aad26c12104586a88916"
  aggr_window = 60

  layers = [
    {
      max_times       = 3
      notify_step     = 10
      escalate_window = 15
      force_escalate  = true
      target = {
        person_ids = [75340551232454]
        by = {
          follow_preference = true
          critical          = ["email", "sms", "voice"]
          warning           = ["email", "sms"]
          info              = ["email"]
        }
        webhooks = [
          {
            type     = "wecom"
            settings = jsonencode({ token = "https://qyapi.weixin.qq.com/xxx" })
          }
        ]
      }
    }
  ]

  # Only apply during weekday business hours
  time_filters = [
    {
      start  = "09:00"
      end    = "18:00"
      repeat = [1, 2, 3, 4, 5]
    }
  ]

  # Only match critical severity alerts
  filters = [
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
}
