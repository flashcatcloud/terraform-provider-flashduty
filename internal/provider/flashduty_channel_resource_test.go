package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccChannelResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccChannelResourceConfig(rName, "Test channel description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_channel.test", "channel_name", rName),
					resource.TestCheckResourceAttr("flashduty_channel.test", "description", "Test channel description"),
					resource.TestCheckResourceAttrSet("flashduty_channel.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "flashduty_channel.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccChannelResourceConfig(rName, "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_channel.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccChannelResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "flashduty_team" "test" {
  team_name   = "%[1]s-team"
  description = "Team for channel test"
}

resource "flashduty_channel" "test" {
  channel_name = %[1]q
  description  = %[2]q
  team_id      = tonumber(flashduty_team.test.id)
}
`, name, description)
}
