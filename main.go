package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-provider-sonarcloud/sonarcloud"
)

func main() {
	tfsdk.Serve(context.Background(), sonarcloud.New, tfsdk.ServeOpts{
		Name: "sonarcloud",
	})
}
