# Simple daily rotation schedule
resource "flashduty_schedule" "primary_oncall" {
  schedule_name = "Primary On-Call"
  description   = "Primary on-call rotation for the SRE team"

  layers = [
    {
      layer_name     = "Primary Layer"
      mode           = 0
      layer_start    = 1704038400
      rotation_unit  = "day"
      rotation_value = 1
      groups = [
        {
          group_name = "SRE Group"
          members = [
            {
              role_id    = 0
              person_ids = [75340551232454]
            }
          ]
        }
      ]
    }
  ]
}

# Schedule with advanced settings
resource "flashduty_schedule" "advanced" {
  schedule_name = "Advanced Schedule"
  description   = "Schedule with restrictions, day mask, and notifications"

  layers = [
    {
      layer_name     = "Business Hours"
      mode           = 0
      layer_start    = 1704038400
      rotation_unit  = "week"
      rotation_value = 1
      fair_rotation  = true
      restrict_mode  = 1
      restrict_periods = [
        {
          restrict_start = 32400
          restrict_end   = 64800
        }
      ]
      day_mask = {
        repeat = [1, 2, 3, 4, 5]
      }
      groups = [
        {
          group_name = "Weekday Team"
          members = [
            {
              role_id    = 0
              person_ids = [75340551232454]
            }
          ]
        }
      ]
    }
  ]

  notify = {
    advance_in_time = 3600
    by = {
      follow_preference = true
      personal_channels = ["email"]
    }
  }
}
