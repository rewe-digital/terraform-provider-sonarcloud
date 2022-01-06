package sonarcloud

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/projects"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/user_groups"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/user_tokens"
	"math/big"
)

// changedAttrs returns a map where the keys are the names of all the attributes that were changed
// Note that the name is not the full path, but only the AttributeName of the last path step.
func changedAttrs(req tfsdk.UpdateResourceRequest, diags diag.Diagnostics) map[string]struct{} {
	diffs, err := req.Plan.Raw.Diff(req.State.Raw)
	if err != nil {
		diags.AddError(
			"Could not diff plan with state",
			"This should not happen and is an error in the provider.",
		)
	}

	changes := make(map[string]struct{})
	for _, diff := range diffs {
		steps := diff.Path.Steps()
		index := len(steps) - 1

		if !diff.Value1.Equal(*diff.Value2) {
			attr := steps[index].(tftypes.AttributeName)
			changes[string(attr)] = struct{}{}
		}
	}
	return changes
}

// findGroup returns the group with the given name if it exists in the response
func findGroup(response *user_groups.SearchResponseAll, name string) (Group, bool) {
	var result Group
	ok := false
	for _, g := range response.Groups {
		if g.Name == name {
			result = Group{
				ID:           types.String{Value: big.NewFloat(g.Id).String()},
				Default:      types.Bool{Value: g.Default},
				Description:  types.String{Value: g.Description},
				MembersCount: types.Number{Value: big.NewFloat(g.MembersCount)},
				Name:         types.String{Value: g.Name},
			}
			ok = true
			break
		}
	}
	return result, ok
}

// findGroup returns the group member with the given login if it exists in the response
func findGroupMember(response *user_groups.UsersResponseAll, group string, login string) (GroupMember, bool) {
	var result GroupMember
	ok := false
	for _, u := range response.Users {
		if u.Login == login {
			result = GroupMember{
				Group: types.String{Value: group},
				Login: types.String{Value: login},
			}
			ok = true
			break
		}
	}
	return result, ok
}

// tokenExists returns whether a token with the given name exists in the response
func tokenExists(response *user_tokens.SearchResponse, name string) bool {
	for _, t := range response.UserTokens {
		if t.Name == name {
			return true
		}
	}
	return false
}

// findProject returns the project with the given key if it exists in the response
func findProject(response *projects.SearchResponseAll, key string) (Project, bool) {
	var result Project
	ok := false
	for _, p := range response.Components {
		if p.Key == key {
			result = Project{
				ID:         types.String{Value: p.Key},
				Name:       types.String{Value: p.Name},
				Key:        types.String{Value: p.Key},
				Visibility: types.String{Value: p.Visibility},
			}
			ok = true
			break
		}
	}
	return result, ok
}
