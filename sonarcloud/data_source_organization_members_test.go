package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceOrganizationMembers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizationMembersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarcloud_organization_members.test_members", "users.#"),
				),
			},
		},
	})
}

func testAccDataSourceOrganizationMembersConfig() string {
	return fmt.Sprintf(`
data "sonarcloud_organization_members" "test_members" {}
`)
}
