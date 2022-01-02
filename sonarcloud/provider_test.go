package sonarcloud

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"testing"
)

var testAccProviderFactories map[string]func() (tfprotov6.ProviderServer, error)

func init() {
	testAccProviderFactories = map[string]func() (tfprotov6.ProviderServer, error) {
		"sonarcloud": func() (tfprotov6.ProviderServer, error) {
			return tfsdk.NewProtocol6Server(New()), nil
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
