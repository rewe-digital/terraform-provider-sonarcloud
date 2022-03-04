package sonarcloud

import (
	"context"
    "strings"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/user_groups"
)

type resourceUserGroupMemberType struct{}

func (r resourceUserGroupMemberType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource manages a single member of a user group.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"group": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The name of the group to which the user should be added.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"login": {
				Type:        types.StringType,
				Required:    true,
				Description: "The login of the user that should be added to the group.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r resourceUserGroupMemberType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceUserGroupMember{
		p: *(p.(*provider)),
	}, nil
}

type resourceUserGroupMember struct {
	p provider
}

func (r resourceUserGroupMember) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan GroupMember
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := user_groups.AddUserRequest{
		Login:        plan.Login.Value,
		Name:         plan.Group.Value,
		Organization: r.p.organization,
	}

	err := r.p.client.UserGroups.AddUser(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create the user_group_member.",
			fmt.Sprintf("The AddUser request returned an error: %+v", err),
		)
		return
	}

	// We have no response, assume the values were set when no error has been returned and just set ID
	state := plan
	state.ID = types.String{Value: fmt.Sprintf("%s%s", plan.Group.Value, plan.Login.Value)}
	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
}

func (r resourceUserGroupMember) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Retrieve values from state
    var state GroupMember
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
    
	// Fill in api action struct
	request := user_groups.UsersRequest{
		Q: state.Login.Value,
        Name: state.Group.Value,
	}

	response, err := r.p.client.UserGroups.UsersAll(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the user_group_member.",
            fmt.Sprintf("The UsersAll request returned an error: %+v", err),
		)
		return
	}
    
	// Check if the resource exists the list of retrieved resources
	if result, ok := findGroupMember(response, state.Group.Value, state.Login.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r resourceUserGroupMember) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// NOOP, we always need to recreate
}

func (r resourceUserGroupMember) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// Retrieve values from state
	var state GroupMember
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := user_groups.RemoveUserRequest{
		Login:        state.Login.Value,
		Name:         state.Group.Value,
		Organization: r.p.organization,
	}

	err := r.p.client.UserGroups.RemoveUser(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete the user_group_member.",
			fmt.Sprintf("The RemoveUser request returned an error: %+v", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceUserGroupMember) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
    idParts := strings.Split(req.ID,",")
    if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
        resp.Diagnostics.AddError(
            "Unexpected Import Identifier",
            fmt.Sprintf("Expected import identifier with format: login,group. Got: %q", req.ID),
        )
        return
    }

    resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("login"), idParts[0])...)
    resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("group"), idParts[1])...)
}
