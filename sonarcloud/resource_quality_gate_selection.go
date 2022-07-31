package sonarcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/qualitygates"
)

type resourceQualityGateSelectionType struct{}

func (r resourceQualityGateSelectionType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource selects a quality gate for a project",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.Float64Type,
				Description: "The ID of the selection.",
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"gate_id": {
				Type:        types.StringType,
				Description: "The ID of the quality gate that is selecting the project(s).",
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"project_key": {
				Type:        types.StringType,
				Description: "The Keys of the projects which have been selected on the referenced quality gate",
				Required:    true,
			},
		},
	}, nil
}

func (r resourceQualityGateSelectionType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceQualityGateSelection{
		p: *(p.(*provider)),
	}, nil
}

type resourceQualityGateSelection struct {
	p provider
}

func (r resourceQualityGateSelection) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unkown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Selection
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct for Quality Gates
	request := qualitygates.SelectRequest{
		GateId:       plan.GateId.Value,
		ProjectKey:   plan.ProjectKey.Value,
		Organization: r.p.organization,
	}
	err := r.p.client.Qualitygates.Select(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create Quality Gate Selection",
			fmt.Sprintf("The Select request returned an error: %+v", err),
		)
		return
	}

	// Query for selection
	searchRequest := qualitygates.SearchRequest{
		GateId:       plan.GateId.Value,
		Organization: r.p.organization,
	}

	res, err := r.p.client.Qualitygates.Search(searchRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read Quality Gate Selection",
			fmt.Sprintf("The Search request reutrned an error: %+v", err),
		)
		return
	}

	var result Selection
	for _, project := range res.Results {
		if project.Key == plan.ProjectKey.Value {
			result = Selection{
				ID:         types.Float64{Value: project.Id},
				GateId:     types.String{Value: plan.GateId.Value},
				ProjectKey: types.String{Value: project.Key},
			}
			break
		}
	}
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceQualityGateSelection) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Selection
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := qualitygates.DeselectRequest{
		Organization: r.p.organization,
		ProjectKey:   state.ProjectKey.Value,
	}
	err := r.p.client.Qualitygates.Deselect(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Deselect the Quality Gate Selection",
			fmt.Sprintf("The Deselect request returned an error: %+v", err),
		)
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceQualityGateSelection) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state Selection
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	searchRequest := qualitygates.SearchRequest{
		GateId:       state.GateId.Value,
		Organization: r.p.organization,
	}
	res, err := r.p.client.Qualitygates.Search(searchRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Read the Quality Gate Selection",
			fmt.Sprintf("The Search request returned an error: %+v", err),
		)
		return
	}
	found := false
	for _, searchR := range res.Results {
		if searchR.Id == state.ID.Value {
			result := Selection{
				ID:         types.Float64{Value: searchR.Id},
				GateId:     types.String{Value: state.GateId.Value},
				ProjectKey: types.String{Value: searchR.Key},
			}
			diags = resp.State.Set(ctx, result)
			resp.Diagnostics.Append(diags...)
			found = true
		}
	}
	if !found {
		resp.State.RemoveResource(ctx)
	}
}

func (r resourceQualityGateSelection) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var state Selection
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan Selection
	diags = req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	changed := changedAttrs(req, diags)

	if _, ok := changed["project_key"]; ok {
		deselectRequest := qualitygates.DeselectRequest{
			Organization: r.p.organization,
			ProjectKey:   state.ProjectKey.Value,
		}
		err := r.p.client.Qualitygates.Deselect(deselectRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not Deselect the Quality Gate selection",
				fmt.Sprintf("The Deselect request returned an error: %+v", err),
			)
			return
		}
		selectRequest := qualitygates.SelectRequest{
			GateId:       state.GateId.Value,
			Organization: r.p.organization,
			ProjectKey:   plan.ProjectKey.Value,
		}
		err = r.p.client.Qualitygates.Select(selectRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not Select the Quality Gate selection",
				fmt.Sprintf("The Select request returned an error: %+v", err),
			)
		}
	}

	request := qualitygates.SearchRequest{
		GateId:       plan.GateId.Value,
		Organization: r.p.organization,
	}
	res, err := r.p.client.Qualitygates.Search(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Read the Quality Gate Selection",
			fmt.Sprintf("The Search request returned an error: %+v", err),
		)
		return
	}
	var result Selection
	for _, searchR := range res.Results {
		if searchR.Id == state.ID.Value {
			result = Selection{
				ID:         types.Float64{Value: searchR.Id},
				GateId:     types.String{Value: plan.GateId.Value},
				ProjectKey: types.String{Value: searchR.Key},
			}
		}
	}
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}
