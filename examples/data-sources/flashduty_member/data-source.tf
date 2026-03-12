# Get member by ID
data "flashduty_member" "john" {
  member_id = 75340551232454
}

# Use the member data
output "member_email" {
  value = data.flashduty_member.john.email
}

output "member_name" {
  value = data.flashduty_member.john.member_name
}

output "member_verified" {
  value = data.flashduty_member.john.email_verified
}
