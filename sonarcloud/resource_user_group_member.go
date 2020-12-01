package sonarcloud

import (
	"context"
	json "encoding/json"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"terraform-provider-sonarcloud/pkg/api"
)

type UserGroupUsersResponse struct {
	Users []UserGroupUser `json:"users"`
}

type UserGroupUser struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

func resourceSourceUserGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupMemberCreate,
		ReadContext:   resourceUserGroupMemberRead,
		DeleteContext: resourceUserGroupMemberDelete,
		Schema: map[string]*schema.Schema{
			"group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group name",
				ForceNew:    true,
			},
			"login": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User login",
				ForceNew:    true,
			},
		},
	}
}

func resourceUserGroupMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get all resource values
	group := d.Get("group").(string)
	login := d.Get("login").(string)

	// Fill in api action struct
	create := api.UserGroupsAddUser{
		Login: login,
		Name:  group,
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&create)
	if err != nil {
		return diag.FromErr(err)
	}

	// Cast m to SonarClient and create POST request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequest("POST", fmt.Sprintf("%s/user_groups/add_user", API), strings.NewReader(values.Encode()))
	if err != nil {
		return diag.FromErr(err)
	}

	// Perform the request
	resp, err := sc.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Check status code and return diagnostics from ErrorResponse if needed
	if resp.StatusCode >= 300 {
		return diagErrorResponse(resp, diags)
	}

	// Set resource id
	d.SetId(group + login)

	return diags
}

func resourceUserGroupMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get resource value that's needed to read the remote resource
	group := d.Get("group").(string)
	login := d.Get("login").(string)

	// Cast m to SonarClient and create GET request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequestWithParameters("GET", fmt.Sprintf("%s/user_groups/users", API), "name", group)
	if err != nil {
		return diag.FromErr(err)
	}

	// Perform the request
	resp, err := sc.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Check status code and return diagnostics from ErrorResponse if needed
	if resp.StatusCode != 200 {
		return diagErrorResponse(resp, diags)
	}

	// Decode JSON response to response struct
	usersResponse := &UserGroupUsersResponse{}
	err = json.NewDecoder(resp.Body).Decode(&usersResponse)
	if err != nil {
		return diag.Errorf("Decode error during Read: %+v", err)
	}

	// Check if the resource exists in the list of retrieved resources
	// TODO: anti-corruption layer that hides this implementation detail
	found := false
	for _, u := range usersResponse.Users {
		if u.Login == login {
			found = true
			break
		}
	}

	// Unset the id if the resource has not been found
	if !found {
		d.SetId("")
	}

	return diags
}

func resourceUserGroupMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get all resource values
	group := d.Get("group").(string)
	login := d.Get("login").(string)

	// Use name because the organization is always set and using an id will then throw an error...
	del := api.UserGroupsRemoveUser{
		Login: login,
		Name:  group,
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&del)
	if err != nil {
		return diag.FromErr(err)
	}

	// Cast m to SonarClient and create GET request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequest("POST", fmt.Sprintf("%s/user_groups/remove_user", API), strings.NewReader(values.Encode()))
	if err != nil {
		return diag.FromErr(err)
	}

	// Perform the request
	resp, err := sc.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Check status code and return diagnostics from ErrorResponse if needed
	if resp.StatusCode >= 300 {
		return diagErrorResponse(resp, diags)
	}

	// Unset the id
	d.SetId("")

	return diags
}
