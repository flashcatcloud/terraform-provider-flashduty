# Get channel by ID
data "flashduty_channel" "production" {
  channel_id = 6148622168454
}

# Use the channel data
output "channel_name" {
  value = data.flashduty_channel.production.channel_name
}

output "channel_team_id" {
  value = data.flashduty_channel.production.team_id
}
