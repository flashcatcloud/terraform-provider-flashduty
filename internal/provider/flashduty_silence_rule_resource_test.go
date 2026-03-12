package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSilenceRuleResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	startTime := time.Now().Unix()
	endTime := startTime + 86400*365

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSilenceRuleResourceConfig(rName, "Test silence rule", startTime, endTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_silence_rule.test", "rule_name", rName),
					resource.TestCheckResourceAttr("flashduty_silence_rule.test", "description", "Test silence rule"),
					resource.TestCheckResourceAttrSet("flashduty_silence_rule.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "flashduty_silence_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["flashduty_silence_rule.test"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["channel_id"], rs.Primary.ID), nil
				},
			},
			// Update and Read testing
			{
				Config: testAccSilenceRuleResourceConfig(rName, "Updated silence rule", startTime, endTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_silence_rule.test", "description", "Updated silence rule"),
				),
			},
		},
	})
}

func testAccSilenceRuleResourceConfig(name, description string, startTime, endTime int64) string {
	return fmt.Sprintf(`
resource "flashduty_team" "test" {
  team_name   = "%[1]s-team"
  description = "Team for silence rule test"
}

resource "flashduty_channel" "test" {
  channel_name = "%[1]s-channel"
  description  = "Channel for silence rule test"
  team_id      = tonumber(flashduty_team.test.id)
}

resource "flashduty_silence_rule" "test" {
  channel_id  = tonumber(flashduty_channel.test.id)
  rule_name   = %[1]q
  description = %[2]q

  filters = [
    {
      conditions = [
        {
          key    = "severity"
          oper   = "IN"
          vals = ["Info", "Warning"]
        }
      ]
    }
  ]

  time_filter = {
    start_time = %[3]d
    end_time   = %[4]d
  }
}
`, name, description, startTime, endTime)
}
