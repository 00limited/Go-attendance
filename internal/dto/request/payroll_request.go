package request

import "time"

// PayrollRequest represents the request to run payroll
type PayrollRequest struct {
	PayPeriodStart time.Time `json:"pay_period_start" validate:"required"`
	PayPeriodEnd   time.Time `json:"pay_period_end" validate:"required"`
	BasicSalary    float64   `json:"basic_salary" validate:"required,min=0"`
	OvertimeRate   float64   `json:"overtime_rate" validate:"required,min=0"` // Rate per hour for overtime
}

// PayrollEmployeeRequest for processing individual employee payroll
type PayrollEmployeeRequest struct {
	EmployeeID     uint      `json:"employee_id" validate:"required"`
	PayPeriodStart time.Time `json:"pay_period_start" validate:"required"`
	PayPeriodEnd   time.Time `json:"pay_period_end" validate:"required"`
	BasicSalary    float64   `json:"basic_salary" validate:"required,min=0"`
	OvertimeRate   float64   `json:"overtime_rate" validate:"required,min=0"`
}

// PayrollSummaryRequest for generating payroll summary reports
type PayrollSummaryRequest struct {
	PayPeriodStart time.Time `json:"pay_period_start" validate:"required"`
	PayPeriodEnd   time.Time `json:"pay_period_end" validate:"required"`
}
