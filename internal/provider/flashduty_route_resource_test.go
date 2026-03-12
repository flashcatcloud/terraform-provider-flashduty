package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRouteResource(t *testing.T) {
	// Skip: Route resource has no delete API, and creating routes
	// associates channels with integrations, making cleanup difficult.
	// Route tests should be performed manually with proper cleanup.
	t.Skip("Route resource tests require manual cleanup due to no delete API")

	integrationID := testAccGetEnv(t, "FLASHDUTY_TEST_INTEGRATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRouteResourceConfig(integrationID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("flashduty_route.test", "integration_id"),
					resource.TestCheckResourceAttrSet("flashduty_route.test", "version"),
				),
			},
		},
	})
}

func testAccRouteResourceConfig(integrationID string) string {
	return fmt.Sprintf(`
resource "flashduty_route" "test" {
  integration_id = %s

  default = {
    enabled = true
  }
}
`, integrationID)
}
