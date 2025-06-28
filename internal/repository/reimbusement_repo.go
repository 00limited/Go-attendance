package repository

import (
	"fmt"
	"time"

	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"

	"gorm.io/gorm"
)

type reimbusement struct {
	db *gorm.DB
}

// NewReimbusementRepository creates a new reimbusement repository
func NewReimbusementRepository(db *gorm.DB) *reimbusement {
	return &reimbusement{db: db}
}

// GetDB returns the underlying GORM DB instance for audit functionality
func (r *reimbusement) GetDB() *gorm.DB {
	return r.db
}

// ReimbusementRepository defines the interface for reimbusement repository
type ReimbusementRepository interface {
	CreateReimbusement(employeeID uint, amount float64, description string) (*model.Reimbursement, error)
	CreateReimbusementWithAudit(employeeID uint, amount float64, description string, auditDB *middleware.AuditableDB) (*model.Reimbursement, error)
	GetDB() *gorm.DB
}

// CreateReimbusement creates a new reimbusement record
func (r *reimbusement) CreateReimbusement(employeeID uint, amount float64, description string) (*model.Reimbursement, error) {

	timeNow := time.Now()
	// Check if the employee exists
	var employee model.Employee
	if err := r.db.First(&employee, employeeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee with ID %d not found", employeeID)
		}
		return nil, err
	}

	//check if employee already claim reimbusement
	var existingReimbusement model.Reimbursement
	err := r.db.Where("employee_id = ? AND DATE(reimbursement_date) = ?", employee.ID, timeNow.Format("2006-01-02")).Find(&existingReimbusement).Error
	if err != nil {
		return nil, fmt.Errorf("reimbusement for employee with ID %d already exists for today", employeeID)
	}
	if existingReimbusement.ID != 0 {
		return nil, fmt.Errorf("reimbusement for employee with name %s already claim for today", employee.Name)
	}

	// Create the reimbusement record
	reimbusementRecord := model.Reimbursement{
		EmployeeID:        employeeID,
		Amount:            amount,
		Reason:            description,
		ReimbursementDate: timeNow,
	}
	err = r.db.Create(&reimbusementRecord).Error
	if err != nil {
		return nil, err
	}

	return &reimbusementRecord, nil
}

// CreateReimbusementWithAudit creates a new reimbusement record with audit trail
func (r *reimbusement) CreateReimbusementWithAudit(employeeID uint, amount float64, description string, auditDB *middleware.AuditableDB) (*model.Reimbursement, error) {
	timeNow := time.Now()
	// Check if the employee exists
	var employee model.Employee
	if err := r.db.First(&employee, employeeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee with ID %d not found", employeeID)
		}
		return nil, err
	}

	//check if employee already claim reimbusement
	var existingReimbusement model.Reimbursement
	err := r.db.Where("employee_id = ? AND DATE(reimbursement_date) = ?", employee.ID, timeNow.Format("2006-01-02")).Find(&existingReimbusement).Error
	if err != nil {
		return nil, fmt.Errorf("reimbusement for employee with ID %d already exists for today", employeeID)
	}
	if existingReimbusement.ID != 0 {
		return nil, fmt.Errorf("reimbusement for employee with name %s already claim for today", employee.Name)
	}

	// Create the reimbusement record with audit fields
	reimbusementRecord := model.Reimbursement{
		EmployeeID:        employeeID,
		Amount:            amount,
		Reason:            description,
		ReimbursementDate: timeNow,
	}
	err = auditDB.Create(&reimbusementRecord).Error
	if err != nil {
		return nil, err
	}

	return &reimbusementRecord, nil
}
