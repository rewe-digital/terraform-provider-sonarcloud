package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-sonarcloud/sonarcloud"
)

// Format examples and generate documentation
//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	providerserver.Serve(context.Background(), sonarcloud.New, providerserver.ServeOpts{
		Address: "sonarcloud",
	})
}
