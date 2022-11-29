package sonarcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud"
)

type dataSourceUserGroupPermissionsType struct{}

func (d dataSourceUserGroupPermissionsType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This data source retrieves all the user groups and their permissions.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Computed:    true,
				Description: "The implicit ID of the data source.",
			},
			"project_key": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The key of the project to read the user group permissions for.",
			},
			"groups": {
				Computed:    true,
				Description: "The groups and their permissions.",
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The ID of the user group.",
					},
					"name": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The name of the user group.",
					},
					"description": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The description of the user group.",
					},
					"permissions": {
						Type:        types.SetType{ElemType: types.StringType},
						Computed:    true,
						Description: "The permissions of this user group.",
					},
				}),
			},
		},
	}, nil
}

func (d dataSourceUserGroupPermissionsType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceUserGroupPermissions{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceUserGroupPermissions struct {
	p provider
}

func (d dataSourceUserGroupPermissions) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var config DataUserGroupPermissions
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Query for permissions
	searchRequest := UserGroupPermissionsSearchRequest{ProjectKey: config.ProjectKey.Value}
	groups, err := sonarcloud.GetAll[UserGroupPermissionsSearchRequest, UserGroupPermissionsSearchResponseGroup](d.p.client, "/permissions/groups", searchRequest, "groups")
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get user group permissions",
			fmt.Sprintf("The request returned an error: %+v", err),
		)
		return
	}

	result := DataUserGroupPermissions{}
	var allGroups []DataUserGroupPermissionsGroup
	for _, group := range groups {
		permissionsElems := make([]attr.Value, len(group.Permissions))
		for i, permission := range group.Permissions {
			permissionsElems[i] = types.String{Value: permission}
		}

		allGroups = append(allGroups, DataUserGroupPermissionsGroup{
			ID:          types.String{Value: group.Id},
			Name:        types.String{Value: group.Name},
			Description: types.String{Value: group.Description},
			Permissions: types.Set{Elems: permissionsElems, ElemType: types.StringType},
		})
	}
	result.Groups = allGroups
	result.ID = types.String{Value: d.p.organization}
	result.ProjectKey = config.ProjectKey

	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)

}
