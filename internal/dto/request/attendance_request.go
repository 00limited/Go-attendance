package request

// CreateAttendanceRequest represents the request body for creating an attendance record.
type CreateAttendanceRequest struct {
	EmployeeID  uint   `json:"employee_id" validate:"required"`
	Checkin     string `json:"checkin" validate:"required"`  // ISO 8601 format
	Checkout    string `json:"checkout" validate:"required"` // ISO 8601 format
	HoursWorked int    `json:"hours_worked" validate:"required"`
	Status      string `json:"status" validate:"required,oneof=present absent leave"`
}
