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
	"time"
)

type UserGroupsSearchResponse struct {
	Groups []UserGroup `json:"groups"`
}

type UserGroupCreateResponse struct {
	UserGroup `json:"group"`
}

type UserGroup struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	MembersCount int    `json:"membersCount"`
	Default      bool   `json:"default"`
}

func dataSourceUserGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"groups": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"members_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	sc := m.(*SonarClient)
	req, err := sc.NewRequest("GET", fmt.Sprintf("%s/user_groups/search", API), nil)
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

	groups := &UserGroupsSearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(&groups)
	if err != nil {
		return diag.Errorf("Decode error: %+v", err)
	}

	g := asLowerCaseMap(&groups.Groups)
	if err := d.Set("groups", g); err != nil {
		return diag.Errorf("Error setting state: %+v", err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func asLowerCaseMap(ug *[]UserGroup) (entries []map[string]interface{}) {
	if ug == nil {
		return
	}

	for _, u := range *ug {
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

func appendDiagErrorFromStr(diags diag.Diagnostics, s string) diag.Diagnostics {
	return append(diags, diag.Diagnostic{
		Severity: 0,
		Summary:  s,
	})
}
