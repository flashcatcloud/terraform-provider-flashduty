package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScheduleResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	startTime := time.Now().Unix()
	memberID := testAccGetEnv(t, "FLASHDUTY_TEST_MEMBER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccScheduleResourceConfig(rName, "Test schedule description", startTime, memberID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_schedule.test", "schedule_name", rName),
					resource.TestCheckResourceAttr("flashduty_schedule.test", "description", "Test schedule description"),
					resource.TestCheckResourceAttrSet("flashduty_schedule.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "flashduty_schedule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccScheduleResourceConfig(rName, "Updated description", startTime, memberID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_schedule.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccScheduleResourceConfig(name, description string, startTime int64, memberID string) string {
	return fmt.Sprintf(`
resource "flashduty_schedule" "test" {
  schedule_name = %[1]q
  description   = %[2]q

  layers = [
    {
      layer_name     = "Primary Layer"
      mode           = 0
      rotation_unit  = "day"
      rotation_value = 1
      layer_start    = %[3]d

      groups = [
        {
          group_name = "Group 1"
          members = [
            {
              role_id    = 0
              person_ids = [%[4]s]
            }
          ]
        }
      ]
    }
  ]
}
`, name, description, startTime, memberID)
}
