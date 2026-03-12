## 0.1.0

Initial release:

* **New Provider:** `flashduty` - Terraform provider for managing Flashduty resources via API
* **New Resource:** `flashduty_team` - Manage teams
* **New Resource:** `flashduty_member_invite` - Invite members to the account
* **New Resource:** `flashduty_channel` - Manage collaboration spaces with alert grouping and flap detection
* **New Resource:** `flashduty_schedule` - Configure on-call schedules with rotation rules
* **New Resource:** `flashduty_incident` - Create and manage incidents programmatically
* **New Resource:** `flashduty_escalate_rule` - Define alert escalation rules with layered notification
* **New Resource:** `flashduty_silence_rule` - Configure alert silence rules during maintenance
* **New Resource:** `flashduty_inhibit_rule` - Set up alert inhibition based on conditions
* **New Resource:** `flashduty_field` - Define custom metadata fields for incidents
* **New Resource:** `flashduty_route` - Configure alert routing for shared integrations
* **New Resource:** `flashduty_template` - Manage notification templates across channels
* **New Resource:** `flashduty_alert_pipeline` - Define alert processing pipeline rules (transform, drop, inhibit)
* **New Data Source:** `flashduty_team` - Look up a team by ID
* **New Data Source:** `flashduty_teams` - List teams with filtering
* **New Data Source:** `flashduty_channel` - Look up a channel by ID
* **New Data Source:** `flashduty_channels` - List channels with filtering
* **New Data Source:** `flashduty_member` - Look up a member by ID
* **New Data Source:** `flashduty_members` - List members with filtering
* **New Data Source:** `flashduty_field` - Look up a custom field by ID
* **New Data Source:** `flashduty_fields` - List custom fields
* **New Data Source:** `flashduty_route` - Look up routing rules by integration ID
* **New Data Source:** `flashduty_route_history` - Look up route change history
* **New Data Source:** `flashduty_template` - Look up a notification template by ID
* **New Data Source:** `flashduty_templates` - List notification templates with filtering
* **New Data Source:** `flashduty_alert_pipeline` - Look up alert pipeline rules by integration ID
