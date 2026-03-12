# List all templates
data "flashduty_templates" "all" {}

# Filter templates by team
data "flashduty_templates" "team_templates" {
  team_ids = [12345]
}

# Search templates by name
data "flashduty_templates" "search" {
  query = "production"
}

output "all_template_names" {
  value = [for t in data.flashduty_templates.all.templates : t.template_name]
}
