package sonarcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/iancoleman/strcase"
	"strconv"
	"terraform-provider-sonarcloud/pkg/api"
	"time"
)

func dataSourceOrganizationMembers() *schema.Resource {
	return &schema.Resource{
		Description: "Data source that retrieves a list of users of the configured organization.",
		ReadContext: dataSourceOrganizationMembersRead,
		Schema: map[string]*schema.Schema{
			"users": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The users that are a member of this organization.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"login": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The login of this user",
						},
					},
				},
			},
		},
	}
}

func dataSourceOrganizationMembersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO: check paging of result to see if we need to fetch more pages
	sc := m.(*SonarClient)
	req, err := sc.NewRequestWithParameters("GET", fmt.Sprintf("%s/organizations/search_members", API), "p", "1", "ps", "500")
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := sc.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return diagErrorResponse(resp, diags)
	}

	members := &api.OrganizationsSearchMembersResponse{}
	err = json.NewDecoder(resp.Body).Decode(&members)
	if err != nil {
		return diag.Errorf("Decode error: %+v", err)
	}

	users := usersAsLowerCaseMap(&members.Users)
	if err := d.Set("users", users); err != nil {
		return diag.Errorf("Error setting state: %+v", err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func usersAsLowerCaseMap(users *[]api.OrganizationSearchMembersResponseUsers) (entries []map[string]interface{}) {
	if users == nil {
		return
	}

	for _, u := range *users {
		m := make(map[string]interface{})
		s := structs.New(u)
		for _, f := range s.Fields() {
			if f.IsExported() {
				name := strcase.ToSnake(f.Name())
				m[name] = f.Value()
			}
		}
		entries = append(entries, m)
	}

	return entries
}
