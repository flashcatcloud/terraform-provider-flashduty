# Get route by integration ID
data "flashduty_route" "main" {
  integration_id = 5592304204454
}

# Use the route data
output "route_version" {
  value = data.flashduty_route.main.version
}

output "route_status" {
  value = data.flashduty_route.main.status
}

output "default_channels" {
  value = data.flashduty_route.main.default != null ? data.flashduty_route.main.default.channel_ids : []
}

output "route_cases_count" {
  value = data.flashduty_route.main.cases != null ? length(data.flashduty_route.main.cases) : 0
}
