package api

// AUTOMATICALLY GENERATED, DO NOT EDIT BY HAND!

// UserGroupsAddUser: Add a user to a group.<br />'id' or 'name' must be provided.<br />Requires the following permission: 'Administer System'.
type UserGroupsAddUser struct {
	Id           string `json:"id"`           // Group id
	Login        string `json:"login"`        // User login
	Name         string `json:"name"`         // Group name
	Organization string `json:"organization"` // Key of organization
}

// UserGroupsCreate: Create a group.<br>Requires the following permission: 'Administer System'.
type UserGroupsCreate struct {
	Description  string `json:"description"`  // Description for the new group. A group description cannot be larger than 200 characters.
	Name         string `json:"name"`         // Name for the new group. A group name cannot be larger than 255 characters and must be unique. The value 'anyone' (whatever the case) is reserved and cannot be used.
	Organization string `json:"organization"` // Key of organization
}

// UserGroupsDelete: Delete a group. The default groups cannot be deleted.<br/>'id' or 'name' must be provided.<br />Requires the following permission: 'Administer System'.
type UserGroupsDelete struct {
	Id           string `json:"id"`           // Group id
	Name         string `json:"name"`         // Group name
	Organization string `json:"organization"` // Key of organization
}

// UserGroupsRemoveUser: Remove a user from a group.<br />'id' or 'name' must be provided.<br>Requires the following permission: 'Administer System'.
type UserGroupsRemoveUser struct {
	Id           string `json:"id"`           // Group id
	Login        string `json:"login"`        // User login
	Name         string `json:"name"`         // Group name
	Organization string `json:"organization"` // Key of organization
}

// UserGroupsSearch: Search for user groups.<br>Requires the following permission: 'Administer System'.
type UserGroupsSearch struct {
	F            string `json:"f"`            // Comma-separated list of the fields to be returned in response. All the fields are returned by default.
	Organization string `json:"organization"` // Key of organization
	P            string `json:"p"`            // 1-based page number
	Ps           string `json:"ps"`           // Page size. Must be greater than 0 and less or equal than 500
	Q            string `json:"q"`            // Limit search to names that contain the supplied string.
}

// UserGroupsUpdate: Update a group.<br>Requires the following permission: 'Administer System'.
type UserGroupsUpdate struct {
	Description string `json:"description"` // New optional description for the group. A group description cannot be larger than 200 characters. If value is not defined, then description is not changed.
	Id          string `json:"id"`          // Identifier of the group.
	Name        string `json:"name"`        // New optional name for the group. A group name cannot be larger than 255 characters and must be unique. Value 'anyone' (whatever the case) is reserved and cannot be used. If value is empty or not defined, then name is not changed.
}

// UserGroupsUsers: Search for users with membership information with respect to a group.<br>Requires the following permission: 'Administer System'.
type UserGroupsUsers struct {
	Id           string `json:"id"`           // Group id
	Name         string `json:"name"`         // Group name
	Organization string `json:"organization"` // Key of organization
	P            string `json:"p"`            // 1-based page number
	Ps           string `json:"ps"`           // Page size. Must be greater than 0.
	Q            string `json:"q"`            // Limit search to names or logins that contain the supplied string.
	Selected     string `json:"selected"`     // Depending on the value, show only selected items (selected=selected), deselected items (selected=deselected), or all items with their selection status (selected=all).
}
