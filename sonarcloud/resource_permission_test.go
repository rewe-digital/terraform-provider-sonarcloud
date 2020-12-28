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
	projects := []string{"", org + "_test"}
	groups := []string{"Members"}

	// Possible values for global permissions: admin, profileadmin, gateadmin, scan, provisioning
	// Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user
	permissions := make([][]string, 0)
	permissions = append(permissions, []string{
		"provisioning",
	})
	permissions = append(permissions, []string{
		"provisioning",
		"scan",
	})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionConfig(projects[0], groups[0], permissions[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "project", projects[0]),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "group", groups[0]),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.0", permissions[0][0]),
				),
			},
			{
				Config: testAccPermissionConfig(projects[0], groups[0], permissions[1]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "project", projects[0]),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "group", groups[0]),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.0", permissions[1][0]),
					resource.TestCheckResourceAttr("sonarcloud_permission.test_permission", "permissions.1", permissions[1][1]),
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
