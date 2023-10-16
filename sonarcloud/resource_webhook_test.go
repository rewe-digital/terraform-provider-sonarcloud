package sonarcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccWebhook(t *testing.T) {
	project := os.Getenv("SONARCLOUD_PROJECT_KEY")
	secret := "ThisIsNotAVeryGoodSecret..."

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectWebhookConfig("test", project, secret, "https://www.example.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "name", "test"),
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "url", "https://www.example.com"),
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "secret", secret),
				),
			},
			webhookImportCheck("sonarcloud_webhook.test", project),
			{
				Config: testAccProjectWebhookConfig("test-two", project, "", "https://www.example.com/test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "name", "test-two"),
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "url", "https://www.example.com/test"),
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "secret", ""),
				),
			},
			webhookImportCheck("sonarcloud_webhook.test", project),
			{
				Config: testAccOrganizationWebhookConfig("test", secret, "https://www.example.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "name", "test"),
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "url", "https://www.example.com"),
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "secret", secret),
				),
			},
			webhookImportCheck("sonarcloud_webhook.test", ""),
			{
				Config: testAccOrganizationWebhookConfig("test-two", "", "https://www.example.com/test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "name", "test-two"),
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "url", "https://www.example.com/test"),
					resource.TestCheckResourceAttr("sonarcloud_webhook.test", "secret", ""),
				),
			},
			webhookImportCheck("sonarcloud_webhook.test", ""),
		},
		CheckDestroy: testAccWebhookDestroy,
	})
}

func testAccWebhookDestroy(s *terraform.State) error {
	return nil
}

func testAccProjectWebhookConfig(name, project, secret, url string) string {
	result := fmt.Sprintf(`
resource "sonarcloud_webhook" "test" {
  name    = "%s"
  project = "%s"
  secret  = "%s"
  url     = "%s"
}
`, name, project, secret, url)
	return result
}

func testAccOrganizationWebhookConfig(name, secret, url string) string {
	result := fmt.Sprintf(`
resource "sonarcloud_webhook" "test" {
  name    = "%s"
  secret  = "%s"
  url     = "%s"
}
`, name, secret, url)
	return result
}

func webhookImportCheck(resourceName, project string) resource.TestStep {
	return resource.TestStep{
		ResourceName: resourceName,
		ImportState:  true,
		ImportStateIdFunc: func(state *terraform.State) (string, error) {
			id := state.RootModule().Resources[resourceName].Primary.ID
			return fmt.Sprintf("%s,%s", id, project), nil
		},
		// We need to set ImportStateVerify to false because we cannot read the secret value from the API.
		ImportStateVerify: false,
	}
}
