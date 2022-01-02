package sonarcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/user_groups"
	"math/big"
)

type resourceUserGroupType struct{}

func (r resourceUserGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource manages user groups for the organization.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "Name of the user group",
			},
			"description": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Description for the user group",
			},
			"default": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "Whether the group is the default group.",
			},
			"members_count": {
				Type:        types.NumberType,
				Computed:    true,
				Description: "The number of members this group has.",
			},
		},
	}, nil
}

func (r resourceUserGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceUserGroup{
		p: *(p.(*provider)),
	}, nil
}

type resourceUserGroup struct {
	p provider
}

func (r resourceUserGroup) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Group
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := user_groups.CreateRequest{
		Name:         plan.Name.Value,
		Description:  plan.Description.Value,
		Organization: r.p.organization,
	}

	res, err := r.p.client.UserGroups.Create(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create the user_group",
			fmt.Sprintf("The Create request returned an error: %+v", err),
		)
		return
	}

	var result = Group{
		Default:      types.Bool{Value: res.Group.Default},
		Description:  types.String{Value: res.Group.Description},
		ID:           types.String{Value: big.NewFloat(res.Group.Id).String()},
		MembersCount: types.Number{Value: big.NewFloat(res.Group.MembersCount)},
		Name:         types.String{Value: res.Group.Name},
	}
	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)
}

func (r resourceUserGroup) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Retrieve values from state
	var state Group
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := user_groups.SearchRequest{
		Q: state.Name.Value,
	}

	response, err := r.p.client.UserGroups.SearchAll(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the user_group",
			fmt.Sprintf("The SearchAll request returned an error: %+v", err),
		)
		return
	}

	// Check if the resource exists the list of retrieved resources
	if result, ok := findGroup(response, state.Name.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r resourceUserGroup) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Retrieve values from state
	var state Group
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan Group
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	changed := changedAttrs(req, diags)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	// Note: we skip values that have not been changed
	request := user_groups.UpdateRequest{
		Id: state.ID.Value,
	}

	if _, ok := changed["name"]; ok {
		request.Name = plan.Name.Value
	}
	if _, ok := changed["description"]; ok {
		request.Description = plan.Description.Value
	}

	err := r.p.client.UserGroups.Update(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not update the user_group",
			fmt.Sprintf("The Update request returned an error: %+v", err),
		)
		return
	}

	// We don't have a return value, so we have to query it again
	// Fill in api action struct
	searchRequest := user_groups.SearchRequest{}

	response, err := r.p.client.UserGroups.SearchAll(searchRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the user_group",
			fmt.Sprintf("The SearchAll request returned an error: %+v", err),
		)
		return
	}

	// Check if the resource exists the list of retrieved resources
	if result, ok := findGroup(response, plan.Name.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	}
}

func (r resourceUserGroup) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// Retrieve values from state
	var state Group
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := user_groups.DeleteRequest{
		Id: state.ID.Value,
	}

	err := r.p.client.UserGroups.Delete(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete the user_group",
			fmt.Sprintf("The SearchAll request returned an error: %+v", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceUserGroup) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

func findGroup(response *user_groups.SearchResponseAll, name string) (Group, bool) {
	var result Group
	ok := false
	for _, g := range response.Groups {
		if g.Name == name {
			result = Group{
				ID:           types.String{Value: big.NewFloat(g.Id).String()},
				Default:      types.Bool{Value: g.Default},
				Description:  types.String{Value: g.Description},
				MembersCount: types.Number{Value: big.NewFloat(g.MembersCount)},
				Name:         types.String{Value: g.Name},
			}
			ok = true
			break
		}
	}
	return result, ok
}
