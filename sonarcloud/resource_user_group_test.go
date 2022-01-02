package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccUserGroup(t *testing.T) {
	names := []string{"test_group", "test_group_updated"}
	descriptions := []string{"A test group for the SonarCloud provider", "A Group"}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupConfig(names[0], descriptions[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_group.test_group", "name", names[0]),
					resource.TestCheckResourceAttr("sonarcloud_user_group.test_group", "description", descriptions[0]),
				),
			},
			{
				Config: testAccUserGroupConfig(names[0], descriptions[1]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_group.test_group", "name", names[0]),
					resource.TestCheckResourceAttr("sonarcloud_user_group.test_group", "description", descriptions[1]),
				),
			},
			{
				Config: testAccUserGroupConfig(names[1], descriptions[1]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_group.test_group", "name", names[1]),
					resource.TestCheckResourceAttr("sonarcloud_user_group.test_group", "description", descriptions[1]),
				),
			},
		},
		CheckDestroy: testAccUserGroupDestroy,
	})
}

func testAccUserGroupDestroy(s *terraform.State) error {
	return nil
}

func testAccUserGroupConfig(name string, description string) string {
	return fmt.Sprintf(`
resource "sonarcloud_user_group" "test_group" {
	name = "%s"
	description = "%s"
}
`, name, description)
}
