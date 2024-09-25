package auth

type SessionManagerKey struct{}

type UserSession struct {
	Id       int
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     []string `json:"role"`
}

type User struct {
	Id       int
	Username string
	Email    string
	Password string
}
