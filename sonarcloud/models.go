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

// Not sure what to do about actions. I haven't set them somewhere in resource_quality_gates.go, but I cannot find where that is.
// Running acceptance tests shows the error with the helpful message "unhandled unknown value"
// More info on the error here: https://github.com/hashicorp/terraform-plugin-framework/issues/191
// It may be okay to leave this commented out, as these values are not user actionable.
// type Action struct {
// 	Copy              types.Bool `tfsdk:"copy"`
// 	Delete            types.Bool `tfsdk:"delete"`
// 	ManageConditions  types.Bool `tfsdk:"manage_conditions"`
// 	Rename            types.Bool `tfsdk:"rename"`
// 	SetAsDefault      types.Bool `tfsdk:"set_as_default"`
// 	AssociateProjects types.Bool `tfsdk:"associate_projects"`
// }

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
	// Actions    Action        `tfsdk:"actions"`
	Conditions []Condition   `tfsdk:"conditions"`
	ID         types.Float64 `tfsdk:"id"`
	IsBuiltIn  types.Bool    `tfsdk:"is_built_in"`
	IsDefault  types.Bool    `tfsdk:"is_default"`
	Name       types.String  `tfsdk:"name"`
}

type QualityGates struct {
	ID           types.String  `tfsdk:"id"`
	QualityGates []QualityGate `tfsdk:"quality_gates"`
}
