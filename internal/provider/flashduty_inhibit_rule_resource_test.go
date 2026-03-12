package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccInhibitRuleResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccInhibitRuleResourceConfig(rName, "Test inhibit rule"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_inhibit_rule.test", "rule_name", rName),
					resource.TestCheckResourceAttr("flashduty_inhibit_rule.test", "description", "Test inhibit rule"),
					resource.TestCheckResourceAttrSet("flashduty_inhibit_rule.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "flashduty_inhibit_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["flashduty_inhibit_rule.test"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["channel_id"], rs.Primary.ID), nil
				},
			},
			// Update and Read testing
			{
				Config: testAccInhibitRuleResourceConfig(rName, "Updated inhibit rule"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_inhibit_rule.test", "description", "Updated inhibit rule"),
				),
			},
		},
	})
}

func testAccInhibitRuleResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "flashduty_team" "test" {
  team_name   = "%[1]s-team"
  description = "Team for inhibit rule test"
}

resource "flashduty_channel" "test" {
  channel_name = "%[1]s-channel"
  description  = "Channel for inhibit rule test"
  team_id      = tonumber(flashduty_team.test.id)
}

resource "flashduty_inhibit_rule" "test" {
  channel_id  = tonumber(flashduty_channel.test.id)
  rule_name   = %[1]q
  description = %[2]q

  source_filters = [
    {
      conditions = [
        {
          key    = "severity"
          oper   = "IN"
          vals = ["Critical"]
        }
      ]
    }
  ]

  target_filters = [
    {
      conditions = [
        {
          key    = "severity"
          oper   = "IN"
          vals = ["Warning", "Info"]
        }
      ]
    }
  ]

  equals = ["labels.host"]
}
`, name, description)
}
