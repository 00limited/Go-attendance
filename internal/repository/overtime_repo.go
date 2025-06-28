package repository

import (
	"fmt"
	"time"

	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"gorm.io/gorm"
)

type overtime struct {
	db *gorm.DB
}

// NewOvertimeRepository creates a new instance of overtime repository.
func NewOvertimeRepository(db *gorm.DB) *overtime {
	return &overtime{db: db}
}

// GetDB returns the underlying GORM DB instance for audit functionality
func (o *overtime) GetDB() *gorm.DB {
	return o.db
}

type OvertimeRepository interface {
	CreateOvertimePeriod(employeeID uint, hours int, reason string) (*model.Overtime, error)
	CreateOvertimePeriodWithAudit(employeeID uint, hours int, reason string, auditDB *middleware.AuditableDB) (*model.Overtime, error)
	GetDB() *gorm.DB
}

func (o *overtime) CreateOvertimePeriod(employeeID uint, hours int, reason string) (*model.Overtime, error) {

	today := time.Now().Format("2006-01-02")

	// Check if the employee exists
	var employee model.Employee
	if err := o.db.First(&employee, employeeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee with ID %d not found", employeeID)
		}
		return nil, err
	}

	//check if employee is already claim overtime
	var existingOvertime model.Overtime
	err := o.db.Where("employee_id = ? AND overtime_date = ?", employee.ID, today).First(&existingOvertime).Error
	if err != nil {
		return nil, fmt.Errorf("overtime for employee with ID %d already exists for today", employeeID)
	}

	if existingOvertime.ID != 0 {
		return nil, fmt.Errorf("overtime for employee with name %s already exists for today", employee.Name)
	}

	//check if not weekend
	dayOfWeek := time.Now().Weekday()
	if dayOfWeek != time.Saturday && dayOfWeek != time.Sunday {

		// check if already checkout attendance
		var attendance model.Attendance
		err := o.db.Where("employee_id = ? AND date = ?", employee.ID, today).First(&attendance).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("attendance for employee with ID %d not found", employeeID)
			}
			return nil, err
		}

		if attendance.Checkout == nil {
			return nil, fmt.Errorf("employee with ID %d has not checked out yet", employeeID)
		}

	}

	if hours > 3 {
		return nil, fmt.Errorf("overtime hours maximum is 3 hours")
	}

	// Create new overtime record
	overtimePeriod := model.Overtime{
		OvertimeDate: today,
		EmployeeID:   employeeID,
		Hours:        hours,
		Reason:       reason,
	}

	err = o.db.Create(&overtimePeriod).Error
	if err != nil {
		return nil, err
	}

	return &overtimePeriod, nil

}

func (o *overtime) CreateOvertimePeriodWithAudit(employeeID uint, hours int, reason string, auditDB *middleware.AuditableDB) (*model.Overtime, error) {
	today := time.Now().Format("2006-01-02")

	// Check if the employee exists
	var employee model.Employee
	if err := o.db.First(&employee, employeeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee with ID %d not found", employeeID)
		}
		return nil, err
	}

	//check if employee is already claim overtime
	var existingOvertime model.Overtime
	err := o.db.Where("employee_id = ? AND overtime_date = ?", employee.ID, today).First(&existingOvertime).Error
	if err != nil {
		return nil, fmt.Errorf("overtime for employee with ID %d already exists for today", employeeID)
	}

	if existingOvertime.ID != 0 {
		return nil, fmt.Errorf("overtime for employee with name %s already exists for today", employee.Name)
	}

	//check if not weekend
	dayOfWeek := time.Now().Weekday()
	if dayOfWeek != time.Saturday && dayOfWeek != time.Sunday {

		// check if already checkout attendance
		var attendance model.Attendance
		err := o.db.Where("employee_id = ? AND date = ?", employee.ID, today).First(&attendance).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("attendance for employee with ID %d not found", employeeID)
			}
			return nil, err
		}

		if attendance.Checkout == nil {
			return nil, fmt.Errorf("employee with ID %d has not checked out yet", employeeID)
		}

	}

	if hours > 3 {
		return nil, fmt.Errorf("overtime hours maximum is 3 hours")
	}

	// Create new overtime record with audit fields
	overtimePeriod := model.Overtime{
		OvertimeDate: today,
		EmployeeID:   employeeID,
		Hours:        hours,
		Reason:       reason,
	}

	err = auditDB.Create(&overtimePeriod).Error
	if err != nil {
		return nil, err
	}

	return &overtimePeriod, nil
}
