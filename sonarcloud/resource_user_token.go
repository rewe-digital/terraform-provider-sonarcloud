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

func resourceUserToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserTokenCreate,
		ReadContext:   resourceUserTokenRead,
		DeleteContext: resourceUserTokenDelete,
		Schema: map[string]*schema.Schema{
			"login": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User login",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the token",
				ForceNew:    true,
			},
			"token": {
				Type:        schema.TypeString,
				Description: "Value of the token",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceUserTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get all resource values
	login := d.Get("login").(string)
	name := d.Get("name").(string)

	// Fill in api action struct
	create := api.UserTokensGenerate{
		Login: login,
		Name:  name,
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&create)
	if err != nil {
		return diag.FromErr(err)
	}

	// Cast m to SonarClient and create POST request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequest("POST", fmt.Sprintf("%s/user_tokens/generate", API), strings.NewReader(values.Encode()))
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
	tokenResponse := &api.UserTokensGenerateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return diag.Errorf("Decode error during Read: %+v", err)
	}

	// Set resource id and token value
	if err = d.Set("token", tokenResponse.Token); err != nil {
		return diag.Errorf("Could not set token value: %+v", err)
	}
	d.SetId(login + name)

	return diags
}

func resourceUserTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get resource value that's needed to read the remote resource
	login := d.Get("login").(string)
	name := d.Get("name").(string)

	// Cast m to SonarClient and create GET request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequestWithParameters("GET", fmt.Sprintf("%s/user_tokens/search", API), "login", login)
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
	tokenResponse := &api.UserTokensSearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return diag.Errorf("Decode error during Read: %+v", err)
	}

	// Check if the resource exists in the list of retrieved resources
	// TODO: anti-corruption layer that hides this implementation detail
	found := false
	for _, t := range tokenResponse.Tokens {
		if t.Name == name {
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

func resourceUserTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get all resource values
	login := d.Get("login").(string)
	name := d.Get("name").(string)

	// Use name because the organization is always set and using an id will then throw an error...
	del := api.UserTokensRevoke{
		Login: login,
		Name:  name,
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&del)
	if err != nil {
		return diag.FromErr(err)
	}

	// Cast m to SonarClient and create GET request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequest("POST", fmt.Sprintf("%s/user_tokens/revoke", API), strings.NewReader(values.Encode()))
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
