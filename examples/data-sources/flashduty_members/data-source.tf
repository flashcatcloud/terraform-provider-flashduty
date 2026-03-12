# Get all members
data "flashduty_members" "all" {}

# Output all member names
output "all_member_names" {
  value = [for member in data.flashduty_members.all.members : member.member_name]
}

# Find active members
locals {
  active_members = [for m in data.flashduty_members.all.members : m if m.status == "active"]
}

output "active_member_count" {
  value = length(local.active_members)
}

# Find member by email
locals {
  john = [for m in data.flashduty_members.all.members : m if m.email == "john@example.com"][0]
}

output "john_id" {
  value = local.john.member_id
}
