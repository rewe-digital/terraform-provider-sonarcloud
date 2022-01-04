package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func testAccPreCheckUserGroupMember(t *testing.T) {
	if v := os.Getenv("SONARCLOUD_TEST_USER_LOGIN"); v == "" {
		t.Fatal("SONARCLOUD_TEST_USER_LOGIN must be set for acceptance tests")
	}
	if v := os.Getenv("SONARCLOUD_TEST_GROUP_NAME"); v == "" {
		t.Fatal("SONARCLOUD_TEST_GROUP_NAME must be set for acceptance tests")
	}
}

func TestAccUserGroupMember(t *testing.T) {
	login := os.Getenv("SONARCLOUD_TEST_USER_LOGIN")
	group := os.Getenv("SONARCLOUD_TEST_GROUP_NAME")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); testAccPreCheckUserGroupMember(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupMemberConfig(group, login),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_user_group_member.test_group_member", "group", group),
					resource.TestCheckResourceAttr("sonarcloud_user_group_member.test_group_member", "login", login),
				),
			},
		},
		CheckDestroy: testAccUserGroupMemberDestroy,
	})
}

func testAccUserGroupMemberDestroy(s *terraform.State) error {
	return nil
}

func testAccUserGroupMemberConfig(group string, login string) string {
	return fmt.Sprintf(`
resource "sonarcloud_user_group_member" "test_group_member" {
	group = "%s"
	login = "%s"
}
`, group, login)
}
