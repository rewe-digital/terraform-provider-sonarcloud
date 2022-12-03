package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccUserPermissions(t *testing.T) {
	projectKey := os.Getenv("SONARCLOUD_PROJECT_KEY")
	login := os.Getenv("SONARCLOUD_TEST_USER_LOGIN")

	// Possible values for global permissions: admin, profileadmin, gateadmin, scan, provisioning
	// Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user
	// Note: some permissions (like codeviewer) are active by default on public projects, and are not returned when reading
	// these should not be used in tests when using a public test project
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserPermissionConfig("", login, []string{
					"provisioning",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "project_key", ""),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "login", login),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "permissions.0", "provisioning"),
					resource.TestCheckResourceAttrSet("sonarcloud_user_permissions.test_permission", "name"),
					resource.TestCheckResourceAttrSet("sonarcloud_user_permissions.test_permission", "avatar"),
				),
			},
			userPermissionsImportCheck("sonarcloud_user_permissions.test_permission", login, ""),
			{
				Config: testAccUserPermissionConfig("", login, []string{
					"provisioning",
					"scan",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "project_key", ""),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "login", login),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "permissions.0", "provisioning"),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "permissions.1", "scan"),
				),
			},
			userPermissionsImportCheck("sonarcloud_user_permissions.test_permission", login, ""),
			{
				Config: testAccUserPermissionConfig("", login, []string{
					"scan",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "project_key", ""),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "login", login),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "permissions.0", "scan"),
				),
			},
			userPermissionsImportCheck("sonarcloud_user_permissions.test_permission", login, ""),
			{
				Config: testAccUserPermissionConfig(projectKey, login, []string{
					"admin",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "project_key", projectKey),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "login", login),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "permissions.0", "admin"),
				),
			},
			userPermissionsImportCheck("sonarcloud_user_permissions.test_permission", login, projectKey),
			{
				Config: testAccUserPermissionConfig(projectKey, login, []string{
					"issueadmin",
					"scan",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "project_key", projectKey),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "login", login),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "permissions.0", "issueadmin"),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "permissions.1", "scan"),
				),
			},
			userPermissionsImportCheck("sonarcloud_user_permissions.test_permission", login, projectKey),
			{
				Config: testAccUserPermissionConfig(projectKey, login, []string{
					"scan",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "project_key", projectKey),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "login", login),
					resource.TestCheckResourceAttr("sonarcloud_user_permissions.test_permission", "permissions.0", "scan"),
				),
			},
			userPermissionsImportCheck("sonarcloud_user_permissions.test_permission", login, projectKey),
		},
		CheckDestroy: testAccUserPermissionDestroy,
	})
}

func testAccUserPermissionDestroy(s *terraform.State) error {
	return nil
}

func testAccUserPermissionConfig(project string, login string, permissions []string) string {
	result := fmt.Sprintf(`
resource "sonarcloud_user_permissions" "test_permission" {
	project_key = "%s"
	login = "%s"
	permissions = %s
}
`, project, login, terraformListString(permissions))
	return result
}

func userPermissionsImportCheck(resourceName, name, projectKey string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateId:     fmt.Sprintf("%s,%s", name, projectKey),
		ImportStateVerify: true,
	}
}
