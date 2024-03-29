package sonarcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceProjectMainBranch(t *testing.T) {
	prefix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	key := prefix + "sonarcloud-provider-acc-test"
	names := []string{"main", "trunk"}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMainBranchConfig(key, names[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_project_main_branch.test", "name", names[0]),
					resource.TestCheckResourceAttr("sonarcloud_project_main_branch.test", "project_key", key),
				),
			},
			projectMainBranchImportCheck("sonarcloud_project_main_branch.test", names[0], key),
			{
				Config: testAccProjectMainBranchConfig(key, names[1]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_project_main_branch.test", "name", names[1]),
					resource.TestCheckResourceAttr("sonarcloud_project_main_branch.test", "project_key", key),
				),
			},
			projectMainBranchImportCheck("sonarcloud_project_main_branch.test", names[1], key),
		},
		CheckDestroy: testAccProjectMainBranchDestroy,
	})
}

func testAccProjectMainBranchDestroy(s *terraform.State) error {
	return nil
}

func testAccProjectMainBranchConfig(project, branchName string) string {
	return fmt.Sprintf(`
resource "sonarcloud_project" "test" {
	name = "%s"
	key = "%s"
	visibility = "public"
}

resource "sonarcloud_project_main_branch" "test" {
	name = "%s"
	project_key = sonarcloud_project.test.key
}
`, project, project, branchName)
}

func projectMainBranchImportCheck(resourceName, name, projectKey string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateId:     fmt.Sprintf("%s,%s", name, projectKey),
		ImportStateVerify: true,
	}
}
