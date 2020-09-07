package sonarcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const API = "https://www.sonarcloud.io/api"

// Provider SonarCloud
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORGANIZATION", nil),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TOKEN", nil),
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"sonarcloud_user_groups": dataSourceUserGroups(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	organization := d.Get("organization").(string)
	token := d.Get("token").(string)

	c := &Config{
		Organization: organization,
		Token:        token,
	}

	var diags diag.Diagnostics

	return c, diags
}

type Config struct {
	Organization string
	Token        string
}
