package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func testAccPreCheckDataSourceProjectLinks(t *testing.T) {
	if v := os.Getenv("SONARCLOUD_PROJECT_KEY"); v == "" {
		t.Fatal("SONARCLOUD_PROJECT_KEY must be set for acceptance tests")
	}
}

func TestAccDataSourceProjectLinks(t *testing.T) {
	expectedNumberOfLinks := "1"

	project := os.Getenv("SONARCLOUD_PROJECT_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); testAccPreCheckDataSourceProjectLinks(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProjectLinksConfig(project),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sonarcloud_project_links.test", "links.#", expectedNumberOfLinks),
					resource.TestCheckResourceAttrSet("data.sonarcloud_project_links.test", "links.0.id"),
					resource.TestCheckResourceAttrSet("data.sonarcloud_project_links.test", "links.0.name"),
					resource.TestCheckResourceAttrSet("data.sonarcloud_project_links.test", "links.0.type"),
					resource.TestCheckResourceAttrSet("data.sonarcloud_project_links.test", "links.0.url"),
				),
			},
		},
	})
}

func testAccDataSourceProjectLinksConfig(projectKey string) string {
	return fmt.Sprintf(`
data "sonarcloud_project_links" "test" {
	project_key = "%s"
}
`, projectKey)
}
