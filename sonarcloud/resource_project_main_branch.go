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
		Description: "This resource manages a project main branch.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name of the project main branch.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
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
			},
		},
	}, nil
}

func (r resourceProjectMainBranchType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceProject{
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

	projectKey := state.ProjectKey.Value
	if _, ok := changed["project_key"]; ok {
		projectKey = plan.ProjectKey.Value
	}

	name := state.Name.Value
	if _, ok := changed["name"]; ok {
		name = plan.Name.Value
	}

	request := project_branches.RenameRequest{
		Project: projectKey,
		Name:    name,
	}

	err := r.p.client.ProjectBranches.Rename(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not update the main project branch",
			fmt.Sprintf("The Rename request returned an error: %+v", err),
		)
		return
	}
}

func (r resourceProjectMainBranch) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state ProjectMainBranch
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: according to docs, this may not work for main branches, and may require admin privilege
	// https://github.com/reinoudk/go-sonarcloud/blob/main/sonarcloud/project_branches/project_branches_gen.go#L5
	request := project_branches.DeleteRequest{
		Project: state.ProjectKey.Value,
		Branch:  state.Name.Value,
	}

	err := r.p.client.ProjectBranches.Delete(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete the project main branch",
			fmt.Sprintf("The Delete request returned an error: %+v", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceProjectMainBranch) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("project_key"), req, resp)
}
