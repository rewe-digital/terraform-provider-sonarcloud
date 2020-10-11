package sonarcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
