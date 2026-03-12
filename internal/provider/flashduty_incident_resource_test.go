package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIncidentResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIncidentResourceConfig(rName, "Test incident description", "Info"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_incident.test", "title", rName),
					resource.TestCheckResourceAttr("flashduty_incident.test", "description", "Test incident description"),
					resource.TestCheckResourceAttr("flashduty_incident.test", "incident_severity", "Info"),
					resource.TestCheckResourceAttrSet("flashduty_incident.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "flashduty_incident.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIncidentResourceConfig(rName, "Updated description", "Warning"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_incident.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("flashduty_incident.test", "incident_severity", "Warning"),
				),
			},
		},
	})
}

func testAccIncidentResourceConfig(title, description, severity string) string {
	return fmt.Sprintf(`
resource "flashduty_team" "test" {
  team_name   = "%[1]s-team"
  description = "Team for incident test"
}

resource "flashduty_channel" "test" {
  channel_name = "%[1]s-channel"
  description  = "Channel for incident test"
  team_id      = tonumber(flashduty_team.test.id)
}

resource "flashduty_incident" "test" {
  title             = %[1]q
  description       = %[2]q
  incident_severity = %[3]q
  channel_id        = tonumber(flashduty_channel.test.id)
}
`, title, description, severity)
}
