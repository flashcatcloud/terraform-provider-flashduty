# Get route history by integration ID
data "flashduty_route_history" "main" {
  integration_id = 5592304204454
}

# Use the route history data
output "history_count" {
  value = length(data.flashduty_route_history.main.items)
}

output "latest_version" {
  value = length(data.flashduty_route_history.main.items) > 0 ? data.flashduty_route_history.main.items[0].version : 0
}

# Get all versions
output "all_versions" {
  value = [for item in data.flashduty_route_history.main.items : item.version]
}
