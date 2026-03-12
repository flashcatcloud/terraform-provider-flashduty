package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEscalateRuleResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	templateID := testAccGetEnv(t, "FLASHDUTY_TEST_TEMPLATE_ID")
	memberID := testAccGetEnv(t, "FLASHDUTY_TEST_MEMBER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEscalateRuleResourceConfig(rName, "Test escalate rule", templateID, memberID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_escalate_rule.test", "rule_name", rName),
					resource.TestCheckResourceAttr("flashduty_escalate_rule.test", "description", "Test escalate rule"),
					resource.TestCheckResourceAttrSet("flashduty_escalate_rule.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "flashduty_escalate_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["flashduty_escalate_rule.test"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["channel_id"], rs.Primary.ID), nil
				},
			},
			// Update and Read testing
			{
				Config: testAccEscalateRuleResourceConfig(rName, "Updated escalate rule", templateID, memberID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_escalate_rule.test", "description", "Updated escalate rule"),
				),
			},
		},
	})
}

func testAccEscalateRuleResourceConfig(name, description, templateID, memberID string) string {
	return fmt.Sprintf(`
resource "flashduty_team" "test" {
  team_name   = "%[1]s-team"
  description = "Team for escalate rule test"
}

resource "flashduty_channel" "test" {
  channel_name = "%[1]s-channel"
  description  = "Channel for escalate rule test"
  team_id      = tonumber(flashduty_team.test.id)
}

resource "flashduty_escalate_rule" "test" {
  channel_id  = tonumber(flashduty_channel.test.id)
  rule_name   = %[1]q
  description = %[2]q
  template_id = %[3]q

  layers = [
    {
      max_times        = 1
      notify_step      = 10
      escalate_window  = 30
      force_escalate   = true
      target = {
        person_ids = [%[4]s]
        by = {
          follow_preference = true
        }
      }
    }
  ]
}
`, name, description, templateID, memberID)
}
