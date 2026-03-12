# Look up a specific template by ID
data "flashduty_template" "example" {
  template_id = "template-abc123"
}

output "template_name" {
  value = data.flashduty_template.example.template_name
}

output "email_content" {
  value = data.flashduty_template.example.email
}
