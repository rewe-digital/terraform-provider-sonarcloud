package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceUserGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserGroupConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sonarcloud_user_group.test_group", "name", "TEST_DONT_REMOVE"),
				),
			},
		},
	})
}

func testAccDataSourceUserGroupConfig() string {
	return fmt.Sprintf(`
data "sonarcloud_user_group" "test_group" {
	name = "TEST_DONT_REMOVE"
}
`)
}
