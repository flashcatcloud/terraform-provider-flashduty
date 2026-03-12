# Get team by ID
data "flashduty_team" "by_id" {
  team_id = 12345
}

# Get team by name
data "flashduty_team" "by_name" {
  team_name = "SRE Team"
}

# Get team by external reference ID
data "flashduty_team" "by_ref" {
  ref_id = "sre-team-ref"
}

output "team_name" {
  value = data.flashduty_team.by_id.team_name
}
