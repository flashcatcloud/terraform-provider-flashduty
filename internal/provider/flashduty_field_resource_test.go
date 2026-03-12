package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFieldResource(t *testing.T) {
	// field_name must match ^[a-zA-Z_][a-zA-Z0-9_]{0,39}$ (no hyphens)
	rName := "tf_test_" + acctest.RandStringFromCharSet(8, "abcdefghijklmnopqrstuvwxyz")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFieldResourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_field.test", "field_name", rName),
					resource.TestCheckResourceAttr("flashduty_field.test", "display_name", "Test Field"),
					resource.TestCheckResourceAttr("flashduty_field.test", "field_type", "single_select"),
					resource.TestCheckResourceAttrSet("flashduty_field.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "flashduty_field.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccFieldResourceConfigUpdated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_field.test", "display_name", "Updated Field"),
				),
			},
		},
	})
}

func testAccFieldResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "flashduty_field" "test" {
  field_name   = %[1]q
  display_name = "Test Field"
  description  = "Test field description"
  field_type   = "single_select"
  value_type   = "string"
  options      = jsonencode(["Option1", "Option2", "Option3"])
}
`, name)
}

func testAccFieldResourceConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "flashduty_field" "test" {
  field_name   = %[1]q
  display_name = "Updated Field"
  description  = "Updated field description"
  field_type   = "single_select"
  value_type   = "string"
  options      = jsonencode(["Option1", "Option2", "Option3", "Option4"])
}
`, name)
}
