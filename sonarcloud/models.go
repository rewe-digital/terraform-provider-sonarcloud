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
