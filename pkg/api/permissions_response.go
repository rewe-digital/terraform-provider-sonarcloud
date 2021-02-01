package api

// PermissionsGroupsSearchResponse: Response for PermissionsGroupsSearch requests
type PermissionsGroupsSearchResponse struct {
	Groups []PermissionsGroup `json:"groups"`
}

// PermissionsGroup: Single group of the PermissionsGroupsSearchResponse
type PermissionsGroup struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}
