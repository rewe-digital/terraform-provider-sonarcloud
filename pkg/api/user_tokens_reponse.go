package api

// UserTokensGenerateResponse: Response for the UserTokensGenerate request
type UserTokensGenerateResponse struct {
	Login string `json:"login,omitempty"` // User login. If not set, the token is generated for the authenticated user.
	Name  string `json:"name,omitempty"`  // Token name
	Token string `json:"token,omitempty"` // Token value
}

// UserTokensSearchResponse: Response for the UserTokensSearch request
type UserTokensSearchResponse struct {
	Login  string                         `json:"login,omitempty"`      // User login
	Tokens []UserTokenSearchResponseToken `json:"userTokens,omitempty"` // User login
}

type UserTokenSearchResponseToken struct {
	Name string `json:"name,omitempty"` // Token name
}
