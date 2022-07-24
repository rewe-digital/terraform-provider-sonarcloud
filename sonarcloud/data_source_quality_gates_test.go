package sonarcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceQualityGates(t *testing.T) {
	defaultQualityGateName := "TEST_QUALITY_GATE"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceQualityGatesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sonarcloud_quality_gates.test_quality_gates", "name", defaultQualityGateName),
				),
			},
		},
	})
}

func TestAccDataSourceQualityGatesConfig() string {
	return fmt.Sprintf(`
data "sonarcloud_quality_gates "test_quality_gates" {
	name = "TEST_QUALITY_GATE"
}
`)
}
