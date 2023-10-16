package sonarcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/webhooks"
)

type resourceWebhookType struct{}

func (r resourceWebhookType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource represents a project or organization webhook.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Computed:    true,
				Description: "ID of the webhook, this is equal to its key.",
			},
			"key": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Key of the webhook.",
			},
			"project": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The key of the project to add the webhook to. If empty, the webhook will be added to the organization.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name of the webhook.",
			},
			"secret": {
				Type:        types.StringType,
				Optional:    true,
				Description: "If provided, secret will be used as the key to generate the HMAC hex (lowercase) digest value in the 'X-Sonar-Webhook-HMAC-SHA256' header.",
				Sensitive:   true,
			},
			"url": {
				Type:        types.StringType,
				Required:    true,
				Description: "The url of the webhook.",
			},
		},
	}, nil
}

func (r resourceWebhookType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceWebhook{
		p: *(p.(*provider)),
	}, nil
}

type resourceWebhook struct {
	p provider
}

func (r resourceWebhook) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Webhook
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := webhooks.CreateRequest{
		Name:         plan.Name.Value,
		Organization: r.p.organization,
		Project:      plan.Project.Value,
		Secret:       plan.Secret.Value,
		Url:          plan.Url.Value,
	}

	res, err := r.p.client.Webhooks.Create(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create the webhook",
			fmt.Sprintf("The Create request returned an error: %+v", err),
		)
		return
	}

	webhook := res.Webhook
	var result = Webhook{
		ID:      types.String{Value: webhook.Key},
		Key:     types.String{Value: webhook.Key},
		Project: plan.Project,
		Name:    types.String{Value: webhook.Name},
		// Just use the secret from the plan, as it's not returned by the API
		Secret: plan.Secret,
		Url:    types.String{Value: webhook.Url},
	}
	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)
}

func (r resourceWebhook) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Retrieve values from state
	var state Webhook
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := webhooks.ListRequest{
		Organization: r.p.organization,
		Project:      state.Project.Value,
	}

	response, err := r.p.client.Webhooks.List(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the webhooks",
			fmt.Sprintf("The List request returned an error: %+v", err),
		)
		return
	}

	// Check if the resource exists the list of retrieved resources
	if result, ok := findWebhook(response, state.ID.Value, state.Project.Value, state.Secret.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r resourceWebhook) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Retrieve values from state
	var state Webhook
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan Webhook
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill in api action struct
	request := webhooks.UpdateRequest{
		Name:   plan.Name.Value,
		Secret: plan.Secret.Value,
		Url:    plan.Url.Value,
		// Note: this is an inconsistency in the API naming...
		Webhook: state.Key.Value,
	}

	err := r.p.client.Webhooks.Update(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not update the webhook",
			fmt.Sprintf("The Update request returned an error: %+v", err),
		)
		return
	}

	// We don't have a return value, so we have to query it again
	// Fill in api action struct
	listRequest := webhooks.ListRequest{
		Organization: r.p.organization,
		Project:      state.Project.Value,
	}

	response, err := r.p.client.Webhooks.List(listRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the webhooks",
			fmt.Sprintf("The List request returned an error: %+v", err),
		)
		return
	}

	// Check if the resource exists the list of retrieved resources
	if result, ok := findWebhook(response, state.Key.Value, state.Project.Value, plan.Secret.Value); ok {
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	}
}

func (r resourceWebhook) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Webhook
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := webhooks.DeleteRequest{
		// Note: this is an inconsistency in the API naming...
		Webhook: state.Key.Value,
	}
	err := r.p.client.Webhooks.Delete(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete the webhook",
			fmt.Sprintf("The Delete request returned an error: %+v", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceWebhook) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) < 1 || len(idParts) > 2 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id OR id,project_key. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
	if len(idParts) == 2 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project"), idParts[1])...)
	}
}

// findWebhook returns the link with the given id, if it exists in the response
func findWebhook(response *webhooks.ListResponse, key, project_key, secret string) (Webhook, bool) {
	var result Webhook
	ok := false

	// If project_key is an empty string, we need to explicitly set 'Null' to 'true' in the types.String struct.
	// Otherwise, it would be in an invalid state, which leads to potentially indeterminate behaviour.
	// This is "fixed" in https://github.com/hashicorp/terraform-plugin-framework/pull/523 with explicit constructor
	// functions that ensure a valid state.
	// TODO: upgrade terraform provider framework dependency so we can use an explicit constructor
	var projectKeyIsNull bool
	if project_key == "" {
		projectKeyIsNull = true
	} else {
		projectKeyIsNull = false
	}

	for _, webhook := range response.Webhooks {
		if webhook.Key == key {
			result = Webhook{
				ID:      types.String{Value: webhook.Key},
				Key:     types.String{Value: webhook.Key},
				Project: types.String{Value: project_key, Null: projectKeyIsNull},
				Name:    types.String{Value: webhook.Name},
				// We have to use the secret from the plan, as it's not returned by the API
				Secret: types.String{Value: secret},
				Url:    types.String{Value: webhook.Url},
			}
			ok = true
			break
		}
	}
	return result, ok
}
