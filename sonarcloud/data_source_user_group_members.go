package sonarcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/user_groups"
)

type dataSourceUserGroupMembersType struct{}

func (d dataSourceUserGroupMembersType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This data source retrieves a list of users for the given group.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"group": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name of the group.",
			},
			"users": {
				Computed:    true,
				Description: "The users of the group.",
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"login": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The login of this user",
					},
					"name": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The name of this user",
					},
				}),
			},
		},
	}, nil
}

func (d dataSourceUserGroupMembersType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceUserGroupMembers{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceUserGroupMembers struct {
	p provider
}

func (d dataSourceUserGroupMembers) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var config Users
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// An empty search request retrieves all members
	request := user_groups.UsersRequest{
		Name: config.Group.Value,
	}

	res, err := d.p.client.UserGroups.UsersAll(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read user_group_members.",
			fmt.Sprintf("The UsersAll request returned an error: %+v", err),
		)
		return
	}

	result := Users{}
	allUsers := make([]User, len(res.Users))
	for i, user := range res.Users {
		allUsers[i] = User{
			Login: types.String{Value: user.Login},
			Name:  types.String{Value: user.Name},
		}
	}
	result.Users = allUsers
	result.ID = config.Group
	result.Group = config.Group

	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)
}
