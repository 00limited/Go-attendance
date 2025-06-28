package request

// LoginRequest represents the login request payload
type LoginRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Password string `json:"password" validate:"required,min=6"`
}
