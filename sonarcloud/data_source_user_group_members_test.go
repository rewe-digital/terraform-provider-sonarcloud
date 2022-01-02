package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceUserGroupMembers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserGroupMembersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarcloud_user_group_members.test_members", "group"),
					resource.TestCheckResourceAttrSet("data.sonarcloud_user_group_members.test_members", "users.#"),
				),
			},
		},
	})
}

func testAccDataSourceUserGroupMembersConfig() string {
	return fmt.Sprintf(`
data "sonarcloud_user_group_members" "test_members" {
	group = "Members"
}
`)
}
