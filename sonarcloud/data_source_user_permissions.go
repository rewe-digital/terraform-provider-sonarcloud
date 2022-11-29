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

type dataSourceUserPermissionsType struct{}

func (d dataSourceUserPermissionsType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This data source retrieves all the users of an organization and their permissions.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Computed:    true,
				Description: "The implicit ID of the data source.",
			},
			"project_key": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The key of the project to read the user permissions for.",
			},
			"users": {
				Computed:    true,
				Description: "The users and their permissions.",
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"login": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The login of the user.",
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.RequiresReplace(),
						},
					},
					"name": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The name of the user.",
					},
					"permissions": {
						Type:        types.SetType{ElemType: types.StringType},
						Required:    true,
						Description: "The permissions of the user.",
					},
					"avatar": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The avatar ID of the user.",
					},
				}),
			},
		},
	}, nil
}

func (d dataSourceUserPermissionsType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceUserPermissions{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceUserPermissions struct {
	p provider
}

func (d dataSourceUserPermissions) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var config DataUserPermissions
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Query for permissions
	searchRequest := UserPermissionsSearchRequest{ProjectKey: config.ProjectKey.Value}
	users, err := sonarcloud.GetAll[UserPermissionsSearchRequest, UserPermissionsSearchResponseUser](d.p.client, "/permissions/users", searchRequest, "users")
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get user permissions",
			fmt.Sprintf("The request returned an error: %+v", err),
		)
		return
	}

	result := DataUserPermissions{}
	var allUsers []DataUserPermissionsUser
	for _, user := range users {
		permissionsElems := make([]attr.Value, len(user.Permissions))
		for i, permission := range user.Permissions {
			permissionsElems[i] = types.String{Value: permission}
		}

		allUsers = append(allUsers, DataUserPermissionsUser{
			Login:       types.String{Value: user.Login},
			Name:        types.String{Value: user.Name},
			Permissions: types.Set{Elems: permissionsElems, ElemType: types.StringType},
			Avatar:      types.String{Value: user.Avatar},
		})
	}
	result.Users = allUsers
	result.ID = types.String{Value: d.p.organization}
	result.ProjectKey = config.ProjectKey

	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)

}
