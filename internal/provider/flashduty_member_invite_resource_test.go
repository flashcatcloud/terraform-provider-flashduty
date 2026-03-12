package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMemberInviteResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tftest")
	email := fmt.Sprintf("%s@example.com", rName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMemberInviteResourceConfig(email, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flashduty_member_invite.test", "email", email),
					resource.TestCheckResourceAttrSet("flashduty_member_invite.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "flashduty_member_invite.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMemberInviteResourceConfig(email, name string) string {
	return fmt.Sprintf(`
resource "flashduty_member_invite" "test" {
  email       = %[1]q
  member_name = %[2]q
}
`, email, name)
}
