package api

// UsersResponse: Response for the UserSearch request
type OrganizationsSearchMembersResponse struct {
	Paging OrganizationSearchSearchMembersResponsePaging `json:"paging,omitempty"` // Paging info of the response
	Users  []OrganizationSearchMembersResponseUsers      `json:"users,omitempty"`  // Users of the organization
}

type OrganizationSearchSearchMembersResponsePaging struct {
	PageIndex int `json:"pageIndex,omitempty"`
	PageSize  int `json:"pageSize,omitempty"`
	Total     int `json:"total,omitempty"`
}

type OrganizationSearchMembersResponseUsers struct {
	Login string `json:"login,omitempty"`
}
