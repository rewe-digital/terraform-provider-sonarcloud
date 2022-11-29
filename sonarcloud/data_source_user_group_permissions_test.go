package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceUserGroupPermissions(t *testing.T) {
	// This is the number of default groups (3) + the "Anyone" permission group (1)
	numberOfDefaultPermissionGroups := "4"
	projectKey := os.Getenv("SONARCLOUD_PROJECT_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserGroupPermissionsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sonarcloud_user_group_permissions.test_groups", "groups.#", numberOfDefaultPermissionGroups),
				),
			},
			{
				Config: testAccDataSourceUserGroupPermissionsConfigForProject(projectKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sonarcloud_user_group_permissions.test_groups", "groups.#", numberOfDefaultPermissionGroups),
				),
			},
		},
	})
}

func testAccDataSourceUserGroupPermissionsConfig() string {
	return fmt.Sprintf(`
data "sonarcloud_user_group_permissions" "test_groups" {}
`)
}

func testAccDataSourceUserGroupPermissionsConfigForProject(projectKey string) string {
	return fmt.Sprintf(`
data "sonarcloud_user_group_permissions" "test_groups" {
  project_key = "%s"
}
`, projectKey)
}
