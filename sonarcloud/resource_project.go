package sonarcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/projects"
)

type resourceProjectType struct{}

func (r resourceProjectType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource manages a project.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name of the project. **Warning:** forces project recreation when changed.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Validators: []tfsdk.AttributeValidator{
					stringLengthBetween(1, 255),
				},
			},
			"key": {
				Type:        types.StringType,
				Required:    true,
				Description: "The key of the project. **Warning**: must be globally unique.",
				Validators: []tfsdk.AttributeValidator{
					stringLengthBetween(1, 400),
				},
			},
			"visibility": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				Description: "The visibility of the project. Use `private` to only share it with your organization." +
					" Use `public` if the project should be visible to everyone. Defaults to the organization's default visibility." +
					" **Note:** private projects are only available when you have a SonarCloud subscription.",
				Validators: []tfsdk.AttributeValidator{
					allowedOptions("public", "private"),
				},
			},
		},
	}, nil
}

func (r resourceProjectType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceProject{
		p: *(p.(*provider)),
	}, nil
}

type resourceProject struct {
	p provider
}

func (r resourceProject) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Project
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := projects.CreateRequest{
		Name:         plan.Name.Value,
		Organization: r.p.organization,
		Project:      plan.Key.Value,
		Visibility:   plan.Visibility.Value,
	}

	res, err := r.p.client.Projects.Create(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create the project",
			fmt.Sprintf("The Create request returned an error: %+v", err),
		)
		return
	}

	var result = Project{
		ID:         types.String{Value: res.Project.Key},
		Name:       types.String{Value: res.Project.Name},
		Key:        types.String{Value: res.Project.Key},
		Visibility: types.String{Value: plan.Visibility.Value},
	}
	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)
}

func (r resourceProject) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Retrieve values from state
	var state Project
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := projects.SearchRequest{
		Projects: state.Key.Value,
	}

	response, err := r.p.client.Projects.SearchAll(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the project",
			fmt.Sprintf("The SearchAll request returned an error: %+v", err),
		)
		return
	}

	// Check if the resource exists the list of retrieved resources
	if result, ok := findProject(response, state.Key.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r resourceProject) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Retrieve values from state
	var state Project
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan Project
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	changed := changedAttrs(req, diags)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, ok := changed["key"]; ok {
		request := projects.UpdateKeyRequest{
			From: state.Key.Value,
			To:   plan.Key.Value,
		}

		err := r.p.client.Projects.UpdateKey(request)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not update the project key",
				fmt.Sprintf("The UpdateKey request returned an error: %+v", err),
			)
			return
		}
	}

	if _, ok := changed["visibility"]; ok {
		request := projects.UpdateVisibilityRequest{
			Project:    plan.Key.Value,
			Visibility: plan.Visibility.Value,
		}

		err := r.p.client.Projects.UpdateVisibility(request)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not update the project visibility",
				fmt.Sprintf("The UpdateVisibility request returned an error: %+v", err),
			)
			return
		}
	}

	// We don't have a return value, so we have to query it again
	// Fill in api action struct
	searchRequest := projects.SearchRequest{}

	response, err := r.p.client.Projects.SearchAll(searchRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the project",
			fmt.Sprintf("The SearchAll request returned an error: %+v", err),
		)
		return
	}

	// Check if the resource exists the list of retrieved resources
	if result, ok := findProject(response, plan.Key.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	}
}

func (r resourceProject) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// Retrieve values from state
	var state Project
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := projects.DeleteRequest{
		Project: state.Key.Value,
	}

	err := r.p.client.Projects.Delete(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete the project",
			fmt.Sprintf("The Delete request returned an error: %+v", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceProject) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("key"), req, resp)
}
