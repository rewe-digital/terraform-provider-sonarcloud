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
				Type:        types.Float64Type,
				Description: "Id computed by SonarCloud servers",
				Computed:    true,
			},
			"name": {
				Type:        types.StringType,
				Description: "Name of the Quality Gate.",
				Required:    true,
			},
			"is_built_in": {
				Type:        types.BoolType,
				Description: "Defines whether the quality gate is built in. ",
				Computed:    true,
			},
			"is_default": {
				Type:        types.BoolType,
				Description: "Defines whether the quality gate is the defualt gate for an organization",
				Optional:    true,
			},
			"actions": {
				Description: "What actions can be performed on this Quality Gate.",
				Computed:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"rename": {
						Type:        types.BoolType,
						Description: "Whether this object can be renamed",
						Computed:    true,
					},
					"set_as_default": {
						Type:        types.BoolType,
						Description: "Whether this object can be set as Default",
						Computed:    true,
					},
					"copy": {
						Type:        types.BoolType,
						Description: "Whether this object can be copied",
						Computed:    true,
					},
					"associate_projects": {
						Type:        types.BoolType,
						Description: "Whether this object can be associated with Projects",
						Computed:    true,
					},
					"delete": {
						Type:        types.BoolType,
						Description: "Whether this object can be deleted",
						Computed:    true,
					},
					"manage_conditions": {
						Type:        types.BoolType,
						Description: "Whether this object can be managed",
						Computed:    true,
					},
				}),
			},
			"conditions": {
				Optional:    true,
				Description: "The conditions of this quality gate.",
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:        types.Float64Type,
						Description: "Index/ID of the Condition.",
						Computed:    true,
					},
					"metric": {
						Type:        types.StringType,
						Description: "The metric on which the condition is based. Must be one of: https://docs.sonarqube.org/latest/user-guide/metric-definitions/",
						Required:    true,
						Validators: []tfsdk.AttributeValidator{
							allowedOptions("security_rating", "ncloc_language_distribution", "test_execution_time", "statements", "lines_to_cover", "quality_gate_details", "new_reliabillity_remediation_effort", "tests", "security_review_rating", "new_xxx_violations", "conditions_by_line", "new_violations", "ncloc", "duplicated_lines", "test_failures", "test_errors", "reopened_issues", "new_vulnerabilities", "duplicated_lines_density", "test_success_density", "sqale_debt_ratio", "security_hotspots_reviewed", "security_remediation_effort", "covered_conditions_by_line", "classes", "sqale_rating", "xxx_violations", "true_positive_issues", "violations", "new_security_review_rating", "new_security_remediation_effort", "vulnerabillities", "new_uncovered_conditions", "files", "branch_coverage_hits_data", "uncovered_lines", "comment_lines_density", "new_uncovered_lines", "complexty", "cognitive_complexity", "uncovered_conditions", "functions", "new_technical_debt", "new_coverage", "coverage", "new_branch_coverage", "confirmed_issues", "reliabillity_remediation_effort", "projects", "coverage_line_hits_data", "code_smells", "directories", "lines", "bugs", "line_coverage", "new_line_coverage", "reliability_rating", "duplicated_blocks", "branch_coverage", "new_code_smells", "new_sqale_debt_ratio", "open_issues", "sqale_index", "new_lines_to_cover", "comment_lines", "skipped_tests"),
						},
					},
					"op": {
						Type:        types.StringType,
						Description: "Operation on which the metric is evaluated must be either: LT, GT",
						Optional:    true,
						Validators: []tfsdk.AttributeValidator{
							allowedOptions("LT", "GT"),
						},
					},
					"error": {
						Type:        types.Float64Type, // TODO: Change to StringType once go-sonarcloud error has been fixed
						Description: "The value on which the condition errors.",
						Required:    true,
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
	//retrieve values from state
	var state QualityGate
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan QualityGate
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	changed := changedAttrs(req, diags)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, ok := changed["name"]; ok {
		request := qualitygates.RenameRequest{
			Id:           state.ID.String(),
			Name:         plan.Name.Value,
			Organization: r.p.organization,
		}

		err := r.p.client.Qualitygates.Rename(request)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not update Quality Gate Name.",
				fmt.Sprintf("The Rename request returned an error: %+v", err),
			)
			return
		}
	}

	if _, ok := changed["isDefault"]; ok {
		if plan.IsDefault.Value {
			request := qualitygates.SetAsDefaultRequest{
				Id:           state.ID.String(),
				Organization: r.p.organization,
			}
			err := r.p.client.Qualitygates.SetAsDefault(request)
			if err != nil {
				resp.Diagnostics.AddError(
					"Could not set Quality Gate as Default.",
					fmt.Sprintf("The SetAsDefault request returned an error %+v", err),
				)
				return
			}
		}
	}

	if _, ok := changed["conditions"]; ok {
		// delete all state conditions if there are none in the plan.
		if len(plan.Conditions) < 0 {
			for _, condition := range state.Conditions {
				request := qualitygates.DeleteConditionRequest{
					Id:           condition.ID.String(),
					Organization: r.p.organization,
				}
				err := r.p.client.Qualitygates.DeleteCondition(request)
				if err != nil {
					resp.Diagnostics.AddError(
						"Could not delete Quality Gate Condition.",
						fmt.Sprintf("The DeleteCondition request returned an error %+v", err),
					)
					return
				}
			}
		} else {
			// This is long, maybe unnecessarily complex, and its effiency is questionable.
			// I'm not to certain how to evaluate state vs. plan changes for nested attribute lists.
			stateUpdate := Conditions{}
			stateCreate := Conditions{}
			stateDelete := state.Conditions
			for _, planCond := range plan.Conditions {
				notFound := true
				for i := len(state.Conditions); i >= 0; i-- {
					if planCond.Metric == state.Conditions[i].Metric {
						notFound = false
						stateDelete = append(stateDelete[:i], stateDelete[i+1:]...)
						updatedCondition := Condition{
							ID:     state.Conditions[i].ID,
							Error:  planCond.Error,
							Metric: planCond.Metric,
							Op:     planCond.Op,
						}
						stateUpdate.Conditions = append(stateUpdate.Conditions, updatedCondition)
					}
				}
				if notFound {
					stateCreate.Conditions = append(stateCreate.Conditions, planCond)
				}
			}

			if len(stateUpdate.Conditions) > 0 {
				for _, updateCondition := range stateUpdate.Conditions {
					request := qualitygates.UpdateConditionRequest{
						Error:        updateCondition.Error.String(), // TODO: Change when go-sonarcloud error has been fixed
						Id:           updateCondition.ID.String(),
						Metric:       updateCondition.Metric.Value,
						Op:           updateCondition.Op.Value,
						Organization: r.p.organization,
					}

					err := r.p.client.Qualitygates.UpdateCondition(request)
					if err != nil {
						resp.Diagnostics.AddError(
							"Could not update QualityGate condition",
							fmt.Sprintf("The UpdateCondition request returned an error %+v", err),
						)
						return
					}
				}
			}
			if len(stateCreate.Conditions) > 0 {
				for _, createCondition := range stateCreate.Conditions {
					request := qualitygates.CreateConditionRequest{
						Error:        createCondition.Error.String(), // TODO: Change when go-sonarcloud error has been fixed
						Metric:       createCondition.Metric.Value,
						Op:           createCondition.Op.Value,
						Organization: r.p.organization,
					}
					_, err := r.p.client.Qualitygates.CreateCondition(request)
					if err != nil {
						resp.Diagnostics.AddError(
							"Could not create QualityGate condition",
							fmt.Sprintf("The CreateCondition request returned an error %+v", err),
						)
						return
					}
				}
			}
			if len(stateDelete) > 0 {
				for _, deleteCondition := range stateDelete {
					request := qualitygates.DeleteConditionRequest{
						Id:           deleteCondition.ID.String(),
						Organization: r.p.organization,
					}
					err := r.p.client.Qualitygates.DeleteCondition(request)
					if err != nil {
						resp.Diagnostics.AddError(
							"Could not delete QualityGate condition",
							fmt.Sprintf("The DeleteCondition request returned an error %+v", err),
						)
						return
					}
				}
			}
		}
	}
	// There aren't any return values for non-create operations.
	listRequest := qualitygates.ListRequest{}

	response, err := r.p.client.Qualitygates.List(listRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the Quality Gate",
			fmt.Sprintf("The List request returned an error: %+v", err),
		)
		return
	}

	if result, ok := findQualityGate(response, plan.Name.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	}
}

func (r resourceQualityGate) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// Retrieve values from state
	var state QualityGate
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := qualitygates.DestroyRequest{
		Id:           state.ID.String(),
		Organization: r.p.organization,
	}

	err := r.p.client.Qualitygates.Destroy(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not destroy the quality gate",
			fmt.Sprintf("The Destroy request returned an error: %+v", err),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r resourceQualityGate) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateRequest) {
	//TODO: Implement Import
}
