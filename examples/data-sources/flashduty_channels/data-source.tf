# Get all channels
data "flashduty_channels" "all" {}

# Output all channel names
output "all_channel_names" {
  value = [for channel in data.flashduty_channels.all.channels : channel.channel_name]
}

# Filter channels by team
locals {
  sre_channels = [for ch in data.flashduty_channels.all.channels : ch if ch.team_id == 12345]
}

output "sre_channel_count" {
  value = length(local.sre_channels)
}
