# Route configuration for a shared alert integration.
# Each integration has exactly one route (upsert semantics),
# so only define ONE flashduty_route resource per integration_id.
resource "flashduty_route" "main" {
  integration_id = 5592304204454

  cases = [
    {
      if = [
        {
          key  = "severity"
          oper = "IN"
          vals = ["Critical"]
        }
      ]
      channel_ids  = [6148622168454]
      routing_mode = "standard"
    },
    {
      if = [
        {
          key  = "labels.environment"
          oper = "IN"
          vals = ["production", "prod"]
        },
        {
          key  = "labels.team"
          oper = "IN"
          vals = ["sre", "platform"]
        }
      ]
      channel_ids  = [6148622168454]
      routing_mode = "standard"
    }
  ]

  default = {
    channel_ids = [6148622168454]
  }
}
