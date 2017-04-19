package authsrv

type Client struct {
	// id stores the User ID.
	Id string `json:"id"`
	// name stores the name.
	Name string `json:"name,omitempty"`
	// email stores the email.
	Email string `json:"email,omitempty"`
	// client_id stores the client id
	ClientID string `json:"client_id,omitempty"`
	// client_secret stores the client secret
	ClientSecret string `json:"client_secret,omitempty"`
	// redirect_uri stores the redirect url
	RedirectURI string `json:"redirect_uri,omitempty"`
	// auth_uri stores the authorize url
	AuthURI string `json:"auth_uri,omitempty"`
	// token_uri stores the token url
	TokenURI string `json:"token_uri,omitempty"`
}
