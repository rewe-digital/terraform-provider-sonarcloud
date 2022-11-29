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

type resourceUserPermissionsType struct{}

func (r resourceUserPermissionsType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource manages the permissions of a user for the whole organization or a specific project.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Description: "The implicit ID of the resource.",
				Computed:    true,
			},
			"project_key": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The key of the project to restrict the permissions to.",
			},
			"login": {
				Type:        types.StringType,
				Required:    true,
				Description: "The login of the user to set the permissions for.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"name": {
				Type:        types.StringType,
				Computed:    true,
				Description: "The name of the user.",
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
			"avatar": {
				Type:        types.StringType,
				Computed:    true,
				Description: "The avatar ID of the user.",
			},
		},
	}, nil
}

func (r resourceUserPermissionsType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceUserPermissions{
		p: *(p.(*provider)),
	}, nil
}

type resourceUserPermissions struct {
	p provider
}

func (r resourceUserPermissions) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unkown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan UserPermissions
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

			request := permissions.AddUserRequest{
				Login:        plan.Login.Value,
				Permission:   permission,
				ProjectKey:   plan.ProjectKey.Value,
				Organization: r.p.organization,
			}
			if err := r.p.client.Permissions.AddUser(request); err != nil {
				resp.Diagnostics.AddError(
					"Could not add user permissions",
					fmt.Sprintf("The AddUser request returned an error: %+v", err),
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
		resp.Diagnostics.AddError("Could not set user user permissions",
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

	user, err := backoff.RetryWithData(
		func() (*UserPermissions, error) {
			user, err := findUserWithPermissionsSet(r.p.client, plan.Login.Value, plan.ProjectKey.Value, plan.Permissions)
			return user, err
		}, backoffConfig)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not find the user with the planned permissions",
			fmt.Sprintf("The findUserWithPermissionsSet call returned an error: %+v ", err),
		)
	} else {
		diags = resp.State.Set(ctx, user)
		resp.Diagnostics.Append(diags...)
	}
}

func (r resourceUserPermissions) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state UserPermissions
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Query for permissions
	searchRequest := UserPermissionsSearchRequest{ProjectKey: state.ProjectKey.Value}
	users, err := sonarcloud.GetAll[UserPermissionsSearchRequest, UserPermissionsSearchResponseUser](r.p.client, "/permissions/users", searchRequest, "users")
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get user permissions",
			fmt.Sprintf("The request returned an error: %+v", err),
		)
		return
	}

	if user, ok := findUser(users, state.Login.Value); ok {
		permissionsElems := make([]attr.Value, len(user.Permissions))

		for i, permission := range user.Permissions {
			permissionsElems[i] = types.String{Value: permission}
		}

		result := UserPermissions{
			ID:          types.String{Value: state.ProjectKey.Value + "-" + state.Login.Value},
			ProjectKey:  state.ProjectKey,
			Login:       types.String{Value: user.Login},
			Name:        types.String{Value: user.Name},
			Permissions: types.Set{Elems: permissionsElems, ElemType: types.StringType},
			Avatar:      types.String{Value: user.Avatar},
		}
		diags = resp.State.Set(ctx, result)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r resourceUserPermissions) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var state UserPermissions
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan UserPermissions
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	toAdd, toRemove := diffAttrSets(state.Permissions, plan.Permissions)

	for _, remove := range toRemove {
		removeRequest := permissions.RemoveUserRequest{
			Login:        state.Login.Value,
			Organization: r.p.organization,
			Permission:   remove.(types.String).Value,
			ProjectKey:   state.ProjectKey.Value,
		}
		err := r.p.client.Permissions.RemoveUser(removeRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not remove the permission",
				fmt.Sprintf("The RemoveUser request returned an error: %+v", err),
			)
			return
		}
	}
	for _, add := range toAdd {
		addRequest := permissions.AddUserRequest{
			Login:        plan.Login.Value,
			Permission:   add.(types.String).Value,
			ProjectKey:   plan.ProjectKey.Value,
			Organization: r.p.organization,
		}
		if err := r.p.client.Permissions.AddUser(addRequest); err != nil {
			resp.Diagnostics.AddError(
				"Could not add the user permission",
				fmt.Sprintf("The AddUser request returned an error: %+v", err),
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

	user, err := backoff.RetryWithData(
		func() (*UserPermissions, error) {
			return findUserWithPermissionsSet(r.p.client, plan.Login.Value, plan.ProjectKey.Value, plan.Permissions)
		}, backoffConfig)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not find the user with the planned permissions",
			fmt.Sprintf("The findUserWithPermissionsSet call returned an error: %+v ", err),
		)
	} else {
		diags = resp.State.Set(ctx, user)
		resp.Diagnostics.Append(diags...)
	}
}

func (r resourceUserPermissions) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state UserPermissions
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, remove := range state.Permissions.Elems {
		removeRequest := permissions.RemoveUserRequest{
			Login:        state.Login.Value,
			Organization: r.p.organization,
			Permission:   remove.(types.String).Value,
			ProjectKey:   state.ProjectKey.Value,
		}
		err := r.p.client.Permissions.RemoveUser(removeRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not remove the user permission",
				fmt.Sprintf("The RemoveUser request returned an error: %+v", err),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)
}

type UserPermissionsSearchRequest struct {
	ProjectKey string
}

type UserPermissionsSearchResponseUser struct {
	Id          string   `json:"id,omitempty"`
	Login       string   `json:"login,omitempty"`
	Name        string   `json:"name,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Avatar      string   `json:"avatar,omitempty"`
}

func findUserWithPermissionsSet(client *sonarcloud.Client, login, projectKey string, expectedPermissions types.Set) (*UserPermissions, error) {
	searchRequest := UserGroupPermissionsSearchRequest{ProjectKey: projectKey}
	users, err := sonarcloud.GetAll[UserGroupPermissionsSearchRequest, UserPermissionsSearchResponseUser](client, "/permissions/users", searchRequest, "users")
	if err != nil {
		return nil, err
	}

	user, ok := findUser(users, login)
	if !ok {
		return nil, fmt.Errorf("user not found in response (login='%s',projectKey='%s')", login, projectKey)
	}

	permissionsElems := make([]attr.Value, len(user.Permissions))
	for i, permission := range user.Permissions {
		permissionsElems[i] = types.String{Value: permission}
	}

	foundPermissions := types.Set{Elems: permissionsElems, ElemType: types.StringType}

	if !foundPermissions.Equal(expectedPermissions) {
		return nil, fmt.Errorf("the returned permissions do not match the expected permissions (login='%s',projectKey='%s, expected='%v', got='%v')",
			login,
			projectKey,
			expectedPermissions,
			foundPermissions)
	}

	return &UserPermissions{
		ID:          types.String{Value: projectKey + "-" + login},
		ProjectKey:  types.String{Value: projectKey},
		Login:       types.String{Value: user.Login},
		Name:        types.String{Value: user.Name},
		Permissions: foundPermissions,
		Avatar:      types.String{Value: user.Avatar},
	}, nil

}
