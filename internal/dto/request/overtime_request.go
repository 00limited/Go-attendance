package request

type CreateOvertimeRequest struct {
	EmployeeID uint   `json:"employee_id" validate:"required"`
	Reason     string `json:"reason" validate:"required"`
	Hours      int    `json:"hours" validate:"required,min=1"`
}
