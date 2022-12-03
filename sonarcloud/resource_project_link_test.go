package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccProjectLink(t *testing.T) {
	projectKey := os.Getenv("SONARCLOUD_PROJECT_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectLinkConfig(projectKey, "test", "https://www.example.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_project_link.test", "project_key", projectKey),
					resource.TestCheckResourceAttr("sonarcloud_project_link.test", "name", "test"),
					resource.TestCheckResourceAttr("sonarcloud_project_link.test", "url", "https://www.example.com"),
				),
			},
			projectLinkImportCheck("sonarcloud_project_link.test", projectKey),
			{
				Config: testAccProjectLinkConfig(projectKey, "test", "https://www.iana.org/domains/reserved"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_project_link.test", "project_key", projectKey),
					resource.TestCheckResourceAttr("sonarcloud_project_link.test", "name", "test"),
					resource.TestCheckResourceAttr("sonarcloud_project_link.test", "url", "https://www.iana.org/domains/reserved"),
				),
			},
			projectLinkImportCheck("sonarcloud_project_link.test", projectKey),
			{
				Config: testAccProjectLinkConfig(projectKey, "test-two", "https://www.iana.org/domains/reserved"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_project_link.test", "project_key", projectKey),
					resource.TestCheckResourceAttr("sonarcloud_project_link.test", "name", "test-two"),
					resource.TestCheckResourceAttr("sonarcloud_project_link.test", "url", "https://www.iana.org/domains/reserved"),
				),
			},
			projectLinkImportCheck("sonarcloud_project_link.test", projectKey),
		},
		CheckDestroy: testAccLinkDestroy,
	})
}

func testAccLinkDestroy(s *terraform.State) error {
	return nil
}

func testAccProjectLinkConfig(projectKey, name, url string) string {
	result := fmt.Sprintf(`
resource "sonarcloud_project_link" "test" {
	project_key = "%s"
	name        = "%s"
    url         = "%s"
}
`, projectKey, name, url)
	return result
}

func projectLinkImportCheck(resourceName, projectKey string) resource.TestStep {
	return resource.TestStep{
		ResourceName: resourceName,
		ImportState:  true,
		ImportStateIdFunc: func(state *terraform.State) (string, error) {
			id := state.RootModule().Resources[resourceName].Primary.ID
			return fmt.Sprintf("%s,%s", id, projectKey), nil
		},
		ImportStateVerify: true,
	}
}
