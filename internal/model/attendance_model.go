package model

import (
	"time"
)

// Attendance represents an attendance record for an employee.
type Attendance struct {
	DefaultAttribute
	EmployeeID  uint       `json:"employee_id" gorm:"not null;index" validate:"required"`
	Checkin     time.Time  `json:"checkin" gorm:"not null" validate:"required"`
	Checkout    *time.Time `json:"checkout" gorm:"default:null"`
	HoursWorked int        `json:"hours_worked" gorm:"not null;default:0" validate:"min=0,max=24"`
	Status      string     `json:"status" gorm:"not null;size:20;check:status IN ('present','absent','leave','holiday')" validate:"required,oneof=present absent leave holiday"`
	Date        time.Time  `json:"date" gorm:"not null;type:date;index" validate:"required"`

	// Relationship
	Employee Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
}

// TableName returns the table name for the Attendance model.
func (Attendance) TableName() string {
	return "attendances"
}

// CalculateHours calculates hours worked between checkin and checkout
func (a *Attendance) CalculateHours() {
	if a.Checkout != nil && a.Checkout.After(a.Checkin) {
		duration := a.Checkout.Sub(a.Checkin)
		a.HoursWorked = int(duration.Hours())
	} else {
		a.HoursWorked = 0
	}
}

// IsPresent checks if the employee was present on this day
func (a *Attendance) IsPresent() bool {
	return a.Status == "present"
}

// IsComplete checks if attendance record is complete (has checkout)
func (a *Attendance) IsComplete() bool {
	return a.Checkout != nil
}
