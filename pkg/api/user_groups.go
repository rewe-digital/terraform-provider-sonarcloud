package api

// AUTOMATICALLY GENERATED, DO NOT EDIT BY HAND!

// UserGroupsAddUser: Add a user to a group.<br />'id' or 'name' must be provided.<br />Requires the following permission: 'Administer System'.
type UserGroupsAddUser struct {
	Id           string `form:"id,omitempty"`           // Group id
	Login        string `form:"login,omitempty"`        // User login
	Name         string `form:"name,omitempty"`         // Group name
	Organization string `form:"organization,omitempty"` // Key of organization
}

// UserGroupsCreate: Create a group.<br>Requires the following permission: 'Administer System'.
type UserGroupsCreate struct {
	Description  string `form:"description,omitempty"`  // Description for the new group. A group description cannot be larger than 200 characters.
	Name         string `form:"name,omitempty"`         // Name for the new group. A group name cannot be larger than 255 characters and must be unique. The value 'anyone' (whatever the case) is reserved and cannot be used.
	Organization string `form:"organization,omitempty"` // Key of organization
}

// UserGroupsDelete: Delete a group. The default groups cannot be deleted.<br/>'id' or 'name' must be provided.<br />Requires the following permission: 'Administer System'.
type UserGroupsDelete struct {
	Id           string `form:"id,omitempty"`           // Group id
	Name         string `form:"name,omitempty"`         // Group name
	Organization string `form:"organization,omitempty"` // Key of organization
}

// UserGroupsRemoveUser: Remove a user from a group.<br />'id' or 'name' must be provided.<br>Requires the following permission: 'Administer System'.
type UserGroupsRemoveUser struct {
	Id           string `form:"id,omitempty"`           // Group id
	Login        string `form:"login,omitempty"`        // User login
	Name         string `form:"name,omitempty"`         // Group name
	Organization string `form:"organization,omitempty"` // Key of organization
}

// UserGroupsSearch: Search for user groups.<br>Requires the following permission: 'Administer System'.
type UserGroupsSearch struct {
	F            string `form:"f,omitempty"`            // Comma-separated list of the fields to be returned in response. All the fields are returned by default.
	Organization string `form:"organization,omitempty"` // Key of organization
	P            string `form:"p,omitempty"`            // 1-based page number
	Ps           string `form:"ps,omitempty"`           // Page size. Must be greater than 0 and less or equal than 500
	Q            string `form:"q,omitempty"`            // Limit search to names that contain the supplied string.
}

// UserGroupsUpdate: Update a group.<br>Requires the following permission: 'Administer System'.
type UserGroupsUpdate struct {
	Description string `form:"description,omitempty"` // New optional description for the group. A group description cannot be larger than 200 characters. If value is not defined, then description is not changed.
	Id          string `form:"id,omitempty"`          // Identifier of the group.
	Name        string `form:"name,omitempty"`        // New optional name for the group. A group name cannot be larger than 255 characters and must be unique. Value 'anyone' (whatever the case) is reserved and cannot be used. If value is empty or not defined, then name is not changed.
}

// UserGroupsUsers: Search for users with membership information with respect to a group.<br>Requires the following permission: 'Administer System'.
type UserGroupsUsers struct {
	Id           string `form:"id,omitempty"`           // Group id
	Name         string `form:"name,omitempty"`         // Group name
	Organization string `form:"organization,omitempty"` // Key of organization
	P            string `form:"p,omitempty"`            // 1-based page number
	Ps           string `form:"ps,omitempty"`           // Page size. Must be greater than 0.
	Q            string `form:"q,omitempty"`            // Limit search to names or logins that contain the supplied string.
	Selected     string `form:"selected,omitempty"`     // Depending on the value, show only selected items (selected=selected), deselected items (selected=deselected), or all items with their selection status (selected=all).
}
