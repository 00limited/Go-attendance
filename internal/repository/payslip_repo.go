package repository

import (
	"time"

	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"gorm.io/gorm"
)

type payslip struct {
	db *gorm.DB
}

// NewPayslipRepository creates a new instance of payslip repository.
func NewPayslipRepository(db *gorm.DB) *payslip {
	return &payslip{db: db}
}

// GetDB returns the underlying GORM DB instance for audit functionality
func (p *payslip) GetDB() *gorm.DB {
	return p.db
}

type PayslipRepository interface {
	CreatePayslip(payslipData *model.Payslip) (*model.Payslip, error)
	CreatePayslipWithAudit(payslipData *model.Payslip, auditDB *middleware.AuditableDB) (*model.Payslip, error)
	GetPayslipByEmployeeAndPeriod(employeeID uint, startDate time.Time, endDate time.Time) (*model.Payslip, error)
	GetPayslipByID(payslipID uint) (*model.Payslip, error)
	GetPayslipsByEmployee(employeeID uint) ([]model.Payslip, error)
	GetPayslipsByPeriod(startDate time.Time, endDate time.Time) ([]model.Payslip, error)
	CheckPayslipExists(employeeID uint, startDate time.Time, endDate time.Time) (bool, error)
	GetAttendanceForPeriod(employeeID uint, startDate time.Time, endDate time.Time) ([]model.Attendance, error)
	GetOvertimeForPeriod(employeeID uint, startDate string, endDate string) ([]model.Overtime, error)
	GetApprovedReimbursementsForPeriod(employeeID uint, startDate, endDate time.Time) ([]model.Reimbursement, error)
	GetEmployeeByID(employeeID uint) (*model.Employee, error)
	GetDB() *gorm.DB
}

func (p *payslip) CreatePayslip(payslipData *model.Payslip) (*model.Payslip, error) {
	err := p.db.Create(payslipData).Error
	if err != nil {
		return nil, err
	}
	return payslipData, nil
}

func (p *payslip) CreatePayslipWithAudit(payslipData *model.Payslip, auditDB *middleware.AuditableDB) (*model.Payslip, error) {
	err := auditDB.Create(payslipData).Error
	if err != nil {
		return nil, err
	}
	return payslipData, nil
}

func (p *payslip) GetPayslipByEmployeeAndPeriod(employeeID uint, startDate time.Time, endDate time.Time) (*model.Payslip, error) {
	var payslip model.Payslip
	err := p.db.Where("employee_id = ? AND pay_period_start = ? AND pay_period_end = ?",
		employeeID, startDate, endDate).First(&payslip).Error
	if err != nil {
		return nil, err
	}
	return &payslip, nil
}

func (p *payslip) GetPayslipByID(payslipID uint) (*model.Payslip, error) {
	var payslip model.Payslip
	err := p.db.Where("id = ?", payslipID).First(&payslip).Error
	if err != nil {
		return nil, err
	}
	return &payslip, nil
}

func (p *payslip) GetPayslipsByEmployee(employeeID uint) ([]model.Payslip, error) {
	var payslips []model.Payslip
	err := p.db.Where("employee_id = ?", employeeID).Order("pay_period_start DESC").Find(&payslips).Error
	if err != nil {
		return nil, err
	}
	return payslips, nil
}

func (p *payslip) GetPayslipsByPeriod(startDate time.Time, endDate time.Time) ([]model.Payslip, error) {
	var payslips []model.Payslip
	err := p.db.Where("pay_period_start >= ? AND pay_period_end <= ?", startDate, endDate).Find(&payslips).Error
	if err != nil {
		return nil, err
	}
	return payslips, nil
}

func (p *payslip) CheckPayslipExists(employeeID uint, startDate time.Time, endDate time.Time) (bool, error) {
	var count int64
	err := p.db.Model(&model.Payslip{}).Where("employee_id = ? AND pay_period_start = ? AND pay_period_end = ?",
		employeeID, startDate, endDate).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (p *payslip) GetAttendanceForPeriod(employeeID uint, startDate time.Time, endDate time.Time) ([]model.Attendance, error) {
	var attendances []model.Attendance
	err := p.db.Where("employee_id = ? AND date >= ? AND date <= ?",
		employeeID, startDate, endDate).Find(&attendances).Error
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

func (p *payslip) GetOvertimeForPeriod(employeeID uint, startDate string, endDate string) ([]model.Overtime, error) {
	var overtimes []model.Overtime
	err := p.db.Debug().Where("employee_id = ? AND overtime_date >= ? AND overtime_date <= ? AND status = ?",
		employeeID, startDate, endDate, "approved").Find(&overtimes).Error
	if err != nil {
		return nil, err
	}
	return overtimes, nil
}

func (p *payslip) GetApprovedReimbursementsForPeriod(employeeID uint, startDate time.Time, endDate time.Time) ([]model.Reimbursement, error) {
	var reimbursements []model.Reimbursement
	err := p.db.Where("employee_id = ? AND reimbursement_date >= ? AND reimbursement_date <= ? AND status = ?",
		employeeID, startDate, endDate, "approved").Find(&reimbursements).Error
	if err != nil {
		return nil, err
	}
	return reimbursements, nil
}

func (p *payslip) GetEmployeeByID(employeeID uint) (*model.Employee, error) {
	var employee model.Employee
	err := p.db.Where("id = ?", employeeID).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}
