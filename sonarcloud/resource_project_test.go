package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccResourceProject(t *testing.T) {
	prefix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	names := []string{"project_a", "project_b"}
	keys := []string{prefix + "sonarcloud-provider-acc-test_a", prefix + "sonarcloud-provider-acc-test_b" + prefix}
	// TODO: use private-enabled organization for acceptance tests so we can verify visibility changes
	visibilities := []string{"public"}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectConfig(names[0], keys[0], visibilities[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_project.test", "name", names[0]),
					resource.TestCheckResourceAttr("sonarcloud_project.test", "key", keys[0]),
					resource.TestCheckResourceAttr("sonarcloud_project.test", "visibility", visibilities[0]),
				),
			},
			{
				Config: testAccProjectConfig(names[1], keys[0], visibilities[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_project.test", "name", names[1]),
					resource.TestCheckResourceAttr("sonarcloud_project.test", "key", keys[0]),
					resource.TestCheckResourceAttr("sonarcloud_project.test", "visibility", visibilities[0]),
				),
			},
			{
				Config: testAccProjectConfig(names[1], keys[1], visibilities[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_project.test", "name", names[1]),
					resource.TestCheckResourceAttr("sonarcloud_project.test", "key", keys[1]),
					resource.TestCheckResourceAttr("sonarcloud_project.test", "visibility", visibilities[0]),
				),
			},
		},
		CheckDestroy: testAccProjectDestroy,
	})
}

func testAccProjectDestroy(s *terraform.State) error {
	return nil
}

func testAccProjectConfig(name, key, visibility string) string {
	return fmt.Sprintf(`
resource "sonarcloud_project" "test" {
	name = "%s"
	key = "%s"
	visibility = "%s"
}
`, name, key, visibility)
}
