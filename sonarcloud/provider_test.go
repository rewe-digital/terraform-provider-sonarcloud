package sonarcloud

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProviderFactories map[string]func() (tfprotov6.ProviderServer, error)

func init() {
	testAccProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"sonarcloud": providerserver.NewProtocol6WithError(New()),
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
