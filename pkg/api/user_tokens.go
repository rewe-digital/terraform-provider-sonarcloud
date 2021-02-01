package api

// AUTOMATICALLY GENERATED, DO NOT EDIT BY HAND!

// UserTokensGenerate: Generate a user access token. <br />Please keep your tokens secret. They enable to authenticate and analyze projects.<br />It requires administration permissions to specify a 'login' and generate a token for another user. Otherwise, a token is generated for the current user.
type UserTokensGenerate struct {
	Login string `form:"login,omitempty"` // User login. If not set, the token is generated for the authenticated user.
	Name  string `form:"name,omitempty"`  // Token name
}

// UserTokensRevoke: Revoke a user access token. <br/>It requires administration permissions to specify a 'login' and revoke a token for another user. Otherwise, the token for the current user is revoked.
type UserTokensRevoke struct {
	Login string `form:"login,omitempty"` // User login
	Name  string `form:"name,omitempty"`  // Token name
}

// UserTokensSearch: List the access tokens of a user.<br>The login must exist and active.<br>Field 'lastConnectionDate' is only updated every hour, so it may not be accurate, for instance when a user is using a token many times in less than one hour.<br/It requires administration permissions to specify a 'login' and list the tokens of another user. Otherwise, tokens for the current user are listed.
type UserTokensSearch struct {
	Login string `form:"login,omitempty"` // User login
}
