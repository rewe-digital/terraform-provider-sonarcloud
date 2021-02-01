package api

// PermissionsGroupSearch: Search for permissions<br />Requires the permission 'Administer' on the specified project.
type PermissionsGroupsSearch struct {
	ProjectKey string `form:"projectKey,omitempty"` // Show permissions for this specific project.
	Q          string `form:"q,omitempty"`          // Limit search to names that contain the supplied string.
}
