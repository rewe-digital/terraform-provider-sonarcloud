package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"strings"
	"testing"
)

func TestAccPermission(t *testing.T) {
	org := os.Getenv("SONARCLOUD_ORGANIZATION")
	project := org + "_test"
	group := os.Getenv("SONARCLOUD_TEST_GROUP_NAME")

	// Possible values for global permissions: admin, profileadmin, gateadmin, scan, provisioning
	// Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user
	// Note: some permissions (like codeviewer) are active by default on public projects, and are not returned when reading
	// these should not be used in tests when using a public test project
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionConfig("", group, []string{
					"provisioning",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "project", ""),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "group", group),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.0", "provisioning"),
				),
			},
			{
				Config: testAccPermissionConfig("", group, []string{
					"provisioning",
					"scan",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "project", ""),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "group", group),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.0", "provisioning"),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.1", "scan"),
				),
			},
			{
				Config: testAccPermissionConfig("", group, []string{
					"scan",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "project", ""),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "group", group),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.0", "scan"),
				),
			},
			{
				Config: testAccPermissionConfig(project, group, []string{
					"admin",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "project", project),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "group", group),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.0", "admin"),
				),
			},
			{
				Config: testAccPermissionConfig(project, group, []string{
					"issueadmin",
					"scan",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "project", project),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "group", group),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.0", "issueadmin"),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.1", "scan"),
				),
			},
			{
				Config: testAccPermissionConfig(project, group, []string{
					"scan",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "project", project),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "group", group),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.0", "scan"),
				),
			},
		},
		CheckDestroy: testAccPermissionDestroy,
	})
}

func testAccPermissionDestroy(s *terraform.State) error {
	return nil
}

func testAccPermissionConfig(project string, group string, permissions []string) string {
	return fmt.Sprintf(`
resource "sonarcloud_permission" "test_permission" {
	project = "%s"	
	group = "%s"
	permissions = %s
}
`, project, group, permissionsListString(permissions))
}

func permissionsListString(permissions []string) string {
	return fmt.Sprintf(`["%s"]`, strings.Join(permissions, `","`))
}
