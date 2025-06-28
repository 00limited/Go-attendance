package usecases

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"gorm.io/gorm"
)

// Mock repositories
type MockPayslipRepository struct {
	mock.Mock
}

func (m *MockPayslipRepository) CheckPayslipExists(employeeID uint, startDate, endDate time.Time) (bool, error) {
	args := m.Called(employeeID, startDate, endDate)
	return args.Bool(0), args.Error(1)
}

func (m *MockPayslipRepository) GetAttendanceForPeriod(employeeID uint, startDate, endDate time.Time) ([]model.Attendance, error) {
	args := m.Called(employeeID, startDate, endDate)
	return args.Get(0).([]model.Attendance), args.Error(1)
}

func (m *MockPayslipRepository) GetOvertimeForPeriod(employeeID uint, startDate, endDate string) ([]model.Overtime, error) {
	args := m.Called(employeeID, startDate, endDate)
	return args.Get(0).([]model.Overtime), args.Error(1)
}

func (m *MockPayslipRepository) GetApprovedReimbursementsForPeriod(employeeID uint, startDate, endDate time.Time) ([]model.Reimbursement, error) {
	args := m.Called(employeeID, startDate, endDate)
	return args.Get(0).([]model.Reimbursement), args.Error(1)
}

func (m *MockPayslipRepository) CreatePayslip(payslip *model.Payslip) (*model.Payslip, error) {
	args := m.Called(payslip)
	return args.Get(0).(*model.Payslip), args.Error(1)
}

func (m *MockPayslipRepository) CreatePayslipWithAudit(payslip *model.Payslip, auditDB *middleware.AuditableDB) (*model.Payslip, error) {
	args := m.Called(payslip, auditDB)
	return args.Get(0).(*model.Payslip), args.Error(1)
}

func (m *MockPayslipRepository) GetEmployeeByID(id uint) (*model.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Employee), args.Error(1)
}

// Additional methods to satisfy PayslipRepository interface
func (m *MockPayslipRepository) GetPayslipByEmployeeAndPeriod(employeeID uint, startDate, endDate time.Time) (*model.Payslip, error) {
	args := m.Called(employeeID, startDate, endDate)
	return args.Get(0).(*model.Payslip), args.Error(1)
}

func (m *MockPayslipRepository) GetPayslipByID(id uint) (*model.Payslip, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Payslip), args.Error(1)
}

func (m *MockPayslipRepository) GetPayslipsByEmployee(employeeID uint) ([]model.Payslip, error) {
	args := m.Called(employeeID)
	return args.Get(0).([]model.Payslip), args.Error(1)
}

func (m *MockPayslipRepository) GetPayslipsByPeriod(startDate, endDate time.Time) ([]model.Payslip, error) {
	args := m.Called(startDate, endDate)
	return args.Get(0).([]model.Payslip), args.Error(1)
}

func (m *MockPayslipRepository) GetDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) GetAllActiveEmployees() ([]model.Employee, error) {
	args := m.Called()
	return args.Get(0).([]model.Employee), args.Error(1)
}

// Additional methods to satisfy EmployeeRepository interface
func (m *MockEmployeeRepository) CreateEmployee(req request.CreateEmployeeRequest) (*model.Employee, error) {
	args := m.Called(req)
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetAllEmployees() ([]model.Employee, error) {
	args := m.Called()
	return args.Get(0).([]model.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) UpdateEmployee(employeeID string, req request.UpdateEmployeeRequest) (*model.Employee, error) {
	args := m.Called(employeeID, req)
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) DeleteEmployee(employeeID string) error {
	args := m.Called(employeeID)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetEmployeeByID(id uint) (*model.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetEmployeeByName(name string) (*model.Employee, error) {
	args := m.Called(name)
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) CreateEmployeeWithAudit(req request.CreateEmployeeRequest, auditDB *middleware.AuditableDB) (*model.Employee, error) {
	args := m.Called(req, auditDB)
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) UpdateEmployeeWithAudit(employeeID string, req request.UpdateEmployeeRequest, auditDB *middleware.AuditableDB) (*model.Employee, error) {
	args := m.Called(employeeID, req, auditDB)
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) DeleteEmployeeWithAudit(employeeID string, auditDB *middleware.AuditableDB) error {
	args := m.Called(employeeID, auditDB)
	return args.Error(0)
}

// Test data creation helpers
func createTestPayrollRequest() request.PayrollRequest {
	return request.PayrollRequest{
		PayPeriodStart: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		BasicSalary:    5000.0,
		OvertimeRate:   50.0,
	}
}

func createTestAttendances() []model.Attendance {
	checkout1 := time.Date(2025, 6, 1, 17, 0, 0, 0, time.UTC)
	checkout2 := time.Date(2025, 6, 2, 17, 0, 0, 0, time.UTC)
	return []model.Attendance{
		{
			EmployeeID:  1,
			Date:        time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			Checkin:     time.Date(2025, 6, 1, 9, 0, 0, 0, time.UTC),
			Checkout:    &checkout1,
			HoursWorked: 8,
			Status:      "present",
		},
		{
			EmployeeID:  1,
			Date:        time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC),
			Checkin:     time.Date(2025, 6, 2, 9, 0, 0, 0, time.UTC),
			Checkout:    &checkout2,
			HoursWorked: 8,
			Status:      "present",
		},
	}
}

func createTestOvertimes() []model.Overtime {
	return []model.Overtime{
		{
			EmployeeID:   1,
			OvertimeDate: "2025-06-01",
			Hours:        2,
			Reason:       "Project deadline",
			Status:       "approved",
		},
		{
			EmployeeID:   1,
			OvertimeDate: "2025-06-02",
			Hours:        3,
			Reason:       "Extra work",
			Status:       "approved",
		},
	}
}

func createTestReimbursements() []model.Reimbursement {
	return []model.Reimbursement{
		{
			EmployeeID:        1,
			ReimbursementDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			Amount:            100.0,
			Reason:            "Travel expenses",
			Status:            model.ReimbursementApproved,
		},
		{
			EmployeeID:        1,
			ReimbursementDate: time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC),
			Amount:            50.0,
			Reason:            "Meal allowance",
			Status:            model.ReimbursementApproved,
		},
	}
}

func createTestEmployees() []model.Employee {
	return []model.Employee{
		{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			Name:             "John Doe",
			Role:             "admin",
			Active:           true,
		},
		{
			DefaultAttribute: model.DefaultAttribute{ID: 2},
			Name:             "Jane Smith",
			Role:             "employee",
			Active:           true,
		},
	}
}

// Test ProcessEmployeePayroll - Success
func TestProcessEmployeePayroll_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()

	// Mock expectations
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)

	expectedPayslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      req.PayPeriodStart,
		PayPeriodEnd:        req.PayPeriodEnd,
		BasicSalary:         5000.0,
		OvertimeHours:       5,      // 2 + 3
		OvertimeAmount:      250.0,  // 5 * 50.0
		ReimbursementAmount: 150.0,  // 100.0 + 50.0
		TotalAmount:         5400.0, // 5000 + 250 + 150
		Status:              "processed",
		AttendanceDays:      2,
	}

	mockPayslipRepo.On("CreatePayslip", mock.MatchedBy(func(p *model.Payslip) bool {
		return p.EmployeeID == 1 &&
			p.BasicSalary == 5000.0 &&
			p.OvertimeHours == 5 &&
			p.OvertimeAmount == 250.0 &&
			p.ReimbursementAmount == 150.0 &&
			p.TotalAmount == 5400.0 &&
			p.AttendanceDays == 2 &&
			p.Status == "processed"
	})).Return(expectedPayslip, nil)

	// Execute
	result, err := usecase.ProcessEmployeePayroll(1, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.EmployeeID)
	assert.Equal(t, 5000.0, result.BasicSalary)
	assert.Equal(t, 5, result.OvertimeHours)
	assert.Equal(t, 250.0, result.OvertimeAmount)
	assert.Equal(t, 150.0, result.ReimbursementAmount)
	assert.Equal(t, 5400.0, result.TotalAmount)
	assert.Equal(t, 2, result.AttendanceDays)
	assert.Equal(t, "processed", result.Status)

	mockPayslipRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessEmployeePayroll - Payslip Already Exists
func TestProcessEmployeePayroll_PayslipAlreadyExists(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()

	// Mock expectations
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(true, nil)

	// Execute
	result, err := usecase.ProcessEmployeePayroll(1, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "payslip already exists for this period")

	mockPayslipRepo.AssertExpectations(t)
}

// Test ProcessEmployeePayroll - Check Payslip Exists Error
func TestProcessEmployeePayroll_CheckPayslipExistsError(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()

	// Mock expectations
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, errors.New("database error"))

	// Execute
	result, err := usecase.ProcessEmployeePayroll(1, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to check existing payslip")

	mockPayslipRepo.AssertExpectations(t)
}

// Test ProcessEmployeePayroll - Get Attendance Error
func TestProcessEmployeePayroll_GetAttendanceError(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()

	// Mock expectations
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return([]model.Attendance{}, errors.New("attendance error"))

	// Execute
	result, err := usecase.ProcessEmployeePayroll(1, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get attendance records")

	mockPayslipRepo.AssertExpectations(t)
}

// Test ProcessEmployeePayroll - Get Overtime Error
func TestProcessEmployeePayroll_GetOvertimeError(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	attendances := createTestAttendances()

	// Mock expectations
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return([]model.Overtime{}, errors.New("overtime error"))

	// Execute
	result, err := usecase.ProcessEmployeePayroll(1, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get overtime records")

	mockPayslipRepo.AssertExpectations(t)
}

// Test ProcessEmployeePayroll - Get Reimbursements Error
func TestProcessEmployeePayroll_GetReimbursementsError(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()

	// Mock expectations
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return([]model.Reimbursement{}, errors.New("reimbursement error"))

	// Execute
	result, err := usecase.ProcessEmployeePayroll(1, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get reimbursement records")

	mockPayslipRepo.AssertExpectations(t)
}

// Test ProcessEmployeePayroll - Create Payslip Error
func TestProcessEmployeePayroll_CreatePayslipError(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()

	// Mock expectations
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)
	mockPayslipRepo.On("CreatePayslip", mock.AnythingOfType("*model.Payslip")).Return(&model.Payslip{}, errors.New("create payslip error"))

	// Execute
	result, err := usecase.ProcessEmployeePayroll(1, req)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result) // CreatePayslip returns empty payslip and error
	assert.Equal(t, "create payslip error", err.Error())

	mockPayslipRepo.AssertExpectations(t)
}

// Test ProcessEmployeePayroll - Zero Values
func TestProcessEmployeePayroll_ZeroValues(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	req.BasicSalary = 0
	req.OvertimeRate = 0

	// Mock expectations with empty arrays
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return([]model.Attendance{}, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return([]model.Overtime{}, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return([]model.Reimbursement{}, nil)

	expectedPayslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      req.PayPeriodStart,
		PayPeriodEnd:        req.PayPeriodEnd,
		BasicSalary:         0,
		OvertimeHours:       0,
		OvertimeAmount:      0,
		ReimbursementAmount: 0,
		TotalAmount:         0,
		Status:              "processed",
		AttendanceDays:      0,
	}

	mockPayslipRepo.On("CreatePayslip", mock.MatchedBy(func(p *model.Payslip) bool {
		return p.EmployeeID == 1 &&
			p.BasicSalary == 0 &&
			p.OvertimeHours == 0 &&
			p.OvertimeAmount == 0 &&
			p.ReimbursementAmount == 0 &&
			p.TotalAmount == 0 &&
			p.AttendanceDays == 0
	})).Return(expectedPayslip, nil)

	// Execute
	result, err := usecase.ProcessEmployeePayroll(1, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0.0, result.BasicSalary)
	assert.Equal(t, 0, result.OvertimeHours)
	assert.Equal(t, 0.0, result.OvertimeAmount)
	assert.Equal(t, 0.0, result.ReimbursementAmount)
	assert.Equal(t, 0.0, result.TotalAmount)
	assert.Equal(t, 0, result.AttendanceDays)

	mockPayslipRepo.AssertExpectations(t)
}

// Test ProcessEmployeePayrollWithAudit - Success
func TestProcessEmployeePayrollWithAudit_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()
	auditDB := &middleware.AuditableDB{}

	// Mock expectations
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)

	expectedPayslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      req.PayPeriodStart,
		PayPeriodEnd:        req.PayPeriodEnd,
		BasicSalary:         5000.0,
		OvertimeHours:       5,
		OvertimeAmount:      250.0,
		ReimbursementAmount: 150.0,
		TotalAmount:         5400.0,
		Status:              "processed",
		AttendanceDays:      2,
	}

	mockPayslipRepo.On("CreatePayslipWithAudit", mock.AnythingOfType("*model.Payslip"), auditDB).Return(expectedPayslip, nil)

	// Execute
	result, err := usecase.ProcessEmployeePayrollWithAudit(1, req, auditDB)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.EmployeeID)
	assert.Equal(t, 5000.0, result.BasicSalary)
	assert.Equal(t, 5, result.OvertimeHours)
	assert.Equal(t, 250.0, result.OvertimeAmount)
	assert.Equal(t, 150.0, result.ReimbursementAmount)
	assert.Equal(t, 5400.0, result.TotalAmount)

	mockPayslipRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayroll - Success
func TestProcessAllEmployeesPayroll_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	employees := createTestEmployees()
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()

	// Mock expectations for getting employees
	mockEmployeeRepo.On("GetAllActiveEmployees").Return(employees, nil)

	// Mock expectations for each employee
	for _, emp := range employees {
		mockPayslipRepo.On("CheckPayslipExists", emp.ID, req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
		mockPayslipRepo.On("GetAttendanceForPeriod", emp.ID, req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
		mockPayslipRepo.On("GetOvertimeForPeriod", emp.ID, "2025-06-01", "2025-06-30").Return(overtimes, nil)
		mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", emp.ID, req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)

		expectedPayslip := &model.Payslip{
			DefaultAttribute:    model.DefaultAttribute{ID: emp.ID},
			EmployeeID:          emp.ID,
			PayPeriodStart:      req.PayPeriodStart,
			PayPeriodEnd:        req.PayPeriodEnd,
			BasicSalary:         5000.0,
			OvertimeHours:       5,
			OvertimeAmount:      250.0,
			ReimbursementAmount: 150.0,
			TotalAmount:         5400.0,
			Status:              "processed",
			AttendanceDays:      2,
		}

		mockPayslipRepo.On("CreatePayslip", mock.MatchedBy(func(p *model.Payslip) bool {
			return p.EmployeeID == emp.ID &&
				p.BasicSalary == 5000.0 &&
				p.OvertimeHours == 5 &&
				p.OvertimeAmount == 250.0 &&
				p.ReimbursementAmount == 150.0 &&
				p.TotalAmount == 5400.0 &&
				p.AttendanceDays == 2 &&
				p.Status == "processed"
		})).Return(expectedPayslip, nil)
	}

	// Execute
	payslips, errors := usecase.ProcessAllEmployeesPayroll(req)

	// Assert
	assert.Len(t, payslips, 2)
	assert.Len(t, errors, 0)
	assert.Equal(t, uint(1), payslips[0].EmployeeID)
	assert.Equal(t, uint(2), payslips[1].EmployeeID)

	mockPayslipRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayroll - Get Employees Error
func TestProcessAllEmployeesPayroll_GetEmployeesError(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()

	// Mock expectations
	mockEmployeeRepo.On("GetAllActiveEmployees").Return([]model.Employee{}, errors.New("employee fetch error"))

	// Execute
	payslips, errors := usecase.ProcessAllEmployeesPayroll(req)

	// Assert
	assert.Nil(t, payslips)
	assert.Len(t, errors, 1)
	assert.Contains(t, errors[0], "Failed to get employees")

	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayroll - Partial Success
func TestProcessAllEmployeesPayroll_PartialSuccess(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	employees := createTestEmployees()
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()

	// Mock expectations for getting employees
	mockEmployeeRepo.On("GetAllActiveEmployees").Return(employees, nil)

	// First employee succeeds
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)

	expectedPayslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      req.PayPeriodStart,
		PayPeriodEnd:        req.PayPeriodEnd,
		BasicSalary:         5000.0,
		OvertimeHours:       5,
		OvertimeAmount:      250.0,
		ReimbursementAmount: 150.0,
		TotalAmount:         5400.0,
		Status:              "processed",
		AttendanceDays:      2,
	}

	mockPayslipRepo.On("CreatePayslip", mock.MatchedBy(func(p *model.Payslip) bool {
		return p.EmployeeID == 1 &&
			p.BasicSalary == 5000.0 &&
			p.OvertimeHours == 5 &&
			p.OvertimeAmount == 250.0 &&
			p.ReimbursementAmount == 150.0 &&
			p.TotalAmount == 5400.0 &&
			p.AttendanceDays == 2 &&
			p.Status == "processed"
	})).Return(expectedPayslip, nil)

	// Second employee fails
	mockPayslipRepo.On("CheckPayslipExists", uint(2), req.PayPeriodStart, req.PayPeriodEnd).Return(false, errors.New("database error"))

	// Execute
	payslips, errors := usecase.ProcessAllEmployeesPayroll(req)

	// Assert
	assert.Len(t, payslips, 1)
	assert.Len(t, errors, 1)
	assert.Equal(t, uint(1), payslips[0].EmployeeID)
	assert.Contains(t, errors[0], "Employee 2:")
	assert.Contains(t, errors[0], "database error")

	mockPayslipRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayrollWithAudit - Success
func TestProcessAllEmployeesPayrollWithAudit_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	employees := createTestEmployees()
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()
	auditDB := &middleware.AuditableDB{}

	// Mock expectations for getting employees
	mockEmployeeRepo.On("GetAllActiveEmployees").Return(employees, nil)

	// Mock expectations for each employee
	for _, emp := range employees {
		mockPayslipRepo.On("CheckPayslipExists", emp.ID, req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
		mockPayslipRepo.On("GetAttendanceForPeriod", emp.ID, req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
		mockPayslipRepo.On("GetOvertimeForPeriod", emp.ID, "2025-06-01", "2025-06-30").Return(overtimes, nil)
		mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", emp.ID, req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)

		expectedPayslip := &model.Payslip{
			DefaultAttribute:    model.DefaultAttribute{ID: emp.ID},
			EmployeeID:          emp.ID,
			PayPeriodStart:      req.PayPeriodStart,
			PayPeriodEnd:        req.PayPeriodEnd,
			BasicSalary:         5000.0,
			OvertimeHours:       5,
			OvertimeAmount:      250.0,
			ReimbursementAmount: 150.0,
			TotalAmount:         5400.0,
			Status:              "processed",
			AttendanceDays:      2,
		}

		mockPayslipRepo.On("CreatePayslipWithAudit", mock.AnythingOfType("*model.Payslip"), auditDB).Return(expectedPayslip, nil)
	}

	// Execute
	payslips, errors := usecase.ProcessAllEmployeesPayrollWithAudit(req, auditDB)

	// Assert
	assert.Len(t, payslips, 2)
	assert.Len(t, errors, 0)
	assert.Equal(t, uint(1), payslips[0].EmployeeID)
	assert.Equal(t, uint(2), payslips[1].EmployeeID)

	mockPayslipRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayrollWithAudit - Get Employees Error
func TestProcessAllEmployeesPayrollWithAudit_GetEmployeesError(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	auditDB := &middleware.AuditableDB{}

	// Mock expectations - employee fetch fails
	mockEmployeeRepo.On("GetAllActiveEmployees").Return([]model.Employee{}, errors.New("database connection error"))

	// Execute
	payslips, errors := usecase.ProcessAllEmployeesPayrollWithAudit(req, auditDB)

	// Assert
	assert.Nil(t, payslips)
	assert.Len(t, errors, 1)
	assert.Contains(t, errors[0], "Failed to get employees")
	assert.Contains(t, errors[0], "database connection error")

	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayrollWithAudit - No Active Employees
func TestProcessAllEmployeesPayrollWithAudit_NoActiveEmployees(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	auditDB := &middleware.AuditableDB{}

	// Mock expectations - return empty employee list
	mockEmployeeRepo.On("GetAllActiveEmployees").Return([]model.Employee{}, nil)

	// Execute
	payslips, errors := usecase.ProcessAllEmployeesPayrollWithAudit(req, auditDB)

	// Assert
	assert.Len(t, payslips, 0)
	assert.Len(t, errors, 0)

	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayrollWithAudit - Partial Success with Audit Failures
func TestProcessAllEmployeesPayrollWithAudit_PartialSuccess(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	employees := createTestEmployees()
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()
	auditDB := &middleware.AuditableDB{}

	// Mock expectations for getting employees
	mockEmployeeRepo.On("GetAllActiveEmployees").Return(employees, nil)

	// First employee succeeds
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)

	expectedPayslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      req.PayPeriodStart,
		PayPeriodEnd:        req.PayPeriodEnd,
		BasicSalary:         5000.0,
		OvertimeHours:       5,
		OvertimeAmount:      250.0,
		ReimbursementAmount: 150.0,
		TotalAmount:         5400.0,
		Status:              "processed",
		AttendanceDays:      2,
	}

	mockPayslipRepo.On("CreatePayslip", mock.MatchedBy(func(p *model.Payslip) bool {
		return p.EmployeeID == 1 &&
			p.BasicSalary == 5000.0 &&
			p.OvertimeHours == 5 &&
			p.OvertimeAmount == 250.0 &&
			p.ReimbursementAmount == 150.0 &&
			p.TotalAmount == 5400.0 &&
			p.AttendanceDays == 2 &&
			p.Status == "processed"
	})).Return(expectedPayslip, nil)

	// Second employee fails at audit payslip creation
	mockPayslipRepo.On("CheckPayslipExists", uint(2), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(2), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(2), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(2), req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)
	mockPayslipRepo.On("CreatePayslipWithAudit", mock.MatchedBy(func(p *model.Payslip) bool {
		return p.EmployeeID == 2
	}), auditDB).Return(&model.Payslip{}, errors.New("audit trail creation failed"))

	// Execute
	payslips, errors := usecase.ProcessAllEmployeesPayrollWithAudit(req, auditDB)

	// Assert
	assert.Len(t, payslips, 1)
	assert.Len(t, errors, 1)
	assert.Equal(t, uint(1), payslips[0].EmployeeID)
	assert.Contains(t, errors[0], "Employee 2:")
	assert.Contains(t, errors[0], "audit trail creation failed")

	mockPayslipRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayrollWithAudit - All Employees Fail
func TestProcessAllEmployeesPayrollWithAudit_AllEmployeesFail(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	employees := createTestEmployees()
	auditDB := &middleware.AuditableDB{}

	// Mock expectations for getting employees
	mockEmployeeRepo.On("GetAllActiveEmployees").Return(employees, nil)

	// Both employees fail at different stages
	// First employee fails because payslip already exists
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(true, nil)

	// Second employee fails at attendance check
	mockPayslipRepo.On("CheckPayslipExists", uint(2), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(2), req.PayPeriodStart, req.PayPeriodEnd).Return([]model.Attendance{}, errors.New("audit database connection failed"))

	// Execute
	payslips, errors := usecase.ProcessAllEmployeesPayrollWithAudit(req, auditDB)

	// Assert
	assert.Len(t, payslips, 0) // No successful payslips
	assert.Len(t, errors, 2)   // Both employees have errors
	assert.Contains(t, errors[0], "Employee 1:")
	assert.Contains(t, errors[0], "payslip already exists for this period")
	assert.Contains(t, errors[1], "Employee 2:")
	assert.Contains(t, errors[1], "failed to get attendance records")

	mockPayslipRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayrollWithAudit - Audit-Specific Failures
func TestProcessAllEmployeesPayrollWithAudit_AuditSpecificFailures(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()

	// Create employees with different audit failure scenarios
	employees := []model.Employee{
		{DefaultAttribute: model.DefaultAttribute{ID: 1}, Name: "John Doe", Role: "employee", Active: true},
		{DefaultAttribute: model.DefaultAttribute{ID: 2}, Name: "Jane Smith", Role: "employee", Active: true},
		{DefaultAttribute: model.DefaultAttribute{ID: 3}, Name: "Bob Wilson", Role: "employee", Active: true},
	}

	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()
	auditDB := &middleware.AuditableDB{}

	// Mock expectations for getting employees
	mockEmployeeRepo.On("GetAllActiveEmployees").Return(employees, nil)

	// Employee 1: Success with audit
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)
	expectedPayslip1 := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{ID: 1},
		EmployeeID:       1, BasicSalary: 5000.0, OvertimeHours: 5, OvertimeAmount: 250.0,
		ReimbursementAmount: 150.0, TotalAmount: 5400.0, Status: "processed", AttendanceDays: 2,
	}
	mockPayslipRepo.On("CreatePayslipWithAudit", mock.MatchedBy(func(p *model.Payslip) bool {
		return p.EmployeeID == 1
	}), auditDB).Return(expectedPayslip1, nil)

	// Employee 2: Fails at audit database constraint violation
	mockPayslipRepo.On("CheckPayslipExists", uint(2), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(2), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(2), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(2), req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)
	mockPayslipRepo.On("CreatePayslipWithAudit", mock.MatchedBy(func(p *model.Payslip) bool {
		return p.EmployeeID == 2
	}), auditDB).Return(&model.Payslip{}, errors.New("audit constraint violation: duplicate entry"))

	// Employee 3: Fails at audit transaction rollback
	mockPayslipRepo.On("CheckPayslipExists", uint(3), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(3), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(3), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(3), req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)
	mockPayslipRepo.On("CreatePayslipWithAudit", mock.MatchedBy(func(p *model.Payslip) bool {
		return p.EmployeeID == 3
	}), auditDB).Return(&model.Payslip{}, errors.New("audit transaction failed: unable to commit"))

	// Execute
	payslips, errors := usecase.ProcessAllEmployeesPayrollWithAudit(req, auditDB)

	// Assert
	assert.Len(t, payslips, 1) // Only employee 1 successful
	assert.Len(t, errors, 2)   // Employees 2, 3 failed
	assert.Equal(t, uint(1), payslips[0].EmployeeID)

	// Check specific audit error messages
	assert.Contains(t, errors[0], "Employee 2:")
	assert.Contains(t, errors[0], "audit constraint violation")
	assert.Contains(t, errors[1], "Employee 3:")
	assert.Contains(t, errors[1], "audit transaction failed")

	mockPayslipRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayrollWithAudit - Large Employee Set with Audit
func TestProcessAllEmployeesPayrollWithAudit_LargeEmployeeSet(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()
	auditDB := &middleware.AuditableDB{}

	// Create 25 employees (smaller than the regular test to focus on audit functionality)
	employees := make([]model.Employee, 25)
	for i := 0; i < 25; i++ {
		employees[i] = model.Employee{
			DefaultAttribute: model.DefaultAttribute{ID: uint(i + 1)},
			Name:             fmt.Sprintf("Employee %d", i+1),
			Role:             "employee",
			Active:           true,
		}
	}

	// Mock expectations for getting employees
	mockEmployeeRepo.On("GetAllActiveEmployees").Return(employees, nil)

	// Mock expectations for each employee (all successful with audit)
	for i := 0; i < 25; i++ {
		empID := uint(i + 1)
		mockPayslipRepo.On("CheckPayslipExists", empID, req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
		mockPayslipRepo.On("GetAttendanceForPeriod", empID, req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
		mockPayslipRepo.On("GetOvertimeForPeriod", empID, "2025-06-01", "2025-06-30").Return(overtimes, nil)
		mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", empID, req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)

		expectedPayslip := &model.Payslip{
			DefaultAttribute: model.DefaultAttribute{ID: empID},
			EmployeeID:       empID, BasicSalary: 5000.0, OvertimeHours: 5, OvertimeAmount: 250.0,
			ReimbursementAmount: 150.0, TotalAmount: 5400.0, Status: "processed", AttendanceDays: 2,
		}

		mockPayslipRepo.On("CreatePayslipWithAudit", mock.AnythingOfType("*model.Payslip"), auditDB).Return(expectedPayslip, nil)
	}

	// Execute
	payslips, errors := usecase.ProcessAllEmployeesPayrollWithAudit(req, auditDB)

	// Assert
	assert.Len(t, payslips, 25)
	assert.Len(t, errors, 0)

	// Verify all employee IDs are present and properly audited
	for i := 0; i < 25; i++ {
		assert.Equal(t, uint(i+1), payslips[i].EmployeeID)
		assert.Equal(t, 5400.0, payslips[i].TotalAmount)
	}

	mockPayslipRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

// Test ProcessAllEmployeesPayrollWithAudit - Nil Audit DB
func TestProcessAllEmployeesPayrollWithAudit_NilAuditDB(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	employees := []model.Employee{
		{DefaultAttribute: model.DefaultAttribute{ID: 1}, Name: "John Doe", Role: "employee", Active: true},
	}
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()

	// Mock expectations for getting employees
	mockEmployeeRepo.On("GetAllActiveEmployees").Return(employees, nil)

	// Mock expectations for employee
	mockPayslipRepo.On("CheckPayslipExists", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(false, nil)
	mockPayslipRepo.On("GetAttendanceForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(attendances, nil)
	mockPayslipRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", uint(1), req.PayPeriodStart, req.PayPeriodEnd).Return(reimbursements, nil)

	expectedPayslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      req.PayPeriodStart,
		PayPeriodEnd:        req.PayPeriodEnd,
		BasicSalary:         5000.0,
		OvertimeHours:       5,
		OvertimeAmount:      250.0,
		ReimbursementAmount: 150.0,
		TotalAmount:         5400.0,
		Status:              "processed",
		AttendanceDays:      2,
	}

	// Pass nil auditDB - this should still work (function should handle nil gracefully)
	mockPayslipRepo.On("CreatePayslipWithAudit", mock.MatchedBy(func(p *model.Payslip) bool {
		return p.EmployeeID == 1
	}), (*middleware.AuditableDB)(nil)).Return(expectedPayslip, nil)

	// Execute with nil auditDB
	payslips, errors := usecase.ProcessAllEmployeesPayrollWithAudit(req, nil)

	// Assert
	assert.Len(t, payslips, 1)
	assert.Len(t, errors, 0)
	assert.Equal(t, uint(1), payslips[0].EmployeeID)
	assert.Equal(t, 5400.0, payslips[0].TotalAmount)

	mockPayslipRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

// Benchmark test for ProcessAllEmployeesPayrollWithAudit
func BenchmarkProcessAllEmployeesPayrollWithAudit(b *testing.B) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	req := createTestPayrollRequest()
	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()
	auditDB := &middleware.AuditableDB{}

	// Create 5 employees for benchmark (smaller set for audit operations)
	employees := make([]model.Employee, 5)
	for i := 0; i < 5; i++ {
		employees[i] = model.Employee{
			DefaultAttribute: model.DefaultAttribute{ID: uint(i + 1)},
			Name:             fmt.Sprintf("Employee %d", i+1),
			Role:             "employee",
			Active:           true,
		}
	}

	// Setup mocks for all iterations
	mockEmployeeRepo.On("GetAllActiveEmployees").Return(employees, nil)

	for i := 0; i < 5; i++ {
		empID := uint(i + 1)
		mockPayslipRepo.On("CheckPayslipExists", empID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil)
		mockPayslipRepo.On("GetAttendanceForPeriod", empID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(attendances, nil)
		mockPayslipRepo.On("GetOvertimeForPeriod", empID, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(overtimes, nil)
		mockPayslipRepo.On("GetApprovedReimbursementsForPeriod", empID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(reimbursements, nil)

		expectedPayslip := &model.Payslip{DefaultAttribute: model.DefaultAttribute{ID: empID}, EmployeeID: empID}
		mockPayslipRepo.On("CreatePayslipWithAudit", mock.AnythingOfType("*model.Payslip"), mock.AnythingOfType("*middleware.AuditableDB")).Return(expectedPayslip, nil)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = usecase.ProcessAllEmployeesPayrollWithAudit(req, auditDB)
	}
}

// Test BuildDetailedPayslipResponse - Success with All Data
func TestBuildDetailedPayslipResponse_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	// Create test data
	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: 1},
		Name:             "John Doe",
		Role:             "employee",
		Active:           true,
	}

	payslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:        time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		BasicSalary:         5000.0,
		OvertimeHours:       8,
		OvertimeAmount:      400.0,
		ReimbursementAmount: 200.0,
		TotalAmount:         5600.0,
		ProcessedAt:         time.Date(2025, 6, 30, 15, 30, 0, 0, time.UTC),
		Status:              "processed",
		AttendanceDays:      22,
	}

	attendances := createTestAttendances()
	overtimes := createTestOvertimes()
	reimbursements := createTestReimbursements()

	// Execute
	result := usecase.BuildDetailedPayslipResponse(payslip, employee, attendances, overtimes, reimbursements)

	// Assert main fields
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result["payslip_id"])
	assert.Equal(t, uint(1), result["employee_id"])
	assert.Equal(t, "John Doe", result["employee_name"])
	assert.Equal(t, payslip.PayPeriodStart, result["pay_period_start"])
	assert.Equal(t, payslip.PayPeriodEnd, result["pay_period_end"])
	assert.Equal(t, payslip.ProcessedAt, result["processed_at"])
	assert.Equal(t, "processed", result["status"])

	// Assert summary exists and has correct structure
	summary := result["summary"].(map[string]interface{})
	assert.NotNil(t, summary)
	assert.Equal(t, 5000.0, summary["basic_salary"])
	assert.Equal(t, 22, summary["total_attendance_days"])
	assert.Equal(t, 8, summary["total_overtime_hours"])
	assert.Equal(t, 400.0, summary["overtime_amount"])
	assert.Equal(t, 200.0, summary["reimbursement_amount"])
	assert.Equal(t, 5600.0, summary["total_take_home_pay"])

	// Assert breakdown arrays exist and have correct structure
	attendanceBreakdown := result["attendance_breakdown"].([]map[string]interface{})
	assert.NotNil(t, attendanceBreakdown)
	assert.Len(t, attendanceBreakdown, 2)

	// Check first attendance entry
	firstAttendance := attendanceBreakdown[0]
	assert.Equal(t, time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC), firstAttendance["date"])
	assert.Equal(t, time.Date(2025, 6, 1, 9, 0, 0, 0, time.UTC), firstAttendance["check_in"])
	assert.Equal(t, 8, firstAttendance["hours_worked"])
	assert.Equal(t, "present", firstAttendance["status"])

	// Check overtime breakdown
	overtimeBreakdown := result["overtime_breakdown"].([]map[string]interface{})
	assert.NotNil(t, overtimeBreakdown)
	assert.Len(t, overtimeBreakdown, 2)

	// Check first overtime entry (rate should be calculated: 400.0 / 8 = 50.0)
	firstOvertime := overtimeBreakdown[0]
	assert.Equal(t, "2025-06-01", firstOvertime["date"])
	assert.Equal(t, 2, firstOvertime["hours"])
	assert.Equal(t, 50.0, firstOvertime["rate"])
	assert.Equal(t, 100.0, firstOvertime["amount"]) // 2 * 50.0
	assert.Equal(t, "Project deadline", firstOvertime["reason"])

	// Check reimbursement breakdown
	reimbursementBreakdown := result["reimbursement_breakdown"].([]map[string]interface{})
	assert.NotNil(t, reimbursementBreakdown)
	assert.Len(t, reimbursementBreakdown, 2)

	// Check first reimbursement entry
	firstReimbursement := reimbursementBreakdown[0]
	assert.Equal(t, time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC), firstReimbursement["date"])
	assert.Equal(t, 100.0, firstReimbursement["amount"])
	assert.Equal(t, "Travel expenses", firstReimbursement["reason"])
	assert.Equal(t, model.ReimbursementApproved, firstReimbursement["status"])
}

// Test BuildDetailedPayslipResponse - Empty Data Arrays
func TestBuildDetailedPayslipResponse_EmptyArrays(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	// Create minimal test data
	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: 1},
		Name:             "Jane Smith",
		Role:             "employee",
		Active:           true,
	}

	payslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:        time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		BasicSalary:         3000.0,
		OvertimeHours:       0,
		OvertimeAmount:      0,
		ReimbursementAmount: 0,
		TotalAmount:         3000.0,
		ProcessedAt:         time.Date(2025, 6, 30, 10, 0, 0, 0, time.UTC),
		Status:              "processed",
		AttendanceDays:      0,
	}

	// Execute with empty arrays
	result := usecase.BuildDetailedPayslipResponse(payslip, employee, []model.Attendance{}, []model.Overtime{}, []model.Reimbursement{})

	// Assert main fields
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result["payslip_id"])
	assert.Equal(t, uint(1), result["employee_id"])
	assert.Equal(t, "Jane Smith", result["employee_name"])
	assert.Equal(t, "processed", result["status"])

	// Assert summary with zero values
	summary := result["summary"].(map[string]interface{})
	assert.Equal(t, 3000.0, summary["basic_salary"])
	assert.Equal(t, 0, summary["total_attendance_days"])
	assert.Equal(t, 0, summary["total_overtime_hours"])
	assert.Equal(t, 0.0, summary["overtime_amount"])
	assert.Equal(t, 0.0, summary["reimbursement_amount"])
	assert.Equal(t, 3000.0, summary["total_take_home_pay"])

	// Assert breakdown arrays are empty but not nil
	attendanceBreakdown := result["attendance_breakdown"].([]map[string]interface{})
	assert.Len(t, attendanceBreakdown, 0)

	overtimeBreakdown := result["overtime_breakdown"].([]map[string]interface{})
	assert.Len(t, overtimeBreakdown, 0)

	reimbursementBreakdown := result["reimbursement_breakdown"].([]map[string]interface{})
	assert.Len(t, reimbursementBreakdown, 0)
}

// Test BuildDetailedPayslipResponse - Zero Overtime Hours Edge Case
func TestBuildDetailedPayslipResponse_ZeroOvertimeHours(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: 1},
		Name:             "Bob Wilson",
		Role:             "employee",
		Active:           true,
	}

	payslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:        time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		BasicSalary:         4000.0,
		OvertimeHours:       0, // Zero overtime hours
		OvertimeAmount:      0.0,
		ReimbursementAmount: 50.0,
		TotalAmount:         4050.0,
		ProcessedAt:         time.Date(2025, 6, 30, 12, 0, 0, 0, time.UTC),
		Status:              "processed",
		AttendanceDays:      20,
	}

	// Create overtime records with zero hours (edge case)
	overtimes := []model.Overtime{
		{
			EmployeeID:   1,
			OvertimeDate: "2025-06-01",
			Hours:        0, // Zero hours
			Reason:       "Cancelled overtime",
			Status:       "approved",
		},
	}

	reimbursements := []model.Reimbursement{
		{
			EmployeeID:        1,
			ReimbursementDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			Amount:            50.0,
			Reason:            "Parking fee",
			Status:            model.ReimbursementApproved,
		},
	}

	// Execute
	result := usecase.BuildDetailedPayslipResponse(payslip, employee, []model.Attendance{}, overtimes, reimbursements)

	// Assert
	summary := result["summary"].(map[string]interface{})
	assert.Equal(t, 0, summary["total_overtime_hours"])
	assert.Equal(t, 0.0, summary["overtime_amount"])

	// Check overtime breakdown - rate should be 0.0 when overtime hours is 0
	overtimeBreakdown := result["overtime_breakdown"].([]map[string]interface{})
	assert.Len(t, overtimeBreakdown, 1)
	firstOvertime := overtimeBreakdown[0]
	assert.Equal(t, 0.0, firstOvertime["rate"])   // Should be 0.0 when no overtime hours
	assert.Equal(t, 0.0, firstOvertime["amount"]) // 0 * 0.0 = 0.0
}

// Test BuildDetailedPayslipResponse - Large Data Sets
func TestBuildDetailedPayslipResponse_LargeDataSets(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: 1},
		Name:             "Alice Johnson",
		Role:             "senior_employee",
		Active:           true,
	}

	payslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:        time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		BasicSalary:         6000.0,
		OvertimeHours:       30,
		OvertimeAmount:      1500.0, // 30 * 50.0
		ReimbursementAmount: 500.0,
		TotalAmount:         8000.0,
		ProcessedAt:         time.Date(2025, 6, 30, 18, 0, 0, 0, time.UTC),
		Status:              "processed",
		AttendanceDays:      30,
	}

	// Create large data sets (30 days of attendance, 10 overtime entries, 5 reimbursements)
	attendances := make([]model.Attendance, 30)
	for i := 0; i < 30; i++ {
		checkout := time.Date(2025, 6, i+1, 17, 0, 0, 0, time.UTC)
		attendances[i] = model.Attendance{
			EmployeeID:  1,
			Date:        time.Date(2025, 6, i+1, 0, 0, 0, 0, time.UTC),
			Checkin:     time.Date(2025, 6, i+1, 9, 0, 0, 0, time.UTC),
			Checkout:    &checkout,
			HoursWorked: 8,
			Status:      "present",
		}
	}

	overtimes := make([]model.Overtime, 10)
	for i := 0; i < 10; i++ {
		overtimes[i] = model.Overtime{
			EmployeeID:   1,
			OvertimeDate: fmt.Sprintf("2025-06-%02d", i+1),
			Hours:        3,
			Reason:       fmt.Sprintf("Project work %d", i+1),
			Status:       "approved",
		}
	}

	reimbursements := make([]model.Reimbursement, 5)
	for i := 0; i < 5; i++ {
		reimbursements[i] = model.Reimbursement{
			EmployeeID:        1,
			ReimbursementDate: time.Date(2025, 6, (i+1)*5, 0, 0, 0, 0, time.UTC),
			Amount:            100.0,
			Reason:            fmt.Sprintf("Business expense %d", i+1),
			Status:            model.ReimbursementApproved,
		}
	}

	// Execute
	result := usecase.BuildDetailedPayslipResponse(payslip, employee, attendances, overtimes, reimbursements)

	// Assert arrays have correct lengths
	attendanceBreakdown := result["attendance_breakdown"].([]map[string]interface{})
	assert.Len(t, attendanceBreakdown, 30)

	overtimeBreakdown := result["overtime_breakdown"].([]map[string]interface{})
	assert.Len(t, overtimeBreakdown, 10)

	reimbursementBreakdown := result["reimbursement_breakdown"].([]map[string]interface{})
	assert.Len(t, reimbursementBreakdown, 5)

	// Assert summary totals are correct
	summary := result["summary"].(map[string]interface{})
	assert.Equal(t, 6000.0, summary["basic_salary"])
	assert.Equal(t, 30, summary["total_attendance_days"])
	assert.Equal(t, 30, summary["total_overtime_hours"])
	assert.Equal(t, 1500.0, summary["overtime_amount"])
	assert.Equal(t, 500.0, summary["reimbursement_amount"])
	assert.Equal(t, 8000.0, summary["total_take_home_pay"])

	// Verify overtime rate calculation (1500.0 / 30 = 50.0)
	firstOvertime := overtimeBreakdown[0]
	assert.Equal(t, 50.0, firstOvertime["rate"])
	assert.Equal(t, 150.0, firstOvertime["amount"]) // 3 * 50.0
}

// Test BuildDetailedPayslipResponse - Different Employee Roles
func TestBuildDetailedPayslipResponse_DifferentEmployeeRoles(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	// Test with different employee roles
	testCases := []struct {
		name     string
		employee *model.Employee
		payslip  *model.Payslip
	}{
		{
			name: "Admin Employee",
			employee: &model.Employee{
				DefaultAttribute: model.DefaultAttribute{ID: 1},
				Name:             "Admin User",
				Role:             "admin",
				Active:           true,
			},
			payslip: &model.Payslip{
				DefaultAttribute: model.DefaultAttribute{ID: 1},
				EmployeeID:       1, BasicSalary: 8000.0, OvertimeHours: 0, OvertimeAmount: 0.0,
				ReimbursementAmount: 0.0, TotalAmount: 8000.0, Status: "processed", AttendanceDays: 25,
			},
		},
		{
			name: "Manager Employee",
			employee: &model.Employee{
				DefaultAttribute: model.DefaultAttribute{ID: 2},
				Name:             "Manager User",
				Role:             "manager",
				Active:           true,
			},
			payslip: &model.Payslip{
				DefaultAttribute: model.DefaultAttribute{ID: 2},
				EmployeeID:       2, BasicSalary: 7000.0, OvertimeHours: 5, OvertimeAmount: 375.0,
				ReimbursementAmount: 200.0, TotalAmount: 7575.0, Status: "processed", AttendanceDays: 23,
			},
		},
		{
			name: "Regular Employee",
			employee: &model.Employee{
				DefaultAttribute: model.DefaultAttribute{ID: 3},
				Name:             "Regular User",
				Role:             "employee",
				Active:           true,
			},
			payslip: &model.Payslip{
				DefaultAttribute: model.DefaultAttribute{ID: 3},
				EmployeeID:       3, BasicSalary: 4000.0, OvertimeHours: 10, OvertimeAmount: 500.0,
				ReimbursementAmount: 100.0, TotalAmount: 4600.0, Status: "processed", AttendanceDays: 22,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.payslip.PayPeriodStart = time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
			tc.payslip.PayPeriodEnd = time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)
			tc.payslip.ProcessedAt = time.Date(2025, 6, 30, 15, 0, 0, 0, time.UTC)

			result := usecase.BuildDetailedPayslipResponse(tc.payslip, tc.employee, []model.Attendance{}, []model.Overtime{}, []model.Reimbursement{})

			// Assert employee-specific fields
			assert.Equal(t, tc.employee.ID, result["employee_id"])
			assert.Equal(t, tc.employee.Name, result["employee_name"])

			// Assert payslip-specific calculations
			summary := result["summary"].(map[string]interface{})
			assert.Equal(t, tc.payslip.BasicSalary, summary["basic_salary"])
			assert.Equal(t, tc.payslip.AttendanceDays, summary["total_attendance_days"])
			assert.Equal(t, tc.payslip.OvertimeHours, summary["total_overtime_hours"])
			assert.Equal(t, tc.payslip.OvertimeAmount, summary["overtime_amount"])
			assert.Equal(t, tc.payslip.ReimbursementAmount, summary["reimbursement_amount"])
			assert.Equal(t, tc.payslip.TotalAmount, summary["total_take_home_pay"])
		})
	}
}

// Test BuildDetailedPayslipResponse - Nil Pointers Edge Case
func TestBuildDetailedPayslipResponse_NilPointers(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: 1},
		Name:             "Test Employee",
		Role:             "employee",
		Active:           true,
	}

	payslip := &model.Payslip{
		DefaultAttribute:    model.DefaultAttribute{ID: 1},
		EmployeeID:          1,
		PayPeriodStart:      time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:        time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		BasicSalary:         5000.0,
		OvertimeHours:       0,
		OvertimeAmount:      0.0,
		ReimbursementAmount: 0.0,
		TotalAmount:         5000.0,
		ProcessedAt:         time.Date(2025, 6, 30, 12, 0, 0, 0, time.UTC),
		Status:              "processed",
		AttendanceDays:      0,
	}

	// Create attendance with nil Checkout pointer
	attendances := []model.Attendance{
		{
			EmployeeID:  1,
			Date:        time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			Checkin:     time.Date(2025, 6, 1, 9, 0, 0, 0, time.UTC),
			Checkout:    nil, // Nil pointer
			HoursWorked: 0,
			Status:      "absent",
		},
	}

	// Execute - should handle nil pointers gracefully
	result := usecase.BuildDetailedPayslipResponse(payslip, employee, attendances, []model.Overtime{}, []model.Reimbursement{})

	// Assert
	assert.NotNil(t, result)
	attendanceBreakdown := result["attendance_breakdown"].([]map[string]interface{})
	assert.Len(t, attendanceBreakdown, 1)

	firstAttendance := attendanceBreakdown[0]
	assert.Equal(t, time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC), firstAttendance["date"])
	assert.Equal(t, time.Date(2025, 6, 1, 9, 0, 0, 0, time.UTC), firstAttendance["check_in"])
	assert.Nil(t, firstAttendance["check_out"]) // Should handle nil gracefully
	assert.Equal(t, 0, firstAttendance["hours_worked"])
	assert.Equal(t, "absent", firstAttendance["status"])
}

// ===========================
// Helper Functions Tests
// ===========================

// Test calculateTotalOvertimeHours - Success with Multiple Overtimes
func TestCalculateTotalOvertimeHours_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	overtimes := []model.Overtime{
		{Hours: 2},
		{Hours: 3},
		{Hours: 1},
		{Hours: 4},
	}

	result := usecase.calculateTotalOvertimeHours(overtimes)
	assert.Equal(t, 10, result) // 2+3+1+4 = 10
}

// Test calculateTotalOvertimeHours - Empty Array
func TestCalculateTotalOvertimeHours_EmptyArray(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	overtimes := []model.Overtime{}

	result := usecase.calculateTotalOvertimeHours(overtimes)
	assert.Equal(t, 0, result)
}

// Test calculateTotalOvertimeHours - Zero Hours
func TestCalculateTotalOvertimeHours_ZeroHours(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	overtimes := []model.Overtime{
		{Hours: 0},
		{Hours: 0},
		{Hours: 5},
		{Hours: 0},
	}

	result := usecase.calculateTotalOvertimeHours(overtimes)
	assert.Equal(t, 5, result) // Only the 5 should count
}

// Test calculateTotalReimbursementAmount - Success with Multiple Reimbursements
func TestCalculateTotalReimbursementAmount_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	reimbursements := []model.Reimbursement{
		{Amount: 100.50},
		{Amount: 250.75},
		{Amount: 75.25},
		{Amount: 300.00},
	}

	result := usecase.calculateTotalReimbursementAmount(reimbursements)
	assert.Equal(t, 726.50, result) // 100.50+250.75+75.25+300.00 = 726.50
}

// Test calculateTotalReimbursementAmount - Empty Array
func TestCalculateTotalReimbursementAmount_EmptyArray(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	reimbursements := []model.Reimbursement{}

	result := usecase.calculateTotalReimbursementAmount(reimbursements)
	assert.Equal(t, 0.0, result)
}

// Test calculateTotalReimbursementAmount - Zero Amounts
func TestCalculateTotalReimbursementAmount_ZeroAmounts(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	reimbursements := []model.Reimbursement{
		{Amount: 0.0},
		{Amount: 100.0},
		{Amount: 0.0},
		{Amount: 50.0},
	}

	result := usecase.calculateTotalReimbursementAmount(reimbursements)
	assert.Equal(t, 150.0, result) // Only 100.0 + 50.0
}

// Test buildAttendanceBreakdown - Success with Multiple Attendances
func TestBuildAttendanceBreakdown_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	attendances := []model.Attendance{
		{
			Date:        time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			Checkin:     time.Date(2025, 6, 1, 9, 0, 0, 0, time.UTC),
			Checkout:    &[]time.Time{time.Date(2025, 6, 1, 17, 0, 0, 0, time.UTC)}[0],
			HoursWorked: 8,
			Status:      "present",
		},
		{
			Date:        time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC),
			Checkin:     time.Date(2025, 6, 2, 9, 30, 0, 0, time.UTC),
			Checkout:    &[]time.Time{time.Date(2025, 6, 2, 17, 30, 0, 0, time.UTC)}[0],
			HoursWorked: 8,
			Status:      "present",
		},
	}

	result := usecase.buildAttendanceBreakdown(attendances)

	assert.Len(t, result, 2)

	// Check first attendance
	first := result[0]
	assert.Equal(t, time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC), first["date"])
	assert.Equal(t, time.Date(2025, 6, 1, 9, 0, 0, 0, time.UTC), first["check_in"])
	assert.NotNil(t, first["check_out"])
	checkoutTime := first["check_out"].(*time.Time)
	assert.Equal(t, time.Date(2025, 6, 1, 17, 0, 0, 0, time.UTC), *checkoutTime)
	assert.Equal(t, 8, first["hours_worked"])
	assert.Equal(t, "present", first["status"])

	// Check second attendance
	second := result[1]
	assert.Equal(t, time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC), second["date"])
	assert.Equal(t, time.Date(2025, 6, 2, 9, 30, 0, 0, time.UTC), second["check_in"])
	assert.NotNil(t, second["check_out"])
	checkoutTime2 := second["check_out"].(*time.Time)
	assert.Equal(t, time.Date(2025, 6, 2, 17, 30, 0, 0, time.UTC), *checkoutTime2)
	assert.Equal(t, 8, second["hours_worked"])
	assert.Equal(t, "present", second["status"])
}

// Test buildAttendanceBreakdown - Empty Array
func TestBuildAttendanceBreakdown_EmptyArray(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	attendances := []model.Attendance{}

	result := usecase.buildAttendanceBreakdown(attendances)
	assert.Len(t, result, 0)
	// The result should be an empty slice, but Go allows either empty slice or nil
	if result != nil {
		assert.Equal(t, []map[string]interface{}{}, result)
	}
}

// Test buildOvertimeBreakdown - Success with Positive Overtime Hours
func TestBuildOvertimeBreakdown_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	payslip := &model.Payslip{
		OvertimeHours:  10,
		OvertimeAmount: 500.0, // Rate = 500/10 = 50.0
	}

	overtimes := []model.Overtime{
		{
			OvertimeDate: "2025-06-01",
			Hours:        2,
			Reason:       "Project deadline",
		},
		{
			OvertimeDate: "2025-06-02",
			Hours:        3,
			Reason:       "Urgent fix",
		},
	}

	result := usecase.buildOvertimeBreakdown(overtimes, payslip)

	assert.Len(t, result, 2)

	// Check first overtime
	first := result[0]
	assert.Equal(t, "2025-06-01", first["date"])
	assert.Equal(t, 2, first["hours"])
	assert.Equal(t, 50.0, first["rate"])
	assert.Equal(t, 100.0, first["amount"]) // 2 * 50.0
	assert.Equal(t, "Project deadline", first["reason"])

	// Check second overtime
	second := result[1]
	assert.Equal(t, "2025-06-02", second["date"])
	assert.Equal(t, 3, second["hours"])
	assert.Equal(t, 50.0, second["rate"])
	assert.Equal(t, 150.0, second["amount"]) // 3 * 50.0
	assert.Equal(t, "Urgent fix", second["reason"])
}

// Test buildOvertimeBreakdown - Zero Overtime Hours in Payslip
func TestBuildOvertimeBreakdown_ZeroOvertimeHours(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	payslip := &model.Payslip{
		OvertimeHours:  0,
		OvertimeAmount: 0.0, // Rate should be 0.0 to avoid division by zero
	}

	overtimes := []model.Overtime{
		{
			OvertimeDate: "2025-06-01",
			Hours:        2,
			Reason:       "Emergency work",
		},
	}

	result := usecase.buildOvertimeBreakdown(overtimes, payslip)

	assert.Len(t, result, 1)

	first := result[0]
	assert.Equal(t, "2025-06-01", first["date"])
	assert.Equal(t, 2, first["hours"])
	assert.Equal(t, 0.0, first["rate"])   // Should be 0.0 when no overtime hours in payslip
	assert.Equal(t, 0.0, first["amount"]) // 2 * 0.0
	assert.Equal(t, "Emergency work", first["reason"])
}

// Test buildOvertimeBreakdown - Empty Array
func TestBuildOvertimeBreakdown_EmptyArray(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	payslip := &model.Payslip{
		OvertimeHours:  10,
		OvertimeAmount: 500.0,
	}

	overtimes := []model.Overtime{}

	result := usecase.buildOvertimeBreakdown(overtimes, payslip)
	assert.Len(t, result, 0)
	assert.NotNil(t, result)
}

// Test buildReimbursementBreakdown - Success with Multiple Reimbursements
func TestBuildReimbursementBreakdown_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	reimbursements := []model.Reimbursement{
		{
			ReimbursementDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			Amount:            100.0,
			Reason:            "Travel expenses",
			Status:            model.ReimbursementApproved,
		},
		{
			ReimbursementDate: time.Date(2025, 6, 5, 0, 0, 0, 0, time.UTC),
			Amount:            50.0,
			Reason:            "Meal allowance",
			Status:            model.ReimbursementApproved,
		},
	}

	result := usecase.buildReimbursementBreakdown(reimbursements)

	assert.Len(t, result, 2)

	// Check first reimbursement
	first := result[0]
	assert.Equal(t, time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC), first["date"])
	assert.Equal(t, 100.0, first["amount"])
	assert.Equal(t, "Travel expenses", first["reason"])
	assert.Equal(t, model.ReimbursementApproved, first["status"])

	// Check second reimbursement
	second := result[1]
	assert.Equal(t, time.Date(2025, 6, 5, 0, 0, 0, 0, time.UTC), second["date"])
	assert.Equal(t, 50.0, second["amount"])
	assert.Equal(t, "Meal allowance", second["reason"])
	assert.Equal(t, model.ReimbursementApproved, second["status"])
}

// Test buildReimbursementBreakdown - Empty Array
func TestBuildReimbursementBreakdown_EmptyArray(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	reimbursements := []model.Reimbursement{}

	result := usecase.buildReimbursementBreakdown(reimbursements)
	assert.Len(t, result, 0)
	assert.NotNil(t, result)
}

// Test calculateEmployeeSummary - Success with Multiple Payslips
func TestCalculateEmployeeSummary_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	empPayslips := []model.Payslip{
		{
			TotalAmount:         5600.0,
			BasicSalary:         5000.0,
			OvertimeAmount:      400.0,
			ReimbursementAmount: 200.0,
			AttendanceDays:      22,
			OvertimeHours:       8,
		},
		{
			TotalAmount:         5350.0,
			BasicSalary:         5000.0,
			OvertimeAmount:      250.0,
			ReimbursementAmount: 100.0,
			AttendanceDays:      21,
			OvertimeHours:       5,
		},
		{
			TotalAmount:         5800.0,
			BasicSalary:         5000.0,
			OvertimeAmount:      600.0,
			ReimbursementAmount: 200.0,
			AttendanceDays:      23,
			OvertimeHours:       12,
		},
	}

	result := usecase.calculateEmployeeSummary(uint(1), empPayslips, "John Doe")

	assert.Equal(t, uint(1), result["employee_id"])
	assert.Equal(t, "John Doe", result["employee_name"])
	assert.Equal(t, 3, result["payslip_count"])
	assert.Equal(t, 16750.0, result["total_take_home_pay"])  // 5600+5350+5800
	assert.Equal(t, 15000.0, result["total_basic_salary"])   // 5000+5000+5000
	assert.Equal(t, 1250.0, result["total_overtime_amount"]) // 400+250+600
	assert.Equal(t, 500.0, result["total_reimbursement"])    // 200+100+200
	assert.Equal(t, 66, result["total_attendance_days"])     // 22+21+23
	assert.Equal(t, 25, result["total_overtime_hours"])      // 8+5+12
}

// Test calculateEmployeeSummary - Single Payslip
func TestCalculateEmployeeSummary_SinglePayslip(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	empPayslips := []model.Payslip{
		{
			TotalAmount:         5600.0,
			BasicSalary:         5000.0,
			OvertimeAmount:      400.0,
			ReimbursementAmount: 200.0,
			AttendanceDays:      22,
			OvertimeHours:       8,
		},
	}

	result := usecase.calculateEmployeeSummary(uint(2), empPayslips, "Jane Smith")

	assert.Equal(t, uint(2), result["employee_id"])
	assert.Equal(t, "Jane Smith", result["employee_name"])
	assert.Equal(t, 1, result["payslip_count"])
	assert.Equal(t, 5600.0, result["total_take_home_pay"])
	assert.Equal(t, 5000.0, result["total_basic_salary"])
	assert.Equal(t, 400.0, result["total_overtime_amount"])
	assert.Equal(t, 200.0, result["total_reimbursement"])
	assert.Equal(t, 22, result["total_attendance_days"])
	assert.Equal(t, 8, result["total_overtime_hours"])
}

// Test calculateEmployeeSummary - Zero Values
func TestCalculateEmployeeSummary_ZeroValues(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	empPayslips := []model.Payslip{
		{
			TotalAmount:         0.0,
			BasicSalary:         0.0,
			OvertimeAmount:      0.0,
			ReimbursementAmount: 0.0,
			AttendanceDays:      0,
			OvertimeHours:       0,
		},
		{
			TotalAmount:         0.0,
			BasicSalary:         0.0,
			OvertimeAmount:      0.0,
			ReimbursementAmount: 0.0,
			AttendanceDays:      0,
			OvertimeHours:       0,
		},
	}

	result := usecase.calculateEmployeeSummary(uint(3), empPayslips, "Test Employee")

	assert.Equal(t, uint(3), result["employee_id"])
	assert.Equal(t, "Test Employee", result["employee_name"])
	assert.Equal(t, 2, result["payslip_count"])
	assert.Equal(t, 0.0, result["total_take_home_pay"])
	assert.Equal(t, 0.0, result["total_basic_salary"])
	assert.Equal(t, 0.0, result["total_overtime_amount"])
	assert.Equal(t, 0.0, result["total_reimbursement"])
	assert.Equal(t, 0, result["total_attendance_days"])
	assert.Equal(t, 0, result["total_overtime_hours"])
}

// Test calculateSummaryTotals - Success with Multiple Employees
func TestCalculateSummaryTotals_Success(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	employeeSummaries := []map[string]interface{}{
		{
			"employee_id":           uint(1),
			"employee_name":         "John Doe",
			"payslip_count":         2,
			"total_take_home_pay":   11000.0,
			"total_basic_salary":    10000.0,
			"total_overtime_amount": 650.0,
			"total_reimbursement":   350.0,
			"total_attendance_days": 43,
			"total_overtime_hours":  13,
		},
		{
			"employee_id":           uint(2),
			"employee_name":         "Jane Smith",
			"payslip_count":         1,
			"total_take_home_pay":   6650.0,
			"total_basic_salary":    6000.0,
			"total_overtime_amount": 500.0,
			"total_reimbursement":   150.0,
			"total_attendance_days": 23,
			"total_overtime_hours":  10,
		},
	}

	payslips := []model.Payslip{
		{}, {}, {}, // 3 payslips total
	}

	result := usecase.calculateSummaryTotals(
		employeeSummaries,
		payslips,
		17650.0, // totalTakeHomePay (11000 + 6650)
		16000.0, // totalBasicSalary (10000 + 6000)
		1150.0,  // totalOvertimeAmount (650 + 500)
		500.0,   // totalReimbursementAmount (350 + 150)
		66,      // totalAttendanceDays (43 + 23)
		23,      // totalOvertimeHours (13 + 10)
	)

	assert.Equal(t, 2, result["total_employees"])
	assert.Equal(t, 3, result["total_payslips"])
	assert.Equal(t, 17650.0, result["total_take_home_pay"])
	assert.Equal(t, 16000.0, result["total_basic_salary"])
	assert.Equal(t, 1150.0, result["total_overtime_amount"])
	assert.Equal(t, 500.0, result["total_reimbursement_amount"])
	assert.Equal(t, 66, result["total_attendance_days"])
	assert.Equal(t, 23, result["total_overtime_hours"])

	// Check averages (totals / 2 employees)
	assert.Equal(t, 8825.0, result["average_take_home_pay"])  // 17650 / 2
	assert.Equal(t, 8000.0, result["average_basic_salary"])   // 16000 / 2
	assert.Equal(t, 575.0, result["average_overtime_amount"]) // 1150 / 2
	assert.Equal(t, 250.0, result["average_reimbursement"])   // 500 / 2
}

// Test calculateSummaryTotals - Zero Employees
func TestCalculateSummaryTotals_ZeroEmployees(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	employeeSummaries := []map[string]interface{}{}
	payslips := []model.Payslip{}

	result := usecase.calculateSummaryTotals(
		employeeSummaries,
		payslips,
		0.0, 0.0, 0.0, 0.0, 0, 0,
	)

	assert.Equal(t, 0, result["total_employees"])
	assert.Equal(t, 0, result["total_payslips"])
	assert.Equal(t, 0.0, result["total_take_home_pay"])
	assert.Equal(t, 0.0, result["total_basic_salary"])
	assert.Equal(t, 0.0, result["total_overtime_amount"])
	assert.Equal(t, 0.0, result["total_reimbursement_amount"])
	assert.Equal(t, 0, result["total_attendance_days"])
	assert.Equal(t, 0, result["total_overtime_hours"])

	// Averages should be 0.0 when no employees
	assert.Equal(t, 0.0, result["average_take_home_pay"])
	assert.Equal(t, 0.0, result["average_basic_salary"])
	assert.Equal(t, 0.0, result["average_overtime_amount"])
	assert.Equal(t, 0.0, result["average_reimbursement"])
}

// Test calculateSummaryTotals - Single Employee
func TestCalculateSummaryTotals_SingleEmployee(t *testing.T) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	employeeSummaries := []map[string]interface{}{
		{
			"employee_id":           uint(1),
			"employee_name":         "John Doe",
			"payslip_count":         1,
			"total_take_home_pay":   5600.0,
			"total_basic_salary":    5000.0,
			"total_overtime_amount": 400.0,
			"total_reimbursement":   200.0,
			"total_attendance_days": 22,
			"total_overtime_hours":  8,
		},
	}

	payslips := []model.Payslip{{}}

	result := usecase.calculateSummaryTotals(
		employeeSummaries,
		payslips,
		5600.0, 5000.0, 400.0, 200.0, 22, 8,
	)

	assert.Equal(t, 1, result["total_employees"])
	assert.Equal(t, 1, result["total_payslips"])

	// Averages should equal totals when only one employee
	assert.Equal(t, 5600.0, result["average_take_home_pay"])
	assert.Equal(t, 5000.0, result["average_basic_salary"])
	assert.Equal(t, 400.0, result["average_overtime_amount"])
	assert.Equal(t, 200.0, result["average_reimbursement"])
}

// Benchmark tests for helper functions
func BenchmarkCalculateTotalOvertimeHours(b *testing.B) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	// Create test data
	overtimes := make([]model.Overtime, 100)
	for i := 0; i < 100; i++ {
		overtimes[i] = model.Overtime{Hours: i % 10}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = usecase.calculateTotalOvertimeHours(overtimes)
	}
}

func BenchmarkCalculateTotalReimbursementAmount(b *testing.B) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	// Create test data
	reimbursements := make([]model.Reimbursement, 100)
	for i := 0; i < 100; i++ {
		reimbursements[i] = model.Reimbursement{Amount: float64(i) * 10.5}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = usecase.calculateTotalReimbursementAmount(reimbursements)
	}
}

func BenchmarkBuildAttendanceBreakdown(b *testing.B) {
	mockPayslipRepo := new(MockPayslipRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	usecase := NewPayrollUsecase(mockPayslipRepo, mockEmployeeRepo)

	// Create test data
	attendances := make([]model.Attendance, 30) // One month of attendance
	for i := 0; i < 30; i++ {
		checkout := time.Date(2025, 6, i+1, 17, 0, 0, 0, time.UTC)
		attendances[i] = model.Attendance{
			Date:        time.Date(2025, 6, i+1, 0, 0, 0, 0, time.UTC),
			Checkin:     time.Date(2025, 6, i+1, 9, 0, 0, 0, time.UTC),
			Checkout:    &checkout,
			HoursWorked: 8,
			Status:      "present",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = usecase.buildAttendanceBreakdown(attendances)
	}
}
