# Get all teams
data "flashduty_teams" "all" {}

# Output all team names
output "all_team_names" {
  value = [for team in data.flashduty_teams.all.teams : team.team_name]
}

# Find a specific team by name using a local
locals {
  sre_team = [for team in data.flashduty_teams.all.teams : team if team.team_name == "SRE Team"][0]
}

output "sre_team_id" {
  value = local.sre_team.team_id
}
