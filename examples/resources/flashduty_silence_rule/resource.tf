# Create a silence rule with a fixed time window
resource "flashduty_silence_rule" "maintenance" {
  channel_id  = 6148622168454
  rule_name   = "Weekly Maintenance Window"
  description = "Silence alerts during weekly maintenance"

  filters = [
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

  time_filter = {
    start_time = 1773121380
    end_time   = 1773207780
  }
}

# Create a silence rule with recurring time filters
resource "flashduty_silence_rule" "known_issue" {
  channel_id  = 6148622168454
  rule_name   = "Known Disk Space Alert"
  description = "Silence known disk space alerts from legacy server"

  filters = [
    {
      conditions = [
        {
          key  = "labels.host"
          oper = "IN"
          vals = ["legacy-server-01"]
        },
        {
          key  = "title"
          oper = "IN"
          vals = ["disk space"]
        }
      ]
    }
  ]

  time_filters = [
    {
      start  = "02:00"
      end    = "04:00"
      repeat = [0]
    }
  ]

  is_directly_discard = true
}
