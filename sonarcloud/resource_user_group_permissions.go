package sonarcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"sync"
	"terraform-provider-sonarcloud/pkg/api"
	"terraform-provider-sonarcloud/pkg/collection"
	"time"
)

func resourceUserGroupPermissions() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource manages the permissions of a user group for the whole organization or a specific project.",
		CreateContext: resourcePermissionCreate,
		ReadContext:   resourcePermissionRead,
		UpdateContext: resourcePermissionUpdate,
		DeleteContext: resourcePermissionDelete,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The key of the project to restrict the permissions to.",
				ForceNew:    true,
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User group to set the permissions for.",
				ForceNew:    true,
			},
			"permissions": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of permissions to grant.",
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						var diags diag.Diagnostics

						allowed := []string{
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
						}
						ok := false
						for _, v := range allowed {
							if v == i.(string) {
								ok = true
								break
							}
						}
						if !ok {
							return diag.Errorf("unsupported permission '%s'", i.(string))
						}

						return diags
					},
				},
			},
		},
	}
}

func resourcePermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get all resource values
	project := d.Get("project").(string)
	group := d.Get("group").(string)
	permissions := d.Get("permissions").([]interface{})

	// Get client
	sc := m.(*SonarClient)

	// Add permissions one by one
	wg := sync.WaitGroup{}
	for _, permission := range permissions {
		permission := permission

		go func() {
			wg.Add(1)
			defer wg.Done()

			if diagnostics := addPermission(sc, project, group, permission); diagnostics != nil {
				diags = append(diags, diagnostics...)
				return
			}
		}()
	}

	// FIXME: use StateChangeConf to handle eventual consistency
	time.Sleep(1 * time.Second)

	// Async wait for all requests to be done
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	// Set ID on success and return error diag on timeout
	select {
	case <-c:
		// Set resource id
		if diags == nil || !diags.HasError() {
			d.SetId(group + project)
			return resourcePermissionRead(ctx, d, m)
		}
	case <-time.After(30 * time.Second):
		return diag.Errorf("requests timed out")
	}

	return diags
}

func resourcePermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get all resource values needed to read the remote resource
	project := d.Get("project").(string)
	group := d.Get("group").(string)
	permissions := d.Get("permissions").([]interface{})
	var params []string
	if project != "" {
		params = []string{"projectKey", project}
	}

	// Cast m to SonarClient and create GET request for URI with encoded values
	sc := m.(*SonarClient)
	req, err := sc.NewRequestWithParameters("GET", fmt.Sprintf("%s/permissions/groups", API), params...)
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
	searchResponse := &api.PermissionsGroupsSearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		fmt.Println(err)
		return diag.Errorf("Decode error: %+v", err)
	}

	// Check if the resource exists in the list of retrieved resources
	// TODO: anti-corruption layer that hides this implementation detail
	groupFound := false
	for _, g := range searchResponse.Groups {
		if g.Name == group {
			groupFound = true

			gp := collection.ToInterfaceSlice(g.Permissions)
			actualPermissions := collection.Ordered(gp, permissions)

			if err := d.Set("permissions", actualPermissions); err != nil {
				return diag.FromErr(err)
			}
			break
		}
	}

	// Unset the id if the resource has not been found
	if !groupFound {
		d.SetId("")
	}

	return diags
}

func resourcePermissionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check if any of the resource's values have been changed
	if !d.HasChange("permissions") {
		return resourcePermissionRead(ctx, d, m)
	}

	// Get all resource values needed to read the remote resource
	project := d.Get("project").(string)
	group := d.Get("group").(string)

	// Find which permissions to add and which to remove
	o, n := d.GetChange("permissions")
	added, removed := collection.Diff(o.([]interface{}), n.([]interface{}))

	// Get client and prepare synchronization
	sc := m.(*SonarClient)
	wg := sync.WaitGroup{}

	// Remove permissions
	for _, permission := range removed {
		permission := permission

		go func() {
			wg.Add(1)
			defer wg.Done()

			if diagnostics := removePermission(sc, project, group, permission); diagnostics != nil {
				diags = append(diags, diagnostics...)
				return
			}
		}()
	}

	// Add permissions
	for _, permission := range added {
		permission := permission

		go func() {
			wg.Add(1)
			defer wg.Done()

			if diagnostics := addPermission(sc, project, group, permission); diagnostics != nil {
				diags = append(diags, diagnostics...)
				return
			}
		}()
	}

	// FIXME: use StateChangeConf to handle eventual consistency
	time.Sleep(1 * time.Second)

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		if diags == nil || !diags.HasError() {
			return resourcePermissionRead(ctx, d, m)
		}
	case <-time.After(30 * time.Second):
		return diag.Errorf("requests timed out")
	}

	return diags
}

func resourcePermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get all resource values
	project := d.Get("project").(string)
	group := d.Get("group").(string)
	permissions := d.Get("permissions").([]interface{})

	// Get client
	sc := m.(*SonarClient)

	wg := sync.WaitGroup{}
	for _, permission := range permissions {
		permission := permission

		go func() {
			wg.Add(1)
			defer wg.Done()

			if diagnostics := removePermission(sc, project, group, permission); diagnostics != nil {
				diags = append(diags, diagnostics...)
				return
			}
		}()
	}

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		// Nothing to do, the ID will be unset depending on the value of diags
		break
	case <-time.After(30 * time.Second):
		return diag.Errorf("requests timed out")
	}

	return diags
}

func addPermission(sc *SonarClient, project string, group string, permission interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Fill in api action struct
	create := api.PermissionsAddGroup{
		ProjectKey: project,
		GroupName:  group,
		Permission: permission.(string),
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&create)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := sc.NewRequest("POST", fmt.Sprintf("%s/permissions/add_group", API), strings.NewReader(values.Encode()))
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

	return diags
}

func removePermission(sc *SonarClient, project string, group string, permission interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Fill in api action struct
	create := api.PermissionsRemoveGroup{
		ProjectKey: project,
		GroupName:  group,
		Permission: permission.(string),
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&create)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := sc.NewRequest("POST", fmt.Sprintf("%s/permissions/remove_group", API), strings.NewReader(values.Encode()))
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

	return diags
}
