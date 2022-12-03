package sonarcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/project_links"
	"strings"
)

type resourceProjectLinkType struct{}

func (r resourceProjectLinkType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource represents a project link.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Computed:    true,
				Description: "ID of the link.",
			},
			"project_key": {
				Type:        types.StringType,
				Required:    true,
				Description: "The key of the project to add the link to.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name the link.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"url": {
				Type:        types.StringType,
				Required:    true,
				Description: "The url of the link.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r resourceProjectLinkType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceProjectLink{
		p: *(p.(*provider)),
	}, nil
}

type resourceProjectLink struct {
	p provider
}

func (r resourceProjectLink) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan ProjectLink
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := project_links.CreateRequest{
		Name:       plan.Name.Value,
		ProjectKey: plan.ProjectKey.Value,
		Url:        plan.Url.Value,
	}

	res, err := r.p.client.ProjectLinks.Create(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create the project link",
			fmt.Sprintf("The Create request returned an error: %+v", err),
		)
		return
	}

	link := res.Link
	var result = ProjectLink{
		ID:         types.String{Value: link.Id},
		ProjectKey: plan.ProjectKey,
		Name:       types.String{Value: link.Name},
		Url:        types.String{Value: link.Url},
	}
	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)
}

func (r resourceProjectLink) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Retrieve values from state
	var state ProjectLink
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := project_links.SearchRequest{
		ProjectKey: state.ProjectKey.Value,
	}

	response, err := r.p.client.ProjectLinks.Search(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the project link",
			fmt.Sprintf("The SearchAll request returned an error: %+v", err),
		)
		return
	}

	// Check if the resource exists the list of retrieved resources
	if result, ok := findProjectLink(response, state.ID.Value, state.ProjectKey.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r resourceProjectLink) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// NOOP, we always need to recreate
}

func (r resourceProjectLink) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state ProjectLink
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := project_links.DeleteRequest{
		Id: state.ID.Value,
	}
	err := r.p.client.ProjectLinks.Delete(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete the project link",
			fmt.Sprintf("The Delete request returned an error: %+v", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceProjectLink) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id,project_key. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_key"), idParts[1])...)
}

// findProjectLink returns the link with the given id, if it exists in the response
func findProjectLink(response *project_links.SearchResponse, id, project_key string) (ProjectLink, bool) {
	var result ProjectLink
	ok := false
	for _, link := range response.Links {
		if link.Id == id {
			result = ProjectLink{
				ID:         types.String{Value: link.Id},
				ProjectKey: types.String{Value: project_key},
				Name:       types.String{Value: link.Name},
				Url:        types.String{Value: link.Url},
			}
			ok = true
			break
		}
	}
	return result, ok
}
