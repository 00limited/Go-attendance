package model

import (
	"time"
)

// OvertimeStatus represents the status of overtime request
type OvertimeStatus string

const (
	OvertimePending  OvertimeStatus = "pending"
	OvertimeApproved OvertimeStatus = "approved"
	OvertimeRejected OvertimeStatus = "rejected"
)

// Overtime represents an overtime record for an employee.
type Overtime struct {
	DefaultAttribute
	EmployeeID   uint           `json:"employee_id" gorm:"not null;index" validate:"required"`
	OvertimeDate string         `json:"overtime_date" validate:"required"`
	StartTime    string         `json:"start_time" `
	EndTime      string         `json:"end_time"`
	Hours        int            `json:"hours" gorm:"not null" validate:"required,min=1,max=12"`
	Reason       string         `json:"reason" gorm:"not null;size:255" validate:"required,min=5,max=255"`
	Status       OvertimeStatus `json:"status" gorm:"not null;default:'pending';size:20" validate:"required,oneof=pending approved rejected"`
	ApprovedBy   *uint          `json:"approved_by" gorm:"default:null"`
	ApprovedAt   *time.Time     `json:"approved_at" gorm:"default:null"`

	// Relationships
	Employee Employee  `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Approver *Employee `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy"`
}

// TableName returns the table name for the Overtime model.
func (Overtime) TableName() string {
	return "overtimes"
}

// Approve marks the overtime as approved
func (o *Overtime) Approve(approverID uint) {
	now := time.Now()
	o.Status = OvertimeApproved
	o.ApprovedBy = &approverID
	o.ApprovedAt = &now
}

// Reject marks the overtime as rejected
func (o *Overtime) Reject(approverID uint) {
	now := time.Now()
	o.Status = OvertimeRejected
	o.ApprovedBy = &approverID
	o.ApprovedAt = &now
}

// IsApproved checks if overtime is approved
func (o *Overtime) IsApproved() bool {
	return o.Status == OvertimeApproved
}

// CalculateHoursFromTime calculates hours from start and end time
func (o *Overtime) CalculateHoursFromTime() {
	// Assuming StartTime and EndTime are in "15:04" format (HH:mm)
	const timeLayout = "15:04"
	start, err1 := time.Parse(timeLayout, o.StartTime)
	end, err2 := time.Parse(timeLayout, o.EndTime)
	if err1 == nil && err2 == nil && end.After(start) {
		duration := end.Sub(start)
		o.Hours = int(duration.Hours())
	}
}
