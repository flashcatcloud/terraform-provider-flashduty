# Create a team first
resource "flashduty_team" "sre" {
  team_name   = "SRE Team"
  description = "Site Reliability Engineering team"
}

# Basic channel
resource "flashduty_channel" "basic" {
  channel_name = "Production Alerts"
  description  = "Channel for production environment alerts"
  team_id      = tonumber(flashduty_team.sre.id)
}

# Channel with alert grouping (intelligent mode)
resource "flashduty_channel" "with_intelligent_grouping" {
  channel_name = "API Service Alerts"
  description  = "Alerts from API services with intelligent grouping"
  team_id      = tonumber(flashduty_team.sre.id)

  group = {
    method            = "i"
    time_window       = 10
    i_score_threshold = 0.85
    i_keys            = ["title", "description", "labels.service"]
    storm_thresholds  = [100, 500]
  }
}

# Channel with rule-based grouping
resource "flashduty_channel" "with_rule_grouping" {
  channel_name = "Database Alerts"
  description  = "Database alerts with rule-based grouping"
  team_id      = tonumber(flashduty_team.sre.id)

  group = {
    method      = "p"
    time_window = 30
    equals      = [["title", "labels.host"]]

    cases = [
      {
        if = [
          {
            key  = "severity"
            oper = "IN"
            vals = ["Critical"]
          }
        ]
        equals = ["title", "labels.host", "labels.database"]
      }
    ]
  }

  flapping = {
    is_disabled = false
    max_changes = 5
    in_mins     = 60
    mute_mins   = 120
  }
}

# Channel with managing teams and auto-resolve
resource "flashduty_channel" "full" {
  channel_name         = "Infrastructure Alerts"
  description          = "Full-featured channel"
  team_id              = tonumber(flashduty_team.sre.id)
  is_private           = true
  auto_resolve_timeout = 3600
  auto_resolve_mode    = "trigger"
}
