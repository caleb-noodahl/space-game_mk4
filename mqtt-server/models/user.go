package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
}

// A profile is a public facing user object
type Profile struct {
	Username string `json:"username"`
	ClientID string `json:"client_id"`
}
