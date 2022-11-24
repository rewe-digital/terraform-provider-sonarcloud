package sonarcloud

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured   bool
	client       *sonarcloud.Client
	organization string
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"organization": {
				Type:     types.StringType,
				Optional: true,
				Description: "The SonarCloud organization to manage the resources for. This value must be set in the" +
					" `SONARCLOUD_ORGANIZATION` environment variable if left empty.",
			},
			"token": {
				Type:      types.StringType,
				Optional:  true,
				Sensitive: true,
				Description: "The token of a user with admin permissions in the organization. This value must be set in" +
					" the `SONARCLOUD_TOKEN` environment variable if left empty.",
			},
		},
	}, nil
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var organization string
	if config.Organization.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as organization",
		)
		return
	}

	if config.Organization.Null {
		organization = os.Getenv("SONARCLOUD_ORGANIZATION")
	} else {
		organization = config.Organization.Value
	}

	var token string
	if config.Token.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as token",
		)
	}

	if config.Token.Null {
		token = os.Getenv("SONARCLOUD_TOKEN")
	} else {
		token = config.Token.Value
	}

	c := sonarcloud.NewClient(organization, token, nil)
	p.client = c
	p.organization = organization
	p.configured = true
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"sonarcloud_user_group":             resourceUserGroupType{},
		"sonarcloud_user_group_member":      resourceUserGroupMemberType{},
		"sonarcloud_project":                resourceProjectType{},
		"sonarcloud_project_main_branch":    resourceProjectMainBranchType{},
		"sonarcloud_user_token":             resourceUserTokenType{},
		"sonarcloud_quality_gate":           resourceQualityGateType{},
		"sonarcloud_quality_gate_selection": resourceQualityGateSelectionType{},
		"sonarcloud_user_group_permissions": resourceUserGroupPermissionsType{},
	}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"sonarcloud_projects":           dataSourceProjectsType{},
		"sonarcloud_user_group":         dataSourceUserGroupType{},
		"sonarcloud_user_groups":        dataSourceUserGroupsType{},
		"sonarcloud_user_group_members": dataSourceUserGroupMembersType{},
		"sonarcloud_quality_gate":       dataSourceQualityGateType{},
		"sonarcloud_quality_gates":      dataSourceQualityGatesType{},
	}, nil
}

type providerData struct {
	Organization types.String `tfsdk:"organization"`
	Token        types.String `tfsdk:"token"`
}
