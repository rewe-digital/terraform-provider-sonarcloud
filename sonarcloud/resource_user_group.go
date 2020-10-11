package sonarcloud

import (
	"context"
	json "encoding/json"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
	"terraform-provider-sonarcloud/pkg/api"
)

func resourceSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the user group",
				ForceNew:    false,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description for the user group",
				ForceNew:    false,
			},
		},
	}
}

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get all resource values
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	// Fill in api action struct
	create := api.UserGroupsCreate{
		Description: description,
		Name:        name,
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&create)
	if err != nil {
		return diag.FromErr(err)
	}

	// Cast m to SonarClient and create POST request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequest("POST", fmt.Sprintf("%s/user_groups/create", API), strings.NewReader(values.Encode()))
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

	// Decode JSON response to response struct
	createResponse := &UserGroupCreateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&createResponse)
	if err != nil {
		return diag.Errorf("Decode error: %+v", err)
	}

	// Set resource id
	d.SetId(strconv.Itoa(createResponse.Id))

	return diags
}

func resourceUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get resource value that's needed to read the remote resource
	name := d.Get("name").(string)

	// Fill in api action struct
	// Note: there is no direct endpoint to get the team, so we have to do a search and filter by name.
	search := api.UserGroupsSearch{
		Q: name,
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&search)
	if err != nil {
		return diag.FromErr(err)
	}

	// Cast m to SonarClient and create GET request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequest("GET", fmt.Sprintf("%s/user_groups/search", API), strings.NewReader(values.Encode()))
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
	searchResponse := &UserGroupsSearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		return diag.Errorf("Decode error: %+v", err)
	}

	// Check if the resource exists the list of retrieved resources
	// TODO: anti-corruption layer that hides this implementation detail
	found := false
	for _, g := range searchResponse.Groups {
		if g.Name == name {
			found = true
			if err := d.Set("description", g.Description); err != nil {
				return diag.FromErr(err)
			}
			break
		}
	}

	// Unset the id if the resource has not been found
	if !found {
		d.SetId("")
	}

	return diags
}

func resourceUserGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Check if any of the resource's values have been changed
	if !d.HasChange("description") && !d.HasChange("name") {
		return resourceUserGroupRead(ctx, d, m)
	}

	// Get all resource values, including id
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	// Fill in api action struct
	// Note: we skip values that have not been changed
	// TODO: create patch function to only update field that have been changed (otherwise the call might fail..)
	update := api.UserGroupsUpdate{
		Id: id,
	}
	if d.HasChange("name") {
		update.Name = name
	}
	if d.HasChange("description") {
		update.Description = description
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&update)
	if err != nil {
		return diag.FromErr(err)
	}

	// Cast m to SonarClient and create POST request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequest("POST", fmt.Sprintf("%s/user_groups/update", API), strings.NewReader(values.Encode()))
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
		diags := diag.Diagnostics{}
		return diagErrorResponse(resp, diags)
	}

	// Update the resource's state by calling the Read function
	return resourceUserGroupRead(ctx, d, m)
}

func resourceUserGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get resource value that's needed to read the remote resource
	name := d.Get("name").(string)

	// Use name because the organization is always set and using an id will then throw an error...
	del := api.UserGroupsDelete{
		Name: name,
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&del)
	if err != nil {
		return diag.FromErr(err)
	}

	// Cast m to SonarClient and create GET request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequest("POST", fmt.Sprintf("%s/user_groups/delete", API), strings.NewReader(values.Encode()))
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
