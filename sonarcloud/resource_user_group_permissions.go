package sonarcloud

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reinoudk/go-sonarcloud/sonarcloud"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/permissions"
	"sync"
	"time"
)

type resourceUserGroupPermissionsType struct{}

func (r resourceUserGroupPermissionsType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource manages the permissions of a user group for the whole organization or a specific project.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Description: "The implicit ID of the resource",
				Computed:    true,
			},
			"project_key": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The key of the project to restrict the permissions to.",
			},
			"group": {
				Type:        types.StringType,
				Required:    true,
				Description: "User group to set the permissions for.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"permissions": {
				Type:     types.SetType{ElemType: types.StringType},
				Required: true,
				Description: "List of permissions to grant." +
					" Available global permissions: [`admin`, `profileadmin`, `gateadmin`, `scan`, `provisioning`]." +
					" Available project permissions: ['admin`, `scan`, `codeviewer`, `issueadmin`, `securityhotspotadmin`, `user`].",
				Validators: []tfsdk.AttributeValidator{
					allowedSetOptions(
						// Global permissions
						"admin",
						"profileadmin",
						"gateadmin",
						"scan",
						"provisioning",
						// Project permissions
						// Note: admin and scan are project permissions as well
						"codeviewer",
						"issueadmin",
						"securityhotspotadmin",
						"user",
					),
				},
			},
		},
	}, nil
}

func (r resourceUserGroupPermissionsType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceUserGroupPermissions{
		p: *(p.(*provider)),
	}, nil
}

type resourceUserGroupPermissions struct {
	p provider
}

func (r resourceUserGroupPermissions) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unkown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan UserGroupPermissions
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Add permissions one by one
	wg := sync.WaitGroup{}
	for _, elem := range plan.Permissions.Elems {
		permission := elem.(types.String).Value

		go func() {
			wg.Add(1)
			defer wg.Done()

			request := permissions.AddGroupRequest{
				GroupName:    plan.Group.Value,
				Permission:   permission,
				ProjectKey:   plan.ProjectKey.Value,
				Organization: r.p.organization,
			}
			if err := r.p.client.Permissions.AddGroup(request); err != nil {
				resp.Diagnostics.AddError(
					"Could not add group permissions",
					fmt.Sprintf("The AddGroup request returned an error: %+v", err),
				)
				return
			}
		}()
	}

	// Async wait for all requests to be done
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	// Set ID on success and return error diag on timeout
	select {
	case <-c:
	case <-time.After(30 * time.Second):
		resp.Diagnostics.AddError("Could not set user group permissions",
			"The requests to set the permissions timed out.",
		)
	}

	plannedPermissions := make([]string, len(plan.Permissions.Elems))
	diags = plan.Permissions.ElementsAs(ctx, &plannedPermissions, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	backoffConfig := defaultBackoffConfig()

	group, err := backoff.RetryWithData(
		func() (*UserGroupPermissions, error) {
			group, err := findUserGroupWithPermissionsSet(r.p.client, plan.Group.Value, plan.ProjectKey.Value, plan.Permissions)
			return group, err
		}, backoffConfig)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not find the user group with the planned permissions",
			fmt.Sprintf("The findUserGroupWithPermissionsSet call returned an error: %+v ", err),
		)
	} else {
		diags = resp.State.Set(ctx, group)
		resp.Diagnostics.Append(diags...)
	}
}

func (r resourceUserGroupPermissions) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state UserGroupPermissions
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Query for permissions
	searchRequest := UserGroupPermissionsSearchRequest{ProjectKey: state.ProjectKey.Value}
	groups, err := sonarcloud.GetAll[UserGroupPermissionsSearchRequest, UserGroupPermissionsSearchResponseGroup](r.p.client, "/permissions/groups", searchRequest, "groups")
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get user group permissions",
			fmt.Sprintf("The request returned an error: %+v", err),
		)
		return
	}

	if group, ok := findUserGroup(groups, state.Group.Value); ok {
		permissionsElems := make([]attr.Value, len(group.Permissions))

		for i, permission := range group.Permissions {
			permissionsElems[i] = types.String{Value: permission}
		}

		result := UserGroupPermissions{
			ID:          state.ID,
			ProjectKey:  state.ProjectKey,
			Group:       state.Group,
			Permissions: types.Set{Elems: permissionsElems, ElemType: types.StringType},
		}
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r resourceUserGroupPermissions) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var state UserGroupPermissions
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan UserGroupPermissions
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	toAdd, toRemove := diffAttrSets(state.Permissions, plan.Permissions)

	for _, remove := range toRemove {
		removeRequest := permissions.RemoveGroupRequest{
			GroupName:    state.Group.Value,
			Permission:   remove.(types.String).Value,
			ProjectKey:   state.ProjectKey.Value,
			Organization: r.p.organization,
		}
		err := r.p.client.Permissions.RemoveGroup(removeRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not remove the user group permission",
				fmt.Sprintf("The RemoveGroup request returned an error: %+v", err),
			)
			return
		}
	}
	for _, add := range toAdd {
		addRequest := permissions.AddGroupRequest{
			GroupName:    plan.Group.Value,
			Permission:   add.(types.String).Value,
			ProjectKey:   plan.ProjectKey.Value,
			Organization: r.p.organization,
		}
		if err := r.p.client.Permissions.AddGroup(addRequest); err != nil {
			resp.Diagnostics.AddError(
				"Could not add the user group permission",
				fmt.Sprintf("The AddGroup request returned an error: %+v", err),
			)
			return
		}
	}

	plannedPermissions := make([]string, len(plan.Permissions.Elems))
	diags = plan.Permissions.ElementsAs(ctx, &plannedPermissions, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	backoffConfig := defaultBackoffConfig()

	group, err := backoff.RetryWithData(
		func() (*UserGroupPermissions, error) {
			group, err := findUserGroupWithPermissionsSet(r.p.client, plan.Group.Value, plan.ProjectKey.Value, plan.Permissions)
			return group, err
		}, backoffConfig)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not find the user group with the planned permissions",
			fmt.Sprintf("The findUserGroupWithPermissionsSet call returned an error: %+v ", err),
		)
	} else {
		diags = resp.State.Set(ctx, group)
		resp.Diagnostics.Append(diags...)
	}
}

func (r resourceUserGroupPermissions) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state UserGroupPermissions
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, remove := range state.Permissions.Elems {
		removeRequest := permissions.RemoveGroupRequest{
			GroupName:    state.Group.Value,
			Permission:   remove.(types.String).Value,
			ProjectKey:   state.ProjectKey.Value,
			Organization: r.p.organization,
		}
		err := r.p.client.Permissions.RemoveGroup(removeRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not remove the user group permission",
				fmt.Sprintf("The RemoveGroup request returned an error: %+v", err),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)
}

type UserGroupPermissionsSearchRequest struct {
	ProjectKey string
}

type UserGroupPermissionsSearchResponseGroup struct {
	Id          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

func findUserGroupWithPermissionsSet(client *sonarcloud.Client, groupName, projectKey string, expectedPermissions types.Set) (*UserGroupPermissions, error) {
	searchRequest := UserGroupPermissionsSearchRequest{ProjectKey: projectKey}
	groups, err := sonarcloud.GetAll[UserGroupPermissionsSearchRequest, UserGroupPermissionsSearchResponseGroup](client, "/permissions/groups", searchRequest, "groups")
	if err != nil {
		return nil, err
	}

	group, ok := findUserGroup(groups, groupName)
	if !ok {
		return nil, fmt.Errorf("group not found in response (groupName='%s',projectKey='%s')", groupName, projectKey)
	}

	permissionsElems := make([]attr.Value, len(group.Permissions))
	for i, permission := range group.Permissions {
		permissionsElems[i] = types.String{Value: permission}
	}

	foundPermissions := types.Set{Elems: permissionsElems, ElemType: types.StringType}

	if !foundPermissions.Equal(expectedPermissions) {
		return nil, fmt.Errorf("the returned permissions do not match the expected permissions (groupName='%s',projectKey='%s, expected='%v', got='%v')",
			groupName,
			projectKey,
			expectedPermissions,
			foundPermissions)
	}

	return &UserGroupPermissions{
		ID:          types.String{Value: projectKey + "-" + groupName},
		ProjectKey:  types.String{Value: projectKey},
		Group:       types.String{Value: groupName},
		Permissions: foundPermissions,
	}, nil
}
