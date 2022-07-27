package sonarcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceQualityGates(t *testing.T) {
	numberOfDefaultQualityGates := "2" //Extra gate is the 'Sonar Way' default quality gate.

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceQualityGatesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sonarcloud_quality_gates.test_quality_gates", "quality_gates.#", numberOfDefaultQualityGates),
				),
			},
		},
	})
}

func testAccDataSourceQualityGatesConfig() string {
	return fmt.Sprintf(`
data "sonarcloud_quality_gates" "test_quality_gates" {}
`)
}
