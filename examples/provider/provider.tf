terraform {
  required_providers {
    flashduty = {
      source = "flashcatcloud/flashduty"
    }
  }
}

provider "flashduty" {
  # APP key can be set via FLASHDUTY_APP_KEY environment variable
  # app_key = "your-app-key"
}
