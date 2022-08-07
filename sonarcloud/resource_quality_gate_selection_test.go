package sonarcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccPreCheckQualityGateSelection(t *testing.T) {
	if v := os.Getenv("SONARCLOUD_QUALITY_GATE_ID"); v == "" {
		t.Fatal("SONARCLOUD_QUALITY_GATE_ID must be set for acceptance tests")
	}
	if v := os.Getenv("SONARCLOUD_PROJECT_KEY"); v == "" {
		t.Fatal("SONARCLOUD_PROJECT_KEY must be set for acceptance tests")
	}
}

func TestAccResourceQualityGateSelection(t *testing.T) {
	gate_id := os.Getenv("SONARCLOUD_QUALITY_GATE_ID")
	project_key := os.Getenv("SONARCLOUD_PROJECT_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); testAccPreCheckQualityGateSelection(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQualityGateSelectionConfig(gate_id, project_key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarcloud_quality_gate_selection.test", "gate_id", gate_id),
					resource.TestCheckResourceAttr("sonarcloud_quality_gate_selection.test", "project_keys.0", project_key),
				),
			},
		},
		CheckDestroy: testAccQualityGateSelectionDestroy,
	})
}

func testAccQualityGateSelectionDestroy(s *terraform.State) error {
	return nil
}

func testAccQualityGateSelectionConfig(gateId, projectKey string) string {
	return fmt.Sprintf(`
resource "sonarcloud_quality_gate_selection" "test" {
	gate_id = "%s"
	project_keys = ["%s"]
}
	`, gateId, projectKey)
}
