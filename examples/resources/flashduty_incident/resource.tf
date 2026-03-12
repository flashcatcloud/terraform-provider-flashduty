# Create a channel first
resource "flashduty_team" "sre" {
  team_name = "SRE Team"
}

resource "flashduty_channel" "production" {
  channel_name = "Production"
  team_id      = tonumber(flashduty_team.sre.id)
}

# Create a critical incident
resource "flashduty_incident" "database_outage" {
  title             = "Database Connection Failure"
  description       = "Primary database is not responding to connection requests"
  incident_severity = "Critical"
  channel_id        = tonumber(flashduty_channel.production.id)
}

# Create a warning-level incident
resource "flashduty_incident" "high_latency" {
  title             = "API Latency Spike"
  description       = "API response times exceeding SLA thresholds"
  incident_severity = "Warning"
  channel_id        = tonumber(flashduty_channel.production.id)
}
