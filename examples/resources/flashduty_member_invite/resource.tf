# Invite a new member with email and phone
resource "flashduty_member_invite" "john" {
  email        = "john.doe@example.com"
  member_name  = "John Doe"
  country_code = "CN"
  phone        = "10000000000"
}

# Invite member with email only
resource "flashduty_member_invite" "alice" {
  email       = "alice@example.com"
  member_name = "Alice Smith"
}
