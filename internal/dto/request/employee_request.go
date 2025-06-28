package request

type CreateEmployeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
	Active   bool   `json:"active" validate:"required"`
}
type UpdateEmployeeRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
	Active   bool   `json:"active" validate:"required"`
}
