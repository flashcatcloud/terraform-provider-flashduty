# Create a basic team
resource "flashduty_team" "sre" {
  team_name   = "SRE Team"
  description = "Site Reliability Engineering team"
}

# Create a team for a specific service
resource "flashduty_team" "payment_service" {
  team_name   = "Payment Service Team"
  description = "Team responsible for payment service operations"
}
