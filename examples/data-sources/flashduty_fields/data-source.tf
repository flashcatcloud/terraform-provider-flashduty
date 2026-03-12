# Get all fields
data "flashduty_fields" "all" {}

# Output all field names
output "all_field_names" {
  value = [for field in data.flashduty_fields.all.fields : field.display_name]
}

# Find select fields
locals {
  select_fields = [for f in data.flashduty_fields.all.fields : f if f.field_type == "single_select" || f.field_type == "multi_select"]
}

output "select_field_count" {
  value = length(local.select_fields)
}
