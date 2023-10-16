package sonarcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccPreCheckDataSourceWebhooks(t *testing.T) {
	if v := os.Getenv("SONARCLOUD_PROJECT_KEY"); v == "" {
		t.Fatal("SONARCLOUD_PROJECT_KEY must be set for acceptance tests")
	}
}

func TestAccDataSourceWebhooks(t *testing.T) {
	expectedNumberOfWebhooks := "1"

	project := os.Getenv("SONARCLOUD_PROJECT_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); testAccPreCheckDataSourceWebhooks(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWebhooksConfig(project),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sonarcloud_webhooks.test", "webhooks.#", expectedNumberOfWebhooks),
					resource.TestCheckResourceAttrSet("data.sonarcloud_webhooks.test", "webhooks.0.key"),
					resource.TestCheckResourceAttrSet("data.sonarcloud_webhooks.test", "webhooks.0.name"),
					resource.TestCheckResourceAttrSet("data.sonarcloud_webhooks.test", "webhooks.0.url"),
					resource.TestCheckResourceAttrSet("data.sonarcloud_webhooks.test", "webhooks.0.has_secret"),
				),
			},
		},
	})
}

func testAccDataSourceWebhooksConfig(projectKey string) string {
	return fmt.Sprintf(`
data "sonarcloud_webhooks" "test" {
	project = "%s"
}
`, projectKey)
}
