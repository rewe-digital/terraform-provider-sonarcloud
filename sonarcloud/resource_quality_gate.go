package sonarcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/qualitygates"
)

type resourceQualityGateType struct{}

func (r resourceQualityGateType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource manages a Quality Gate",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:          types.Float64Type,
				Attributes:    nil,
				Description:   "Id computed by SonarCloud servers",
				Required:      false,
				Optional:      false,
				Computed:      true,
				Sensitive:     false,
				Validators:    []tfsdk.AttributeValidator{},
				PlanModifiers: []tfsdk.AttributePlanModifier{},
			},
			"name": {
				Type:          types.StringType,
				Attributes:    nil,
				Description:   "Name of the Quality Gate.",
				Required:      true,
				Optional:      false,
				Computed:      false,
				Sensitive:     false,
				Validators:    []tfsdk.AttributeValidator{},
				PlanModifiers: []tfsdk.AttributePlanModifier{},
			},
			"isBuiltIn": {
				Type:          types.BoolType,
				Attributes:    nil,
				Description:   "Defines whether the quality gate is built in. ",
				Required:      false,
				Optional:      true,
				Computed:      true,
				Sensitive:     false,
				Validators:    []tfsdk.AttributeValidator{},
				PlanModifiers: []tfsdk.AttributePlanModifier{},
			},
			"isDefault": {
				Type:          types.BoolType,
				Attributes:    nil,
				Description:   "Defines whether the quality gate is the defualt gate for an organization",
				Required:      false,
				Optional:      true,
				Computed:      false,
				Sensitive:     false,
				Validators:    []tfsdk.AttributeValidator{},
				PlanModifiers: []tfsdk.AttributePlanModifier{},
			},
			"actions": {
				Description:   "What actions can be performed on this Quality Gate.",
				Required:      false,
				Optional:      true,
				Computed:      true,
				Sensitive:     false,
				Validators:    []tfsdk.AttributeValidator{},
				PlanModifiers: []tfsdk.AttributePlanModifier{},
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"rename": {
						Type:          types.BoolType,
						Description:   "Whether this object can be renamed",
						Required:      false,
						Optional:      false,
						Computed:      true,
						Sensitive:     false,
						Validators:    []tfsdk.AttributeValidator{},
						PlanModifiers: []tfsdk.AttributePlanModifier{},
					},
					"setAsDefault": {
						Type:          types.BoolType,
						Description:   "Whether this object can be set as Default",
						Required:      false,
						Optional:      false,
						Computed:      true,
						Sensitive:     false,
						Validators:    []tfsdk.AttributeValidator{},
						PlanModifiers: []tfsdk.AttributePlanModifier{},
					},
					"copy": {
						Type:          types.BoolType,
						Description:   "Whether this object can be copied",
						Required:      false,
						Optional:      false,
						Computed:      true,
						Sensitive:     false,
						Validators:    []tfsdk.AttributeValidator{},
						PlanModifiers: []tfsdk.AttributePlanModifier{},
					},
					"associateProjects": {
						Type:          types.BoolType,
						Description:   "Whether this object can be associated with Projects",
						Required:      false,
						Optional:      false,
						Computed:      true,
						Sensitive:     false,
						Validators:    []tfsdk.AttributeValidator{},
						PlanModifiers: []tfsdk.AttributePlanModifier{},
					},
					"delete": {
						Type:          types.BoolType,
						Description:   "Whether this object can be deleted",
						Required:      false,
						Optional:      false,
						Computed:      true,
						Sensitive:     false,
						Validators:    []tfsdk.AttributeValidator{},
						PlanModifiers: []tfsdk.AttributePlanModifier{},
					},
					"manageConditions": {
						Type:          types.BoolType,
						Description:   "Whether this object can be managed",
						Required:      false,
						Optional:      false,
						Computed:      true,
						Sensitive:     false,
						Validators:    []tfsdk.AttributeValidator{},
						PlanModifiers: []tfsdk.AttributePlanModifier{},
					},
				}),
			},
			"conditions": {
				Computed:    true,
				Optional:    true,
				Description: "The conditions of this quality gate.",
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:          types.Float64Type,
						Description:   "Index/ID of the Condition.",
						Required:      false,
						Optional:      false,
						Computed:      true,
						Sensitive:     false,
						Validators:    []tfsdk.AttributeValidator{},
						PlanModifiers: []tfsdk.AttributePlanModifier{},
					},
					"metric": {
						Type:        types.StringType,
						Attributes:  nil,
						Description: "The metric on which the condition is based. Must be one of: https://docs.sonarqube.org/latest/user-guide/metric-definitions/",
						Required:    true,
						Optional:    false,
						Computed:    false,
						Sensitive:   false,
						Validators: []tfsdk.AttributeValidator{
							allowedOptions("security_rating", "ncloc_language_distribution", "test_execution_time", "statements", "lines_to_cover", "quality_gate_details", "new_reliabillity_remediation_effort", "tests", "security_review_rating", "new_xxx_violations", "conditions_by_line", "new_violations", "ncloc", "duplicated_lines", "test_failures", "test_errors", "reopened_issues", "new_vulnerabilities", "duplicated_lines_density", "test_success_density", "sqale_debt_ratio", "security_hotspots_reviewed", "security_remediation_effort", "covered_conditions_by_line", "classes", "sqale_rating", "xxx_violations", "true_positive_issues", "violations", "new_security_review_rating", "new_security_remediation_effort", "vulnerabillities", "new_uncovered_conditions", "files", "branch_coverage_hits_data", "uncovered_lines", "comment_lines_density", "new_uncovered_lines", "complexty", "cognitive_complexity", "uncovered_conditions", "functions", "new_technical_debt", "new_coverage", "coverage", "new_branch_coverage", "confirmed_issues", "reliabillity_remediation_effort", "projects", "coverage_line_hits_data", "code_smells", "directories", "lines", "bugs", "line_coverage", "new_line_coverage", "reliability_rating", "duplicated_blocks", "branch_coverage", "new_code_smells", "new_sqale_debt_ratio", "open_issues", "sqale_index", "new_lines_to_cover", "comment_lines", "skipped_tests"),
						},
						PlanModifiers: []tfsdk.AttributePlanModifier{},
					},
					"op": {
						Type:        types.StringType,
						Attributes:  nil,
						Description: "Operation on which the metric is evaluated must be either: LT, GT",
						Required:    false,
						Optional:    true,
						Computed:    false,
						Sensitive:   false,
						Validators: []tfsdk.AttributeValidator{
							allowedOptions("LT", "GT"),
						},
						PlanModifiers: []tfsdk.AttributePlanModifier{},
					},
					"error": {
						Type:          types.Float64Type,
						Attributes:    nil,
						Description:   "The value on which the condition errors.",
						Required:      true,
						Optional:      false,
						Computed:      false,
						Sensitive:     false,
						Validators:    []tfsdk.AttributeValidator{},
						PlanModifiers: []tfsdk.AttributePlanModifier{},
					},
				}),
			},
		},
	}, nil
}

func (r resourceQualityGateType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceQualityGate{
		p: *(p.(*provider)),
	}, nil
}

type resourceQualityGate struct {
	p provider
}

func (r resourceQualityGate) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unkown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan QualityGate
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct for Quality Gates
	request := qualitygates.CreateRequest{
		Name:         plan.Name.Value,
		Organization: r.p.organization,
	}

	res, err := r.p.client.Qualitygates.Create(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create the Quality Gate",
			fmt.Sprintf("The Quality Gate create request returned an error: %+v", err),
		)
		return
	}

	var result = QualityGate{
		ID:   types.Float64{Value: res.Id},
		Name: types.String{Value: res.Name},
	}
	tempQualityGateId := res.Id

	conditionRequests := qualitygates.CreateConditionRequest{}
	for i, conditionPlan := range plan.Conditions {
		conditionRequests = qualitygates.CreateConditionRequest{
			Error:        fmt.Sprintf("%f", conditionPlan.Error.Value),
			GateId:       fmt.Sprintf("%f", tempQualityGateId),
			Metric:       conditionPlan.Metric.Value,
			Op:           conditionPlan.Op.Value,
			Organization: r.p.organization,
		}
		res, err := r.p.client.Qualitygates.CreateCondition(conditionRequests)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not create a Condition",
				fmt.Sprintf("The Condition Create Request returned an error: %+v", err),
			)
			return
		}
		// didn't implement warning
		result.Conditions[i] = Condition{
			Error:  types.Float64{Value: res.Error},
			ID:     types.Float64{Value: res.Id},
			Metric: types.String{Value: res.Metric},
			Op:     types.String{Value: res.Op},
		}
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceQualityGate) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	//Retrieve values from state
	var state QualityGate
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := qualitygates.ListRequest{
		Organization: r.p.organization,
	}

	response, err := r.p.client.Qualitygates.List(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the Quality Gate(s)",
			fmt.Sprintf("The List request returned an error: %+v", err),
		)
		return
	}

	// Check if the resource exists in the list of retrieved resources
	if result, ok := findQualityGate(response, state.Name.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r resourceQualityGate) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {

}

func (r resourceQualityGate) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	//TODO: Implement Delete
}

func (r resourceQualityGate) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateRequest) {
	//TODO: Implement Import
}
