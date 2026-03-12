# Get field by ID
data "flashduty_field" "priority" {
  id = "field-abc123"
}

# Use the field data
output "field_name" {
  value = data.flashduty_field.priority.field_name
}

output "field_display_name" {
  value = data.flashduty_field.priority.display_name
}

output "field_options" {
  value = data.flashduty_field.priority.options != null ? jsondecode(data.flashduty_field.priority.options) : []
}
