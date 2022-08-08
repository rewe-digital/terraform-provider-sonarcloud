package sonarcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/qualitygates"
)

type resourceQualityGateSelectionType struct{}

func (r resourceQualityGateSelectionType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource selects a quality gate for one or more projects",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Description: "The implicit ID of the resource",
				Computed:    true,
			},
			"gate_id": {
				Type:        types.StringType,
				Description: "The ID of the quality gate that is selected for the project(s).",
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"project_keys": {
				Type:        types.SetType{ElemType: types.StringType},
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

	for _, s := range plan.ProjectKeys.Elems {
		// Fill in api action struct for Quality Gates
		request := qualitygates.SelectRequest{
			GateId:       plan.GateId.Value,
			ProjectKey:   s.(types.String).Value,
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
			fmt.Sprintf("The Search request returned an error: %+v", err),
		)
		return
	}

	if result, ok := findSelection(res, plan.ProjectKeys.Elems); ok {
		result.GateId = types.String{Value: plan.GateId.Value}
		result.ID = types.String{Value: plan.GateId.Value}
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.Diagnostics.AddError(
			"Could not find Quality Gate Selection",
			fmt.Sprintf("The findSelection function was unable to find the project keys: %+v in the response: %+v", plan.ProjectKeys.Elems, res),
		)
		return
	}
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
	if result, ok := findSelection(res, state.ProjectKeys.Elems); ok {
		result.GateId = types.String{Value: state.GateId.Value}
		result.ID = types.String{Value: state.GateId.Value}
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
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
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sel, rem := diffSelection(state, plan)

	for _, s := range rem {
		deselectRequest := qualitygates.DeselectRequest{
			Organization: r.p.organization,
			ProjectKey:   s.(types.String).Value,
		}
		err := r.p.client.Qualitygates.Deselect(deselectRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not Deselect the Quality Gate selection",
				fmt.Sprintf("The Deselect request returned an error: %+v", err),
			)
			return
		}
	}
	for _, s := range sel {
		selectRequest := qualitygates.SelectRequest{
			GateId:       state.GateId.Value,
			Organization: r.p.organization,
			ProjectKey:   s.(types.String).Value,
		}
		err := r.p.client.Qualitygates.Select(selectRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not Select the Quality Gate selection",
				fmt.Sprintf("The Select request returned an error: %+v", err),
			)
			return
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
	if result, ok := findSelection(res, plan.ProjectKeys.Elems); ok {
		result.GateId = types.String{Value: state.GateId.Value}
		result.ID = types.String{Value: state.GateId.Value}
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.Diagnostics.AddError(
			"Could not find Quality Gate Selection",
			fmt.Sprintf("The findSelection function was unable to find the project keys: %+v in the response: %+v", plan.ProjectKeys.Elems, res),
		)
		return
	}
}

func (r resourceQualityGateSelection) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Selection
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, s := range state.ProjectKeys.Elems {
		request := qualitygates.DeselectRequest{
			Organization: r.p.organization,
			ProjectKey:   s.(types.String).Value,
		}
		err := r.p.client.Qualitygates.Deselect(request)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not Deselect the Quality Gate Selection",
				fmt.Sprintf("The Deselect request returned an error: %+v", err),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)
}

func diffSelection(state, plan Selection) (sel, rem []attr.Value) {
	for _, old := range state.ProjectKeys.Elems {
		// assume that old is a string
		if !containSelection(plan.ProjectKeys.Elems, old.(types.String).Value) {
			rem = append(rem, types.String{Value: old.(types.String).Value})
		}
	}
	for _, new := range plan.ProjectKeys.Elems {
		// assume that new is a string
		if !containSelection(state.ProjectKeys.Elems, new.(types.String).Value) {
			sel = append(sel, types.String{Value: new.(types.String).Value})
		}
	}

	return sel, rem
}

// Check if a condition is contained in a condition list
func containSelection(list []attr.Value, item string) bool {
	for _, c := range list {
		if c.Equal(types.String{Value: item}) {
			return true
		}
	}
	return false
}
