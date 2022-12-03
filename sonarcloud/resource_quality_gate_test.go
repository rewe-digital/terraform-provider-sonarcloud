package sonarcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceQualityGate(t *testing.T) {
	names := []string{"quality_gate_a", "quality_gate_b"}
	def := []string{"true", "false"}
	metrics := []string{"coverage", "duplicated_lines_density"}
	testError := []string{"10", "11"}
	Op := []string{"LT", "GT"}

	// TODO: use fixed test organization so that changes can be verified.

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQualityGateConfig(names[0], def[0], metrics[0], testError[0], Op[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_quality_gate.test", "name", names[0]),
					resource.TestCheckResourceAttr("sonarcloud_quality_gate.test", "conditions.0.metric", metrics[0]),
					resource.TestCheckResourceAttr("sonarcloud_quality_gate.test", "conditions.0.error", testError[0]),
					resource.TestCheckResourceAttr("sonarcloud_quality_gate.test", "conditions.0.op", Op[0]),
				),
			},
			qualityGateImportCheck("sonarcloud_quality_gate.test", names[0]),
			{
				Config: testAccQualityGateConfig(names[1], def[1], metrics[1], testError[1], Op[1]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_quality_gate.test", "name", names[1]),
					resource.TestCheckResourceAttr("sonarcloud_quality_gate.test", "conditions.0.metric", metrics[1]),
					resource.TestCheckResourceAttr("sonarcloud_quality_gate.test", "conditions.0.error", testError[1]),
					resource.TestCheckResourceAttr("sonarcloud_quality_gate.test", "conditions.0.op", Op[1]),
				),
			},
			qualityGateImportCheck("sonarcloud_quality_gate.test", names[1]),
		},
		CheckDestroy: testAccQualityGateDestroy,
	})
}

func testAccQualityGateDestroy(s *terraform.State) error {
	return nil
}

func testAccQualityGateConfig(name, def, metric, err, op string) string {
	return fmt.Sprintf(`
resource "sonarcloud_quality_gate" "test" {
	name = "%s"
	is_default = "%s"
	conditions = [
		{
			metric = "%s"
			error = "%s"
			op = "%s"
		}
	]
}
	`, name, def, metric, err, op)

}

func qualityGateImportCheck(resourceName, name string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateId:     name,
		ImportStateVerify: true,
	}
}
