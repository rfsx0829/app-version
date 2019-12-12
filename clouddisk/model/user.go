package model

// User means a user
type User struct {
	UID      int    `json:"uid"`
	Uname    string `json:"uname"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
