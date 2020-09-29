package sonarcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

var testAccProviderFactories map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"sonarcloud": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SONARCLOUD_ORGANIZATION"); v == "" {
		t.Fatal("SONARCLOUD_ORGANIZATION must be set for acceptance tests")
	}
	if v := os.Getenv("SONARCLOUD_TOKEN"); v == "" {
		t.Fatal("SONARCLOUD_TOKEN must be set for acceptance tests")
	}
}

type testStep struct {
	TerraformConfigFunc func() string
	WantGroups          string
}

type testCase struct {
	Name         string
	TerraformDir string
	CheckDestroy func(s *terraform.State) error
	Steps        []testStep
}

func TestAccDataSourceUserGroups(t *testing.T) {
	testCases := []testCase{
		{
			Name: "Teams",
			Steps: []testStep{
				{TerraformConfigFunc: testAccSonarCloudUserGroupsConfig, WantGroups: "2"},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: testAccProviderFactories,
				Steps:             testSteps(tc.Steps),
			})
		})
	}
}

func testSteps(testSteps []testStep) (steps []resource.TestStep) {
	for _, step := range testSteps {
		steps = append(steps, resource.TestStep{
			Config: step.TerraformConfigFunc(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.sonarcloud_user_groups.test_groups", "groups.#", step.WantGroups),
			),
		})
	}
	return
}

func testAccSonarCloudUserGroupsConfig() string {
	return fmt.Sprintf(`
data "sonarcloud_user_groups" "test_groups" {}
`)
}
