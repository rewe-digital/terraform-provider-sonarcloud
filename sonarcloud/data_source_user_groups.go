package sonarcloud

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/user_groups"
)

type dataSourceUserGroupsType struct{}

func (d dataSourceUserGroupsType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This data source retrieves a list of user groups for the configured organization.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"groups": {
				Computed:    true,
				Description: "The groups of this organization.",
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
					"members_count": {
						Type:        types.Float64Type,
						Computed:    true,
						Description: "The number of members in this user group.",
					},
					"default": {
						Type:        types.BoolType,
						Computed:    true,
						Description: "Whether new members are added to this user group per default or not.",
					},
				}),
			},
		},
	}, nil
}

func (d dataSourceUserGroupsType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceUserGroups{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceUserGroups struct {
	p provider
}

func (d dataSourceUserGroups) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var diags diag.Diagnostics

	request := user_groups.SearchRequest{}

	res, err := d.p.client.UserGroups.SearchAll(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read user_groups",
			fmt.Sprintf("The SearchAll request returned an error: %+v", err),
		)
		return
	}

	result := Groups{}
	allGroups := make([]Group, len(res.Groups))
	for i, group := range res.Groups {
		allGroups[i] = Group{
			ID:           types.String{Value: big.NewFloat(group.Id).String()},
			Default:      types.Bool{Value: group.Default},
			Description:  types.String{Value: group.Description},
			MembersCount: types.Number{Value: big.NewFloat(group.MembersCount)},
			Name:         types.String{Value: group.Name},
		}
	}
	result.Groups = allGroups
	result.ID = types.String{Value: d.p.organization}

	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)
}
