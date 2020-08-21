package sonarcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider SonarCloud
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap:   map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}