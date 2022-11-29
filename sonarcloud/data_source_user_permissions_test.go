package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceUserPermissions(t *testing.T) {
	projectKey := os.Getenv("SONARCLOUD_PROJECT_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserPermissionsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarcloud_user_permissions.test_users", "users.#"),
				),
			},
			{
				Config: testAccDataSourceUserPermissionsConfigForProject(projectKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarcloud_user_permissions.test_users", "project_key"),
					resource.TestCheckResourceAttrSet("data.sonarcloud_user_permissions.test_users", "users.#"),
				),
			},
		},
	})
}

func testAccDataSourceUserPermissionsConfig() string {
	return fmt.Sprintf(`
data "sonarcloud_user_permissions" "test_users" {}
`)
}

func testAccDataSourceUserPermissionsConfigForProject(projectKey string) string {
	return fmt.Sprintf(`
data "sonarcloud_user_permissions" "test_users" {
  project_key = "%s"
}
`, projectKey)
}
