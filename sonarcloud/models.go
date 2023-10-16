package sonarcloud

import "github.com/hashicorp/terraform-plugin-framework/types"

type Groups struct {
	ID     types.String `tfsdk:"id"`
	Groups []Group      `tfsdk:"groups"`
}

type Group struct {
	ID           types.String `tfsdk:"id"`
	Default      types.Bool   `tfsdk:"default"`
	Description  types.String `tfsdk:"description"`
	MembersCount types.Number `tfsdk:"members_count"`
	Name         types.String `tfsdk:"name"`
}

type GroupMember struct {
	ID    types.String `tfsdk:"id"`
	Group types.String `tfsdk:"group"`
	Login types.String `tfsdk:"login"`
}

type User struct {
	Login types.String `tfsdk:"login"`
	Name  types.String `tfsdk:"name"`
}

type Users struct {
	ID    types.String `tfsdk:"id"`
	Group types.String `tfsdk:"group"`
	Users []User       `tfsdk:"users"`
}

type Token struct {
	ID    types.String `tfsdk:"id"`
	Login types.String `tfsdk:"login"`
	Name  types.String `tfsdk:"name"`
	Token types.String `tfsdk:"token"`
}

type Projects struct {
	ID       types.String `tfsdk:"id"`
	Projects []Project    `tfsdk:"projects"`
}

type Project struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Key        types.String `tfsdk:"key"`
	Visibility types.String `tfsdk:"visibility"`
}

type ProjectMainBranch struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ProjectKey types.String `tfsdk:"project_key"`
}

type Condition struct {
	Error  types.String  `tfsdk:"error"`
	ID     types.Float64 `tfsdk:"id"`
	Metric types.String  `tfsdk:"metric"`
	Op     types.String  `tfsdk:"op"`
}

type Conditions struct {
	ID         types.Float64 `tfsdk:"id"`
	Conditions []Condition   `tfsdk:"condition"`
}

type QualityGate struct {
	ID         types.String  `tfsdk:"id"`
	GateId     types.Float64 `tfsdk:"gate_id"`
	Conditions []Condition   `tfsdk:"conditions"`
	IsBuiltIn  types.Bool    `tfsdk:"is_built_in"`
	IsDefault  types.Bool    `tfsdk:"is_default"`
	Name       types.String  `tfsdk:"name"`
}

type QualityGates struct {
	ID           types.String  `tfsdk:"id"`
	QualityGates []QualityGate `tfsdk:"quality_gates"`
}

type Selection struct {
	ID          types.String `tfsdk:"id"`
	GateId      types.String `tfsdk:"gate_id"`
	ProjectKeys types.Set    `tfsdk:"project_keys"`
}

type DataUserGroupPermissionsGroup struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Permissions types.Set    `tfsdk:"permissions"`
}

type DataUserGroupPermissions struct {
	ID         types.String                    `tfsdk:"id"`
	ProjectKey types.String                    `tfsdk:"project_key"`
	Groups     []DataUserGroupPermissionsGroup `tfsdk:"groups"`
}

type UserGroupPermissions struct {
	ID          types.String `tfsdk:"id"`
	ProjectKey  types.String `tfsdk:"project_key"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Permissions types.Set    `tfsdk:"permissions"`
}

type DataUserPermissionsUser struct {
	Login       types.String `tfsdk:"login"`
	Name        types.String `tfsdk:"name"`
	Permissions types.Set    `tfsdk:"permissions"`
	Avatar      types.String `tfsdk:"avatar"`
}

type DataUserPermissions struct {
	ID         types.String              `tfsdk:"id"`
	ProjectKey types.String              `tfsdk:"project_key"`
	Users      []DataUserPermissionsUser `tfsdk:"users"`
}

type UserPermissions struct {
	ID          types.String `tfsdk:"id"`
	ProjectKey  types.String `tfsdk:"project_key"`
	Login       types.String `tfsdk:"login"`
	Name        types.String `tfsdk:"name"`
	Permissions types.Set    `tfsdk:"permissions"`
	Avatar      types.String `tfsdk:"avatar"`
}

type DataProjectLinks struct {
	ID         types.String      `tfsdk:"id"`
	ProjectKey types.String      `tfsdk:"project_key"`
	Links      []DataProjectLink `tfsdk:"links"`
}

type DataProjectLink struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	Url  types.String `tfsdk:"url"`
}

type ProjectLink struct {
	ID         types.String `tfsdk:"id"`
	ProjectKey types.String `tfsdk:"project_key"`
	Name       types.String `tfsdk:"name"`
	Url        types.String `tfsdk:"url"`
}

type DataWebhooks struct {
	ID       types.String  `tfsdk:"id"`
	Project  types.String  `tfsdk:"project"`
	Webhooks []DataWebhook `tfsdk:"webhooks"`
}

type DataWebhook struct {
	Key       types.String `tfsdk:"key"`
	Name      types.String `tfsdk:"name"`
	HasSecret types.Bool   `tfsdk:"has_secret"`
	Url       types.String `tfsdk:"url"`
}

type Webhook struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Project types.String `tfsdk:"project"`
	Name    types.String `tfsdk:"name"`
	Secret  types.String `tfsdk:"secret"`
	Url     types.String `tfsdk:"url"`
}
