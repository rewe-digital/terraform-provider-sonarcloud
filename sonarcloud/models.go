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
