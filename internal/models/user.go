package models

type UserType string

const (
	UserTypeAdmin  UserType = "admin"
	UserTypeNormal UserType = "normal"
)

type User struct {
	ID        string   `json:"id"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	UserType  UserType `json:"user_type"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
}

type CurrentUser struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	UserType UserType `json:"user_type"`
}
