package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func testAccPreCheckUserToken(t *testing.T) {
	if v := os.Getenv("SONARCLOUD_TOKEN_TEST_USER_LOGIN"); v == "" {
		t.Fatal("SONARCLOUD_TOKEN_TEST_USER_LOGIN must be set for acceptance tests")
	}
}

func TestAccUserToken(t *testing.T) {
	login := os.Getenv("SONARCLOUD_TOKEN_TEST_USER_LOGIN")
	name := "TEST TOKEN"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckUserToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserTokenConfig(login, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_token.test_token", "login", login),
					resource.TestCheckResourceAttr("sonarcloud_user_token.test_token", "name", name),
					resource.TestCheckResourceAttrSet("sonarcloud_user_token.test_token", "token"),
				),
			},
		},
		CheckDestroy: testAccUserTokenDestroy,
	})
}

func testAccUserTokenDestroy(s *terraform.State) error {
	return nil
}

func testAccUserTokenConfig(login string, name string) string {
	return fmt.Sprintf(`
resource "sonarcloud_user_token" "test_token" {
	login = "%s"
	name  = "%s"
}
`, login, name)
}
