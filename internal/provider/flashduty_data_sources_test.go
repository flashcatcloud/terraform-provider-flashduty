package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeamsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.flashduty_teams.all", "teams.#"),
				),
			},
		},
	})
}

func testAccTeamsDataSourceConfig() string {
	return `
data "flashduty_teams" "all" {}
`
}

func TestAccMembersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMembersDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.flashduty_members.all", "members.#"),
				),
			},
		},
	})
}

func testAccMembersDataSourceConfig() string {
	return `
data "flashduty_members" "all" {}
`
}

func TestAccChannelsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChannelsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.flashduty_channels.all", "channels.#"),
				),
			},
		},
	})
}

func testAccChannelsDataSourceConfig() string {
	return `
data "flashduty_channels" "all" {}
`
}

func TestAccFieldsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFieldsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.flashduty_fields.all", "fields.#"),
				),
			},
		},
	})
}

func testAccFieldsDataSourceConfig() string {
	return `
data "flashduty_fields" "all" {}
`
}

func TestAccTeamDataSource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.flashduty_team.test", "team_name", rName),
				),
			},
		},
	})
}

func testAccTeamDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "flashduty_team" "test" {
  team_name   = %[1]q
  description = "Test team for data source"
}

data "flashduty_team" "test" {
  team_id = tonumber(flashduty_team.test.id)
}
`, name)
}

func TestAccChannelDataSource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChannelDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.flashduty_channel.test", "channel_name", rName),
				),
			},
		},
	})
}

func testAccChannelDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "flashduty_team" "test" {
  team_name   = "%[1]s-team"
  description = "Test team for channel data source"
}

resource "flashduty_channel" "test" {
  channel_name = %[1]q
  description  = "Test channel for data source"
  team_id      = tonumber(flashduty_team.test.id)
}

data "flashduty_channel" "test" {
  channel_id = tonumber(flashduty_channel.test.id)
}
`, name)
}

func TestAccMemberDataSource(t *testing.T) {
	memberID := testAccGetEnv(t, "FLASHDUTY_TEST_MEMBER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMemberDataSourceConfig(memberID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.flashduty_member.test", "member_name"),
				),
			},
		},
	})
}

func testAccMemberDataSourceConfig(memberID string) string {
	return fmt.Sprintf(`
data "flashduty_member" "test" {
  member_id = %s
}
`, memberID)
}

func TestAccFieldDataSource(t *testing.T) {
	// field_name must match ^[a-zA-Z_][a-zA-Z0-9_]{0,39}$ (no hyphens)
	rName := "tf_test_" + acctest.RandStringFromCharSet(8, "abcdefghijklmnopqrstuvwxyz")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFieldDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.flashduty_field.test", "field_name", rName),
				),
			},
		},
	})
}

func testAccFieldDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "flashduty_field" "test" {
  field_name   = %[1]q
  display_name = "Test Field"
  description  = "Test field for data source"
  field_type   = "text"
  value_type   = "string"
}

data "flashduty_field" "test" {
  id = flashduty_field.test.id
}
`, name)
}

func TestAccRouteDataSource(t *testing.T) {
	integrationID := testAccGetEnv(t, "FLASHDUTY_TEST_INTEGRATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteDataSourceConfig(integrationID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.flashduty_route.test", "version"),
				),
			},
		},
	})
}

func testAccRouteDataSourceConfig(integrationID string) string {
	return fmt.Sprintf(`
data "flashduty_route" "test" {
  integration_id = %s
}
`, integrationID)
}

func TestAccRouteHistoryDataSource(t *testing.T) {
	integrationID := testAccGetEnv(t, "FLASHDUTY_TEST_INTEGRATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteHistoryDataSourceConfig(integrationID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.flashduty_route_history.test", "items.#"),
				),
			},
		},
	})
}

func testAccRouteHistoryDataSourceConfig(integrationID string) string {
	return fmt.Sprintf(`
data "flashduty_route_history" "test" {
  integration_id = %s
}
`, integrationID)
}
