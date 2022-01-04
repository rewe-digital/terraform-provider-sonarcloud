package sonarcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/user_groups"
)

type dataSourceUserGroupType struct{}

func (d dataSourceUserGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This data source retrieves a single user group.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Computed:    true,
				Description: "The ID of the user group.",
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name of the user group.",
			},
			"description": {
				Type:        types.StringType,
				Computed:    true,
				Description: "The description of the user group.",
			},
			"members_count": {
				Type:        types.NumberType,
				Computed:    true,
				Description: "The number of members in this user group.",
			},
			"default": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "Whether new members are added to this user group per default or not.",
			},
		},
	}, nil
}

func (d dataSourceUserGroupType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceUserGroup{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceUserGroup struct {
	p provider
}

func (d dataSourceUserGroup) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	// Retrieve values from config
	var config Group
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := user_groups.SearchRequest{
		Q: config.Name.Value,
	}

	response, err := d.p.client.UserGroups.SearchAll(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the user_group",
			fmt.Sprintf("The SearchAll request returned an error: %+v", err),
		)
		return
	}

	// Check if the resource exists the list of retrieved resources
	if result, ok := findGroup(response, config.Name.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}
