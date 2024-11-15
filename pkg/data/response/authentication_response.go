package response

type LoginResponse struct {
	TokenType string `json:"token_type"`
	Token     string `json:"token"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
