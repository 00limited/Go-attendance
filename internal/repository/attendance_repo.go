package repository

import (
	"fmt"
	"time"

	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"gorm.io/gorm"
)

type attendance struct {
	db *gorm.DB
}

// NewAttendanceRepository creates a new instance of attendance repository.
func NewAttendanceRepository(db *gorm.DB) *attendance {
	return &attendance{db: db}
}

type AttendanceRepository interface {
	CreateAttendancePeriod(employeID uint, checkin time.Time, checkout *time.Time) (*model.Attendance, error)
	CheckinAttendancePeriod(employeID uint) (*model.Attendance, error)
	CheckOutAttendancePeriod(employeID uint) (*model.Attendance, error)
	GetTodayAttendance(employeID uint) (*model.Attendance, error)
	UpdateOrCreateAttendance(employeID uint, date time.Time, checkin time.Time, checkout *time.Time) (*model.Attendance, error)

	// Audit-enabled methods
	CheckinAttendancePeriodWithAudit(employeID uint, auditDB *middleware.AuditableDB) (*model.Attendance, error)
	CheckOutAttendancePeriodWithAudit(employeID uint, auditDB *middleware.AuditableDB) (*model.Attendance, error)
}

func (a *attendance) CreateAttendancePeriod(employeID uint, checkin time.Time, checkout *time.Time) (*model.Attendance, error) {
	// First, check if the employee exists
	var employee model.Employee
	if err := a.db.First(&employee, employeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee with ID %d not found", employeID)
		}
		return nil, err
	}

	date := time.Date(checkin.Year(), checkin.Month(), checkin.Day(), 0, 0, 0, 0, checkin.Location())

	// Check if attendance record exists for this employee and date
	var existingAttendance model.Attendance
	err := a.db.Where("employee_id = ? AND date = ?", employeID, date).First(&existingAttendance).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		// Create new attendance record
		attendance := model.Attendance{
			EmployeeID: employeID,
			Checkin:    checkin,
			Checkout:   checkout,
			Status:     "present",
			Date:       date,
		}

		// Calculate hours if checkout is provided
		attendance.CalculateHours()

		err := a.db.Create(&attendance).Error
		if err != nil {
			return nil, err
		}
		return &attendance, nil
	} else {
		// Update existing record
		existingAttendance.Checkin = checkin
		existingAttendance.Checkout = checkout
		existingAttendance.Status = "present"

		// Calculate hours
		existingAttendance.CalculateHours()

		err := a.db.Save(&existingAttendance).Error
		if err != nil {
			return nil, err
		}
		return &existingAttendance, nil
	}
}

func (a *attendance) CheckinAttendancePeriod(employeID uint) (*model.Attendance, error) {
	// Check if attendance on weekend then return error
	now := time.Now()
	dayOfWeek := now.Weekday()
	if dayOfWeek == time.Saturday || dayOfWeek == time.Sunday {
		return nil, fmt.Errorf("attendance cannot be created on weekends")
	}

	// Check if already checked in today
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var existingAttendance model.Attendance
	err := a.db.Where("employee_id = ? AND date = ?", employeID, today).First(&existingAttendance).Error

	if err == nil {
		return nil, fmt.Errorf("already checked in today")
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new attendance record
	attendance := model.Attendance{
		EmployeeID:  employeID,
		Checkin:     now,
		Status:      "present",
		Date:        today,
		HoursWorked: 0,
	}

	err = a.db.Create(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (a *attendance) CheckOutAttendancePeriod(employeID uint) (*model.Attendance, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Find today's attendance record
	var attendance model.Attendance
	err := a.db.Where("employee_id = ? AND date = ?", employeID, today).First(&attendance).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("no check-in record found for today")
	}

	if err != nil {
		return nil, err
	}

	if attendance.Checkout != nil {
		return nil, fmt.Errorf("already checked out today")
	}

	// Update with checkout time
	attendance.Checkout = &now
	attendance.CalculateHours()

	err = a.db.Save(&attendance).Error
	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

func (a *attendance) GetTodayAttendance(employeID uint) (*model.Attendance, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var attendance model.Attendance
	err := a.db.Where("employee_id = ? AND date = ?", employeID, today).First(&attendance).Error

	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

func (a *attendance) UpdateOrCreateAttendance(employeID uint, date time.Time, checkin time.Time, checkout *time.Time) (*model.Attendance, error) {
	// Normalize date (remove time component)
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Check if attendance record exists
	var existingAttendance model.Attendance
	err := a.db.Where("employee_id = ? AND date = ?", employeID, normalizedDate).First(&existingAttendance).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		// Create new record
		attendance := model.Attendance{
			EmployeeID: employeID,
			Checkin:    checkin,
			Checkout:   checkout,
			Status:     "present",
			Date:       normalizedDate,
		}
		attendance.CalculateHours()

		err := a.db.Create(&attendance).Error
		if err != nil {
			return nil, err
		}
		return &attendance, nil
	} else {
		// Update existing record
		existingAttendance.Checkin = checkin
		existingAttendance.Checkout = checkout
		existingAttendance.Status = "present"
		existingAttendance.CalculateHours()

		err := a.db.Save(&existingAttendance).Error
		if err != nil {
			return nil, err
		}
		return &existingAttendance, nil
	}
}

// CheckinAttendancePeriodWithAudit creates a check-in attendance record with audit tracking
func (a *attendance) CheckinAttendancePeriodWithAudit(employeID uint, auditDB *middleware.AuditableDB) (*model.Attendance, error) {
	// First, check if the employee exists
	var employee model.Employee
	if err := a.db.First(&employee, employeID).Error; err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	// Check if there's already an attendance record for today
	today := time.Now().Format("2006-01-02")
	var existingAttendance model.Attendance
	err := a.db.Where("employee_id = ? AND DATE(date) = ?", employeID, today).First(&existingAttendance).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == nil {
		// Record already exists, return error
		return nil, fmt.Errorf("attendance already recorded for today")
	}

	// Create new attendance record
	attendance := model.Attendance{
		EmployeeID: employeID,
		Status:     "present",
		Date:       time.Now(),
		Checkin:    time.Now(),
	}

	err = auditDB.Create(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

// CheckOutAttendancePeriodWithAudit updates attendance record with checkout time and audit tracking
func (a *attendance) CheckOutAttendancePeriodWithAudit(employeID uint, auditDB *middleware.AuditableDB) (*model.Attendance, error) {
	// Find today's attendance record
	today := time.Now().Format("2006-01-02")
	var attendance model.Attendance
	err := a.db.Where("employee_id = ? AND DATE(date) = ?", employeID, today).First(&attendance).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no check-in record found for today")
		}
		return nil, err
	}

	if attendance.Checkout != nil {
		return nil, fmt.Errorf("already checked out for today")
	}

	// Update checkout time
	now := time.Now()
	attendance.Checkout = &now
	attendance.CalculateHours()

	err = auditDB.Save(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}
