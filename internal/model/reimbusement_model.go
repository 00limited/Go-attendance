package model

import (
	"time"
)

// ReimbursementStatus represents the status of reimbursement request
type ReimbursementStatus string

const (
	ReimbursementPending  ReimbursementStatus = "pending"
	ReimbursementApproved ReimbursementStatus = "approved"
	ReimbursementRejected ReimbursementStatus = "rejected"
	ReimbursementPaid     ReimbursementStatus = "paid"
)

// ReimbursementCategory represents the category of reimbursement
type ReimbursementCategory string

const (
	ReimbursementTravel    ReimbursementCategory = "travel"
	ReimbursementMeals     ReimbursementCategory = "meals"
	ReimbursementEquipment ReimbursementCategory = "equipment"
	ReimbursementTraining  ReimbursementCategory = "training"
	ReimbursementMedical   ReimbursementCategory = "medical"
	ReimbursementOther     ReimbursementCategory = "other"
)

// Reimbursement represents a reimbursement request made by an employee.
type Reimbursement struct {
	DefaultAttribute
	EmployeeID        uint                  `json:"employee_id" gorm:"not null;index" validate:"required"`
	ReimbursementDate time.Time             `json:"reimbursement_date" gorm:"not null;type:date;index" validate:"required"`
	Amount            float64               `json:"amount" gorm:"not null;type:decimal(12,2)" validate:"required,min=0.01,max=999999.99"`
	Category          ReimbursementCategory `json:"category" gorm:"not null;size:50;default:'other'" validate:"required,oneof=travel meals equipment training medical other"`
	Reason            string                `json:"reason" gorm:"not null;size:255" validate:"required,min=5,max=255"`
	Status            ReimbursementStatus   `json:"status" gorm:"not null;default:'pending';size:50" validate:"required,oneof=pending approved rejected paid"`
	ApprovedBy        *uint                 `json:"approved_by" gorm:"default:null"`
	ApprovedAt        *time.Time            `json:"approved_at" gorm:"default:null"`
	// Relationships
	Employee Employee  `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Approver *Employee `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy"`
}

// TableName returns the table name for the Reimbursement model.
func (Reimbursement) TableName() string {
	return "reimbursements"
}

// Approve marks the reimbursement as approved
func (r *Reimbursement) Approve(approverID uint) {
	now := time.Now()
	r.Status = ReimbursementApproved
	r.ApprovedBy = &approverID
	r.ApprovedAt = &now
}

// Reject marks the reimbursement as rejected
func (r *Reimbursement) Reject(approverID uint, reason string) {
	now := time.Now()
	r.Status = ReimbursementRejected
	r.ApprovedBy = &approverID
	r.ApprovedAt = &now
}

// MarkAsPaid marks the reimbursement as paid
func (r *Reimbursement) MarkAsPaid() {
	r.Status = ReimbursementPaid
}

// IsApproved checks if reimbursement is approved
func (r *Reimbursement) IsApproved() bool {
	return r.Status == ReimbursementApproved
}

// IsPaid checks if reimbursement is paid
func (r *Reimbursement) IsPaid() bool {
	return r.Status == ReimbursementPaid
}

// CanBeProcessedInPayroll checks if reimbursement can be included in payroll
func (r *Reimbursement) CanBeProcessedInPayroll() bool {
	return r.Status == ReimbursementApproved
}
