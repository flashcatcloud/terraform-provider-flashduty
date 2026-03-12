# Create a priority field with single select
resource "flashduty_field" "priority" {
  field_name    = "priority"
  display_name  = "Priority Level"
  description   = "The priority level of the incident"
  field_type    = "single_select"
  value_type    = "string"
  options       = jsonencode(["P0", "P1", "P2", "P3"])
  default_value = jsonencode("P2")
}

# Create an affected services multi-select field
resource "flashduty_field" "affected_services" {
  field_name   = "affected_services"
  display_name = "Affected Services"
  description  = "Services impacted by this incident"
  field_type   = "multi_select"
  value_type   = "string"
  options      = jsonencode(["API", "Database", "Cache", "CDN", "Auth"])
}

# Create a customer ID text field
resource "flashduty_field" "customer_id" {
  field_name   = "customer_id"
  display_name = "Customer ID"
  description  = "The ID of the affected customer"
  field_type   = "text"
  value_type   = "string"
}

# Create a requires postmortem checkbox
resource "flashduty_field" "requires_postmortem" {
  field_name    = "requires_postmortem"
  display_name  = "Requires Postmortem"
  description   = "Whether this incident requires a postmortem review"
  field_type    = "checkbox"
  value_type    = "bool"
  default_value = jsonencode(false)
}
