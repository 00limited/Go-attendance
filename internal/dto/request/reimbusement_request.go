package request

// CreateReimbusementRequest represents the request payload for creating a reimbusement.
type CreateReimbusementRequest struct {
	EmployeeID  uint    `json:"employee_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,min=0"`
	Description string  `json:"description" validate:"required"`
}
