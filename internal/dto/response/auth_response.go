package response

import "time"

// LoginResponse represents the login response payload
type LoginResponse struct {
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

// UserInfo represents user information in the token response
type UserInfo struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}
