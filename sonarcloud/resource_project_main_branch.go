package sonarcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/project_branches"
)

type resourceProjectMainBranchType struct{}

func (r resourceProjectMainBranchType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `This resource manages a project main branch.

Note that certain operations, such as the deletion of a project's main branch configuration, may
not be permitted by the SonarCloud web API, or may require admin permissions.
		`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name of the project main branch.",
				Validators: []tfsdk.AttributeValidator{
					stringLengthBetween(1, 255),
				},
			},
			"project_key": {
				Type:        types.StringType,
				Required:    true,
				Description: "The key of the project.",
				Validators: []tfsdk.AttributeValidator{
					stringLengthBetween(1, 400),
				},
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r resourceProjectMainBranchType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceProjectMainBranch{
		p: *(p.(*provider)),
	}, nil
}

type resourceProjectMainBranch struct {
	p provider
}

func (r resourceProjectMainBranch) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan ProjectMainBranch
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := project_branches.RenameRequest{
		Project: plan.ProjectKey.Value,
		Name:    plan.Name.Value,
	}

	err := r.p.client.ProjectBranches.Rename(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create the main project branch",
			fmt.Sprintf("The Rename request returned an error: %+v", err),
		)
		return
	}

	var result = ProjectMainBranch{
		ID:         types.String{Value: plan.Name.Value},
		Name:       types.String{Value: plan.Name.Value},
		ProjectKey: types.String{Value: plan.ProjectKey.Value},
	}
	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)
}

func (r resourceProjectMainBranch) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Retrieve values from state
	var state ProjectMainBranch
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := project_branches.ListRequest{
		Project: state.ProjectKey.Value,
	}

	response, err := r.p.client.ProjectBranches.List(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the project branches",
			fmt.Sprintf("The List request returned an error: %+v", err),
		)
		return
	}

	// Check if the main branch matches the declared main branch
	if result, ok := findProjectMainBranch(response, state.Name.Value, state.ProjectKey.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r resourceProjectMainBranch) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Retrieve values from state
	var state ProjectMainBranch
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan ProjectMainBranch
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	changed := changedAttrs(req, diags)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, ok := changed["name"]; !ok {
		resp.Diagnostics.AddError(
			"Name from plan does not differ from state.",
			"This should not be possible and indicates an issue with the provider. Please contact the developers.")
	}

	request := project_branches.RenameRequest{
		Project: plan.ProjectKey.Value,
		Name:    plan.Name.Value,
	}

	err := r.p.client.ProjectBranches.Rename(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not update the main project branch",
			fmt.Sprintf("The Rename request returned an error: %+v", err),
		)
		return
	}

	// In the absence of an error we assume that the main project branch was updated.
	// The alternative would be to query the API again to verify this.
	// (The rename-response does not have a return value.)
	// As the API seems to be eventually consistent, this results in flaky behaviour, so we just keep it simple for now.
	var result = ProjectMainBranch{
		ID:         types.String{Value: plan.Name.Value},
		Name:       types.String{Value: plan.Name.Value},
		ProjectKey: types.String{Value: plan.ProjectKey.Value},
	}
	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)
}

func (r resourceProjectMainBranch) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state ProjectMainBranch
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// NOTE: according to docs, it's not possible to DELETE main branches, at least not without admin privilege
	// Therefore, we simply remove the main branch from state, per recommendation of @reinoudk.
	resp.State.RemoveResource(ctx)
}

func (r resourceProjectMainBranch) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("project_key"), req, resp)
}
