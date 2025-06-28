package model

import "time"

// Payslip represents a payslip record for an employee.
type Payslip struct {
	DefaultAttribute
	EmployeeID          uint      `json:"employee_id" gorm:"not null"`
	PayPeriodStart      time.Time `json:"pay_period_start" gorm:"not null"`
	PayPeriodEnd        time.Time `json:"pay_period_end" gorm:"not null"`
	BasicSalary         float64   `json:"basic_salary" gorm:"not null"`
	OvertimeHours       int       `json:"overtime_hours" gorm:"default:0"`
	OvertimeAmount      float64   `json:"overtime_amount" gorm:"default:0"`
	ReimbursementAmount float64   `json:"reimbursement_amount" gorm:"default:0"`
	TotalAmount         float64   `json:"total_amount" gorm:"not null"`
	ProcessedAt         time.Time `json:"processed_at" gorm:"not null"`
	Status              string    `json:"status" gorm:"not null;default:'processed'"` // processed, paid
	AttendanceDays      int       `json:"attendance_days" gorm:"default:0"`
}

// TableName returns the table name for the Payslip model.
func (Payslip) TableName() string {
	return "payslips"
}
