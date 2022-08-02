package sonarcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testDataAccPreCheckQualityGate(t *testing.T) {
	if v := os.Getenv("SONARCLOUD_QUALITY_GATE_NAME"); v == "" {
		t.Fatal("SONARCLOUD_QUALITY_GATE_NAME must be set for acceptance tests")
	}
}

func TestAccDataSourceQualityGate(t *testing.T) {
	qualityGateName := os.Getenv("SONARCLOUD_QUALITY_GATE_NAME")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); testDataAccPreCheckQualityGate(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceQualityGateConfig(qualityGateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sonarcloud_quality_gate.test_quality_gate", "name", qualityGateName),
				),
			},
		},
	})
}

func testAccDataSourceQualityGateConfig(qualityGateName string) string {
	return fmt.Sprintf(`
data "sonarcloud_quality_gate" "test_quality_gate" {
	name = "%s"
}
`, qualityGateName)
}
