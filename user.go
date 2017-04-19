package authsrv

type User struct {
	// id stores the User ID.
	Id string `json:"id"`
	// name stores the name.
	Name string `json:"name,omitempty"`
	// email stores the email.
	Email string `json:"email,omitempty"`
	// client_id stores the client id
	ClientID string `json:"client_id,omitempty"`
}
