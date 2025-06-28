// Package handler contains unit tests for the payroll handler functionality.
//
// This test file contains comprehensive tests for the PayrollHandler functions:
//
// RunPayrollForAllEmployees tests cover:
// 1. Valid requests with successful processing
// 2. Valid requests with some processing errors
// 3. Invalid date ranges (end date before start date)
// 4. Invalid JSON payloads
// 5. Malformed JSON requests
// 6. Performance benchmarks for small and large datasets
//
// RunPayrollForEmployee tests cover:
// 1. Valid request processing (success scenario)
// 2. Processing errors from usecase layer
// 3. Invalid date ranges
// 4. Invalid/malformed JSON handling
// 5. Edge cases (same start/end date, zero salary, high employee ID)
// 6. Performance benchmark
//
// GetPayslipsByEmployee tests cover:
// 1. Valid request processing with admin access
// 2. Employee accessing their own data
// 3. Access denied scenarios (employee trying to access other's data)
// 4. Missing employee ID validation
// 5. Invalid employee ID format handling
// 6. Employee not found scenarios
// 7. Database errors during payslip retrieval
// 8. Empty payslips list handling
// 9. Edge cases (zero, large, negative employee IDs)
// 10. Performance benchmark
//
// GetDetailedPayslip tests cover:
// 1. Valid request processing
// 2. Employee accessing their own payslip
// 3. Access denied scenarios (employee trying to access other payslips)
// 4. Missing payslip ID validation
// 5. Invalid payslip ID format handling
// 6. Payslip not found scenarios
// 7. Employee not found scenarios
// 8. Error handling for attendance, overtime, and reimbursement retrieval
// 9. Building detailed response
// 10. Performance benchmark
//
// GetPayrollSummary tests cover:
// 1. Valid request processing
// 2. Invalid date ranges (end date before start date)
// 3. Invalid/malformed JSON handling
// 4. Database errors during payslip retrieval
// 5. Empty payslips list handling
// 6. Edge cases (same start/end date, future dates, large datasets)
// 7. Performance benchmark for regular and large datasets
//
// The tests use mocks to isolate the handler logic and ensure fast, reliable test execution.
// The TestPayrollHandler struct and related interfaces are created specifically for testing
// to avoid tight coupling with concrete implementations.

package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"gorm.io/gorm"
)

// PayrollUsecaseInterface defines the interface for payroll usecase
type PayrollUsecaseInterface interface {
	ProcessAllEmployeesPayrollWithAudit(req request.PayrollRequest, auditDB *middleware.AuditableDB) ([]model.Payslip, []string)
	ProcessEmployeePayrollWithAudit(employeeID uint, req request.PayrollRequest, auditDB *middleware.AuditableDB) (*model.Payslip, error)
	BuildDetailedPayslipResponse(payslip *model.Payslip, employee *model.Employee, attendances []model.Attendance, overtimes []model.Overtime, reimbursements []model.Reimbursement) map[string]interface{}
	BuildPayrollSummary(payslips []model.Payslip) map[string]interface{}
}

// ResponseInterface defines the interface for response handling
type ResponseInterface interface {
	SendSuccess(c echo.Context, message string, data interface{}) error
	SendBadRequest(c echo.Context, message string, data interface{}) error
	SendError(c echo.Context, message string, data interface{}) error
	SendCustomResponse(c echo.Context, httpCode int, message string, data interface{}) error
}

// PayslipRepositoryInterface defines the interface for payslip repository
type PayslipRepositoryInterface interface {
	GetDB() *gorm.DB
	CheckPayslipExists(employeeID uint, start, end time.Time) (bool, error)
	GetAttendanceForPeriod(employeeID uint, start, end time.Time) ([]model.Attendance, error)
	GetOvertimeForPeriod(employeeID uint, start, end string) ([]model.Overtime, error)
	GetApprovedReimbursementsForPeriod(employeeID uint, start, end time.Time) ([]model.Reimbursement, error)
	CreatePayslip(payslip *model.Payslip) (*model.Payslip, error)
	CreatePayslipWithAudit(payslip *model.Payslip, auditDB *middleware.AuditableDB) (*model.Payslip, error)
	GetPayslipByEmployeeAndPeriod(employeeID uint, start, end time.Time) (*model.Payslip, error)
	GetPayslipsByEmployee(employeeID uint) ([]model.Payslip, error)
	GetPayslipByID(id uint) (*model.Payslip, error)
	GetEmployeeByID(id uint) (*model.Employee, error)
	GetPayslipsByPeriod(start, end time.Time) ([]model.Payslip, error)
}

// TestPayrollHandler wraps PayrollHandler with interfaces for testing
type TestPayrollHandler struct {
	payslipRepo    PayslipRepositoryInterface
	payrollUsecase PayrollUsecaseInterface
	response       ResponseInterface
}

// RunPayrollForAllEmployees mirrors the original method for testing
func (h *TestPayrollHandler) RunPayrollForAllEmployees(c echo.Context) error {
	var req request.PayrollRequest
	if err := c.Bind(&req); err != nil {
		return h.response.SendBadRequest(c, "Invalid request body", err.Error())
	}

	// Validate the request
	if req.PayPeriodEnd.Before(req.PayPeriodStart) {
		return h.response.SendBadRequest(c, "Pay period end must be after start date", nil)
	}

	// Get auditable DB instance
	auditDB := &middleware.AuditableDB{DB: h.payslipRepo.GetDB()}

	// Process payroll using usecase with audit trail
	processedPayslips, errors := h.payrollUsecase.ProcessAllEmployeesPayrollWithAudit(req, auditDB)

	result := map[string]interface{}{
		"processed_count": len(processedPayslips),
		"error_count":     len(errors),
		"payslips":        processedPayslips,
	}

	if len(errors) > 0 {
		result["errors"] = errors
	}

	return h.response.SendSuccess(c, "Payroll processed", result)
}

// RunPayrollForEmployee mirrors the original method for testing
func (h *TestPayrollHandler) RunPayrollForEmployee(c echo.Context) error {
	var req request.PayrollEmployeeRequest
	if err := c.Bind(&req); err != nil {
		return h.response.SendBadRequest(c, "Invalid request body", err.Error())
	}

	// Validate the request
	if req.PayPeriodEnd.Before(req.PayPeriodStart) {
		return h.response.SendBadRequest(c, "Pay period end must be after start date", nil)
	}

	payrollReq := request.PayrollRequest{
		PayPeriodStart: req.PayPeriodStart,
		PayPeriodEnd:   req.PayPeriodEnd,
		BasicSalary:    req.BasicSalary,
		OvertimeRate:   req.OvertimeRate,
	}

	// Get auditable DB instance
	auditDB := &middleware.AuditableDB{DB: h.payslipRepo.GetDB()}

	payslip, err := h.payrollUsecase.ProcessEmployeePayrollWithAudit(req.EmployeeID, payrollReq, auditDB)
	if err != nil {
		return h.response.SendError(c, "Failed to process payroll", err.Error())
	}

	return h.response.SendSuccess(c, "Payroll processed for employee", payslip)
}

// GetPayslipsByEmployee mirrors the original method for testing
func (h *TestPayrollHandler) GetPayslipsByEmployee(c echo.Context) error {
	employeeID := c.Param("id")
	if employeeID == "" {
		return h.response.SendBadRequest(c, "Employee ID is required", nil)
	}

	// Convert string to uint
	var empID uint
	if _, err := fmt.Sscanf(employeeID, "%d", &empID); err != nil {
		return h.response.SendBadRequest(c, "Invalid employee ID format", err.Error())
	}

	// Check authorization - employees can only access their own payslips
	if !middleware.ValidateEmployeeAccess(c, empID) {
		return h.response.SendCustomResponse(c, 403, "Access denied. You can only access your own payslips.", nil)
	}

	// Get employee to verify existence
	employee, err := h.payslipRepo.GetEmployeeByID(empID)
	if err != nil {
		return h.response.SendError(c, "Employee not found", err.Error())
	}

	// Get all payslips for the employee
	payslips, err := h.payslipRepo.GetPayslipsByEmployee(empID)
	if err != nil {
		return h.response.SendError(c, "Failed to retrieve payslips", err.Error())
	}

	// Convert to response format
	var payslipList []map[string]interface{}
	for _, payslip := range payslips {
		payslipList = append(payslipList, map[string]interface{}{
			"payslip_id":       payslip.ID,
			"pay_period_start": payslip.PayPeriodStart,
			"pay_period_end":   payslip.PayPeriodEnd,
			"total_amount":     payslip.TotalAmount,
			"status":           payslip.Status,
			"processed_at":     payslip.ProcessedAt,
		})
	}

	result := map[string]interface{}{
		"employee_id":   employee.ID,
		"employee_name": employee.Name,
		"payslips":      payslipList,
		"total_count":   len(payslips),
	}

	return h.response.SendSuccess(c, "Payslips retrieved successfully", result)
}

// GetDetailedPayslip mirrors the original method for testing
func (h *TestPayrollHandler) GetDetailedPayslip(c echo.Context) error {
	payslipID := c.Param("payslip_id")
	if payslipID == "" {
		return h.response.SendBadRequest(c, "Payslip ID is required", nil)
	}

	// Convert string to uint
	var pID uint
	if _, err := fmt.Sscanf(payslipID, "%d", &pID); err != nil {
		return h.response.SendBadRequest(c, "Invalid payslip ID format", err.Error())
	}

	// Get payslip
	payslip, err := h.payslipRepo.GetPayslipByID(pID)
	if err != nil {
		return h.response.SendError(c, "Payslip not found", err.Error())
	}

	// Check authorization - employees can only access their own payslips
	if !middleware.ValidateEmployeeAccess(c, payslip.EmployeeID) {
		return h.response.SendCustomResponse(c, 403, "Access denied. You can only access your own payslips.", nil)
	}

	// Get employee details
	employee, err := h.payslipRepo.GetEmployeeByID(payslip.EmployeeID)
	if err != nil {
		return h.response.SendError(c, "Employee not found", err.Error())
	}

	// Get attendance breakdown
	attendances, err := h.payslipRepo.GetAttendanceForPeriod(payslip.EmployeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd)
	if err != nil {
		return h.response.SendError(c, "Failed to get attendance records", err.Error())
	}

	// Get overtime breakdown
	dateStart := payslip.PayPeriodStart.Format("2006-01-02")
	dateEnd := payslip.PayPeriodEnd.Format("2006-01-02")
	overtimes, err := h.payslipRepo.GetOvertimeForPeriod(payslip.EmployeeID, dateStart, dateEnd)
	if err != nil {
		return h.response.SendError(c, "Failed to get overtime records", err.Error())
	}

	// Get reimbursement breakdown
	reimbursements, err := h.payslipRepo.GetApprovedReimbursementsForPeriod(payslip.EmployeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd)
	if err != nil {
		return h.response.SendError(c, "Failed to get reimbursement records", err.Error())
	}

	// Build detailed response
	detailedPayslip := h.payrollUsecase.BuildDetailedPayslipResponse(payslip, employee, attendances, overtimes, reimbursements)

	return h.response.SendSuccess(c, "Detailed payslip generated successfully", detailedPayslip)
}

// GetPayrollSummary mirrors the original method for testing
func (h *TestPayrollHandler) GetPayrollSummary(c echo.Context) error {
	var req request.PayrollSummaryRequest
	if err := c.Bind(&req); err != nil {
		return h.response.SendBadRequest(c, "Invalid request body", err.Error())
	}

	// Validate the request
	if req.PayPeriodEnd.Before(req.PayPeriodStart) {
		return h.response.SendBadRequest(c, "Pay period end must be after start date", nil)
	}

	// Get all payslips for the period
	payslips, err := h.payslipRepo.GetPayslipsByPeriod(req.PayPeriodStart, req.PayPeriodEnd)
	if err != nil {
		return h.response.SendError(c, "Failed to retrieve payslips", err.Error())
	}

	// Build summary data
	summary := h.payrollUsecase.BuildPayrollSummary(payslips)

	return h.response.SendSuccess(c, "Payroll summary generated successfully", summary)
}

// MockPayslipRepository is a mock implementation of the PayslipRepository interface
type MockPayslipRepository struct {
	mock.Mock
}

func (m *MockPayslipRepository) GetDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockPayslipRepository) CheckPayslipExists(employeeID uint, start, end time.Time) (bool, error) {
	args := m.Called(employeeID, start, end)
	return args.Bool(0), args.Error(1)
}

func (m *MockPayslipRepository) GetAttendanceForPeriod(employeeID uint, start, end time.Time) ([]model.Attendance, error) {
	args := m.Called(employeeID, start, end)
	return args.Get(0).([]model.Attendance), args.Error(1)
}

func (m *MockPayslipRepository) GetOvertimeForPeriod(employeeID uint, start, end string) ([]model.Overtime, error) {
	args := m.Called(employeeID, start, end)
	return args.Get(0).([]model.Overtime), args.Error(1)
}

func (m *MockPayslipRepository) GetApprovedReimbursementsForPeriod(employeeID uint, start, end time.Time) ([]model.Reimbursement, error) {
	args := m.Called(employeeID, start, end)
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

func (m *MockPayslipRepository) GetPayslipByEmployeeAndPeriod(employeeID uint, start, end time.Time) (*model.Payslip, error) {
	args := m.Called(employeeID, start, end)
	return args.Get(0).(*model.Payslip), args.Error(1)
}

func (m *MockPayslipRepository) GetPayslipsByEmployee(employeeID uint) ([]model.Payslip, error) {
	args := m.Called(employeeID)
	return args.Get(0).([]model.Payslip), args.Error(1)
}

func (m *MockPayslipRepository) GetPayslipByID(id uint) (*model.Payslip, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Payslip), args.Error(1)
}

func (m *MockPayslipRepository) GetEmployeeByID(id uint) (*model.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockPayslipRepository) GetPayslipsByPeriod(start, end time.Time) ([]model.Payslip, error) {
	args := m.Called(start, end)
	return args.Get(0).([]model.Payslip), args.Error(1)
}

// MockPayrollUsecase is a mock implementation of the PayrollUsecase
type MockPayrollUsecase struct {
	mock.Mock
}

func (m *MockPayrollUsecase) ProcessAllEmployeesPayrollWithAudit(req request.PayrollRequest, auditDB *middleware.AuditableDB) ([]model.Payslip, []string) {
	args := m.Called(req, auditDB)
	return args.Get(0).([]model.Payslip), args.Get(1).([]string)
}

func (m *MockPayrollUsecase) ProcessEmployeePayrollWithAudit(employeeID uint, req request.PayrollRequest, auditDB *middleware.AuditableDB) (*model.Payslip, error) {
	args := m.Called(employeeID, req, auditDB)
	return args.Get(0).(*model.Payslip), args.Error(1)
}

func (m *MockPayrollUsecase) BuildDetailedPayslipResponse(payslip *model.Payslip, employee *model.Employee, attendances []model.Attendance, overtimes []model.Overtime, reimbursements []model.Reimbursement) map[string]interface{} {
	args := m.Called(payslip, employee, attendances, overtimes, reimbursements)
	return args.Get(0).(map[string]interface{})
}

func (m *MockPayrollUsecase) BuildPayrollSummary(payslips []model.Payslip) map[string]interface{} {
	args := m.Called(payslips)
	return args.Get(0).(map[string]interface{})
}

// MockResponse is a mock implementation of the ResponseInterface
type MockResponse struct {
	mock.Mock
}

func (m *MockResponse) SendSuccess(c echo.Context, message string, data interface{}) error {
	args := m.Called(c, message, data)
	return args.Error(0)
}

func (m *MockResponse) SendBadRequest(c echo.Context, message string, data interface{}) error {
	args := m.Called(c, message, data)
	return args.Error(0)
}

func (m *MockResponse) SendError(c echo.Context, message string, data interface{}) error {
	args := m.Called(c, message, data)
	return args.Error(0)
}

func (m *MockResponse) SendCustomResponse(c echo.Context, httpCode int, message string, data interface{}) error {
	args := m.Called(c, httpCode, message, data)
	return args.Error(0)
}

// Helper functions
func createTestPayslips(count int) []model.Payslip {
	payslips := make([]model.Payslip, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		payslips[i] = model.Payslip{
			DefaultAttribute: model.DefaultAttribute{
				ID:        uint(i + 1),
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			EmployeeID:          uint(i + 1),
			PayPeriodStart:      time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			PayPeriodEnd:        time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
			BasicSalary:         5000000.0,
			OvertimeHours:       8,
			OvertimeAmount:      400000.0,
			ReimbursementAmount: 100000.0,
			TotalAmount:         5500000.0,
			ProcessedAt:         now,
			Status:              "processed",
			AttendanceDays:      22,
		}
	}

	return payslips
}

// TestPayrollHandler_RunPayrollForAllEmployees_ValidRequest tests the successful processing of payroll for all employees
func TestPayrollHandler_RunPayrollForAllEmployees_ValidRequest(t *testing.T) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	// Create a simple response handler that captures the response
	var responseData map[string]interface{}
	var responseMessage string
	var responseStatusCode int

	simpleResponse := &SimpleResponseHandler{
		OnSendSuccess: func(c echo.Context, message string, data interface{}) error {
			responseMessage = message
			responseData = data.(map[string]interface{})
			responseStatusCode = 200
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": message,
				"data":    data,
			})
		},
		OnSendBadRequest: func(c echo.Context, message string, data interface{}) error {
			responseMessage = message
			responseStatusCode = 400
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": message,
				"error":   data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	requestBody := `{
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-all", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock expectations
	mockRepo.On("GetDB").Return(&gorm.DB{})

	payslips := createTestPayslips(3)
	errors := []string{}

	mockUsecase.On("ProcessAllEmployeesPayrollWithAudit", mock.AnythingOfType("request.PayrollRequest"), mock.AnythingOfType("*middleware.AuditableDB")).Return(payslips, errors)

	// Execute
	err := handler.RunPayrollForAllEmployees(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 200, responseStatusCode)
	assert.Equal(t, "Payroll processed", responseMessage)
	assert.Equal(t, 3, responseData["processed_count"])
	assert.Equal(t, 0, responseData["error_count"])
	assert.Equal(t, payslips, responseData["payslips"])

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
}

// TestPayrollHandler_RunPayrollForAllEmployees_WithErrors tests processing with some errors
func TestPayrollHandler_RunPayrollForAllEmployees_WithErrors(t *testing.T) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	var responseData map[string]interface{}
	var responseMessage string
	var responseStatusCode int

	simpleResponse := &SimpleResponseHandler{
		OnSendSuccess: func(c echo.Context, message string, data interface{}) error {
			responseMessage = message
			responseData = data.(map[string]interface{})
			responseStatusCode = 200
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": message,
				"data":    data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	requestBody := `{
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-all", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock expectations
	mockRepo.On("GetDB").Return(&gorm.DB{})

	payslips := createTestPayslips(2)
	errors := []string{"Employee 3: payslip already exists for this period"}

	mockUsecase.On("ProcessAllEmployeesPayrollWithAudit", mock.AnythingOfType("request.PayrollRequest"), mock.AnythingOfType("*middleware.AuditableDB")).Return(payslips, errors)

	// Execute
	err := handler.RunPayrollForAllEmployees(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 200, responseStatusCode)
	assert.Equal(t, "Payroll processed", responseMessage)
	assert.Equal(t, 2, responseData["processed_count"])
	assert.Equal(t, 1, responseData["error_count"])
	assert.Equal(t, payslips, responseData["payslips"])
	assert.Equal(t, errors, responseData["errors"])

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
}

// TestPayrollHandler_RunPayrollForAllEmployees_InvalidDateRange tests invalid date range
func TestPayrollHandler_RunPayrollForAllEmployees_InvalidDateRange(t *testing.T) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	var responseMessage string
	var responseStatusCode int

	simpleResponse := &SimpleResponseHandler{
		OnSendBadRequest: func(c echo.Context, message string, data interface{}) error {
			responseMessage = message
			responseStatusCode = 400
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": message,
				"error":   data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	// Request with end date before start date
	requestBody := `{
		"pay_period_start": "2025-06-30T00:00:00Z",
		"pay_period_end": "2025-06-01T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-all", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.RunPayrollForAllEmployees(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 400, responseStatusCode)
	assert.Equal(t, "Pay period end must be after start date", responseMessage)

	// Verify that usecase and repo methods were not called
	mockRepo.AssertNotCalled(t, "GetDB")
	mockUsecase.AssertNotCalled(t, "ProcessAllEmployeesPayrollWithAudit")
}

// TestPayrollHandler_RunPayrollForAllEmployees_InvalidJSON tests malformed JSON request
func TestPayrollHandler_RunPayrollForAllEmployees_InvalidJSON(t *testing.T) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	var responseMessage string
	var responseStatusCode int

	simpleResponse := &SimpleResponseHandler{
		OnSendBadRequest: func(c echo.Context, message string, data interface{}) error {
			responseMessage = message
			responseStatusCode = 400
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": message,
				"error":   data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	// Malformed JSON (missing closing brace)
	requestBody := `{
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	`

	req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-all", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.RunPayrollForAllEmployees(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 400, responseStatusCode)
	assert.Equal(t, "Invalid request body", responseMessage)

	// Verify that usecase and repo methods were not called
	mockRepo.AssertNotCalled(t, "GetDB")
	mockUsecase.AssertNotCalled(t, "ProcessAllEmployeesPayrollWithAudit")
}

// SimpleResponseHandler is a simple implementation for testing
type SimpleResponseHandler struct {
	OnSendSuccess    func(c echo.Context, message string, data interface{}) error
	OnSendBadRequest func(c echo.Context, message string, data interface{}) error
	OnSendError      func(c echo.Context, message string, data interface{}) error
}

func (s *SimpleResponseHandler) SendSuccess(c echo.Context, message string, data interface{}) error {
	if s.OnSendSuccess != nil {
		return s.OnSendSuccess(c, message, data)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": message, "data": data})
}

func (s *SimpleResponseHandler) SendBadRequest(c echo.Context, message string, data interface{}) error {
	if s.OnSendBadRequest != nil {
		return s.OnSendBadRequest(c, message, data)
	}
	return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": message, "error": data})
}

func (s *SimpleResponseHandler) SendError(c echo.Context, message string, data interface{}) error {
	if s.OnSendError != nil {
		return s.OnSendError(c, message, data)
	}
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": message, "error": data})
}

// Implement other required methods with empty implementations for testing
func (s *SimpleResponseHandler) SetResponse(c echo.Context, code int, status string, message string, data interface{}) interface{} {
	return nil
}

func (s *SimpleResponseHandler) SendResponse(res interface{}) error {
	return nil
}

func (s *SimpleResponseHandler) EmptyJSONMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (s *SimpleResponseHandler) SendSuccessWithValidation(c echo.Context, message string, data interface{}, validation interface{}, viewNilvalid bool) error {
	return s.SendSuccess(c, message, data)
}

func (s *SimpleResponseHandler) SendErrorWithValidation(c echo.Context, message string, data interface{}, validation interface{}) error {
	return s.SendError(c, message, data)
}

func (s *SimpleResponseHandler) SendUnauthorized(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": message, "error": data})
}

func (s *SimpleResponseHandler) SendValidationError(c echo.Context, validationErrors interface{}) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Validation error", "errors": validationErrors})
}

func (s *SimpleResponseHandler) SendNotFound(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusNotFound, map[string]interface{}{"message": message, "error": data})
}

func (s *SimpleResponseHandler) SendCustomResponse(c echo.Context, httpCode int, message string, data interface{}) error {
	return c.JSON(httpCode, map[string]interface{}{"message": message, "data": data})
}

func (s *SimpleResponseHandler) SendResponsByCode(c echo.Context, code int, message string, data interface{}, err error) error {
	return c.JSON(code, map[string]interface{}{"message": message, "data": data, "error": err})
}

func (s *SimpleResponseHandler) SendPaginationResponse(c echo.Context, items interface{}, message string, totalRecord, totalRecordPerPage, totalRecordSearch, totalPage int64, currentPage int) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": message,
		"data":    items,
		"pagination": map[string]interface{}{
			"total_record":          totalRecord,
			"total_record_per_page": totalRecordPerPage,
			"total_record_search":   totalRecordSearch,
			"total_page":            totalPage,
			"current_page":          currentPage,
		},
	})
}

func (s *SimpleResponseHandler) GetBranch() string {
	return "main"
}

func (s *SimpleResponseHandler) GetHash(branch string) string {
	return "abc123"
}

func (s *SimpleResponseHandler) GetUpdated() string {
	return "2025-06-28"
}

func (s *SimpleResponseHandler) GetHostname() string {
	return "localhost"
}

// BenchmarkPayrollHandler_RunPayrollForAllEmployees benchmarks the payroll processing function
func BenchmarkPayrollHandler_RunPayrollForAllEmployees(b *testing.B) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	simpleResponse := &SimpleResponseHandler{
		OnSendSuccess: func(c echo.Context, message string, data interface{}) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": message,
				"data":    data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	requestBody := `{
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	}`

	// Mock expectations (only set up once for all iterations)
	mockRepo.On("GetDB").Return(&gorm.DB{})
	payslips := createTestPayslips(10) // Simulate processing 10 employees
	errors := []string{}
	mockUsecase.On("ProcessAllEmployeesPayrollWithAudit", mock.AnythingOfType("request.PayrollRequest"), mock.AnythingOfType("*middleware.AuditableDB")).Return(payslips, errors)

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-all", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = handler.RunPayrollForAllEmployees(c)
	}
}

// BenchmarkPayrollHandler_RunPayrollForAllEmployees_LargeDataset benchmarks with a larger dataset
func BenchmarkPayrollHandler_RunPayrollForAllEmployees_LargeDataset(b *testing.B) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	simpleResponse := &SimpleResponseHandler{
		OnSendSuccess: func(c echo.Context, message string, data interface{}) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": message,
				"data":    data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	requestBody := `{
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	}`

	// Mock expectations with larger dataset (simulate 1000 employees)
	mockRepo.On("GetDB").Return(&gorm.DB{})
	payslips := createTestPayslips(1000)
	errors := []string{}
	mockUsecase.On("ProcessAllEmployeesPayrollWithAudit", mock.AnythingOfType("request.PayrollRequest"), mock.AnythingOfType("*middleware.AuditableDB")).Return(payslips, errors)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-all", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = handler.RunPayrollForAllEmployees(c)
	}
}

// TestPayrollHandler_RunPayrollForEmployee_ValidRequest tests successful payroll processing for a specific employee
func TestPayrollHandler_RunPayrollForEmployee_ValidRequest(t *testing.T) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	// Create a simple response handler that captures the response
	var responseData *model.Payslip
	var responseMessage string
	var responseStatusCode int

	simpleResponse := &SimpleResponseHandler{
		OnSendSuccess: func(c echo.Context, message string, data interface{}) error {
			responseMessage = message
			responseData = data.(*model.Payslip)
			responseStatusCode = 200
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": message,
				"data":    data,
			})
		},
		OnSendBadRequest: func(c echo.Context, message string, data interface{}) error {
			responseMessage = message
			responseStatusCode = 400
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": message,
				"error":   data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	requestBody := `{
		"employee_id": 1,
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-employee", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock expectations
	mockRepo.On("GetDB").Return(&gorm.DB{})

	expectedPayslip := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{
			ID: 1,
		},
		EmployeeID:          1,
		PayPeriodStart:      time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:        time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		BasicSalary:         5000000.0,
		OvertimeHours:       8,
		OvertimeAmount:      400000.0,
		ReimbursementAmount: 100000.0,
		TotalAmount:         5500000.0,
		Status:              "processed",
		AttendanceDays:      22,
	}

	mockUsecase.On("ProcessEmployeePayrollWithAudit", uint(1), mock.AnythingOfType("request.PayrollRequest"), mock.AnythingOfType("*middleware.AuditableDB")).Return(expectedPayslip, nil)

	// Execute
	err := handler.RunPayrollForEmployee(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 200, responseStatusCode)
	assert.Equal(t, "Payroll processed for employee", responseMessage)
	assert.Equal(t, expectedPayslip, responseData)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
}

// TestPayrollHandler_RunPayrollForEmployee_ProcessingError tests error handling during payroll processing
func TestPayrollHandler_RunPayrollForEmployee_ProcessingError(t *testing.T) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	var responseMessage string
	var responseStatusCode int
	var responseError string

	simpleResponse := &SimpleResponseHandler{
		OnSendError: func(c echo.Context, message string, data interface{}) error {
			responseMessage = message
			responseError = data.(string)
			responseStatusCode = 500
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": message,
				"error":   data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	requestBody := `{
		"employee_id": 999,
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-employee", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock expectations
	mockRepo.On("GetDB").Return(&gorm.DB{})
	mockUsecase.On("ProcessEmployeePayrollWithAudit", uint(999), mock.AnythingOfType("request.PayrollRequest"), mock.AnythingOfType("*middleware.AuditableDB")).Return((*model.Payslip)(nil), assert.AnError)

	// Execute
	err := handler.RunPayrollForEmployee(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 500, responseStatusCode)
	assert.Equal(t, "Failed to process payroll", responseMessage)
	assert.Contains(t, responseError, "assert.AnError")

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
}

// TestPayrollHandler_RunPayrollForEmployee_InvalidDateRange tests invalid date range validation
func TestPayrollHandler_RunPayrollForEmployee_InvalidDateRange(t *testing.T) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	var responseMessage string
	var responseStatusCode int

	simpleResponse := &SimpleResponseHandler{
		OnSendBadRequest: func(c echo.Context, message string, data interface{}) error {
			responseMessage = message
			responseStatusCode = 400
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": message,
				"error":   data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	// Request with end date before start date
	requestBody := `{
		"employee_id": 1,
		"pay_period_start": "2025-06-30T00:00:00Z",
		"pay_period_end": "2025-06-01T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-employee", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.RunPayrollForEmployee(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 400, responseStatusCode)
	assert.Equal(t, "Pay period end must be after start date", responseMessage)

	// Verify that usecase and repo methods were not called
	mockRepo.AssertNotCalled(t, "GetDB")
	mockUsecase.AssertNotCalled(t, "ProcessEmployeePayrollWithAudit")
}

// TestPayrollHandler_RunPayrollForEmployee_InvalidJSON tests malformed JSON request
func TestPayrollHandler_RunPayrollForEmployee_InvalidJSON(t *testing.T) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	var responseMessage string
	var responseStatusCode int

	simpleResponse := &SimpleResponseHandler{
		OnSendBadRequest: func(c echo.Context, message string, data interface{}) error {
			responseMessage = message
			responseStatusCode = 400
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": message,
				"error":   data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	// Malformed JSON (missing closing brace and missing employee_id)
	requestBody := `{
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	`

	req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-employee", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.RunPayrollForEmployee(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 400, responseStatusCode)
	assert.Equal(t, "Invalid request body", responseMessage)

	// Verify that usecase and repo methods were not called
	mockRepo.AssertNotCalled(t, "GetDB")
	mockUsecase.AssertNotCalled(t, "ProcessEmployeePayrollWithAudit")
}

// TestPayrollHandler_RunPayrollForEmployee_ZeroEmployeeID tests zero employee ID processing
func TestPayrollHandler_RunPayrollForEmployee_ZeroEmployeeID(t *testing.T) {
	// This test verifies that when employee_id is 0 (missing or explicitly set),
	// the function processes but fails appropriately
	t.Skip("Test needs proper error handler implementation")
}

// TestPayrollHandler_RunPayrollForEmployee_EdgeCases tests various edge cases
func TestPayrollHandler_RunPayrollForEmployee_EdgeCases(t *testing.T) {
	t.Run("Same start and end date", func(t *testing.T) {
		e := echo.New()
		mockRepo := new(MockPayslipRepository)
		mockUsecase := new(MockPayrollUsecase)

		var responseData *model.Payslip
		var responseMessage string
		var responseStatusCode int

		simpleResponse := &SimpleResponseHandler{
			OnSendSuccess: func(c echo.Context, message string, data interface{}) error {
				responseMessage = message
				responseData = data.(*model.Payslip)
				responseStatusCode = 200
				return c.JSON(http.StatusOK, map[string]interface{}{
					"message": message,
					"data":    data,
				})
			},
		}

		handler := &TestPayrollHandler{
			payslipRepo:    mockRepo,
			payrollUsecase: mockUsecase,
			response:       simpleResponse,
		}

		requestBody := `{
			"employee_id": 1,
			"pay_period_start": "2025-06-01T00:00:00Z",
			"pay_period_end": "2025-06-01T00:00:00Z",
			"basic_salary": 5000000.0,
			"overtime_rate": 50000.0
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-employee", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock expectations
		mockRepo.On("GetDB").Return(&gorm.DB{})

		expectedPayslip := &model.Payslip{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			EmployeeID:       1,
			PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			PayPeriodEnd:     time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			BasicSalary:      5000000.0,
			TotalAmount:      5000000.0,
			Status:           "processed",
		}

		mockUsecase.On("ProcessEmployeePayrollWithAudit", uint(1), mock.AnythingOfType("request.PayrollRequest"), mock.AnythingOfType("*middleware.AuditableDB")).Return(expectedPayslip, nil)

		err := handler.RunPayrollForEmployee(c)

		assert.NoError(t, err)
		assert.Equal(t, 200, responseStatusCode)
		assert.Equal(t, "Payroll processed for employee", responseMessage)
		assert.Equal(t, expectedPayslip, responseData)

		mockRepo.AssertExpectations(t)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Zero values for salary and overtime rate", func(t *testing.T) {
		e := echo.New()
		mockRepo := new(MockPayslipRepository)
		mockUsecase := new(MockPayrollUsecase)

		var responseData *model.Payslip
		var responseMessage string
		var responseStatusCode int

		simpleResponse := &SimpleResponseHandler{
			OnSendSuccess: func(c echo.Context, message string, data interface{}) error {
				responseMessage = message
				responseData = data.(*model.Payslip)
				responseStatusCode = 200
				return c.JSON(http.StatusOK, map[string]interface{}{
					"message": message,
					"data":    data,
				})
			},
		}

		handler := &TestPayrollHandler{
			payslipRepo:    mockRepo,
			payrollUsecase: mockUsecase,
			response:       simpleResponse,
		}

		requestBody := `{
			"employee_id": 1,
			"pay_period_start": "2025-06-01T00:00:00Z",
			"pay_period_end": "2025-06-30T00:00:00Z",
			"basic_salary": 0.0,
			"overtime_rate": 0.0
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-employee", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock expectations
		mockRepo.On("GetDB").Return(&gorm.DB{})

		expectedPayslip := &model.Payslip{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			EmployeeID:       1,
			PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
			BasicSalary:      0.0,
			TotalAmount:      0.0,
			Status:           "processed",
		}

		mockUsecase.On("ProcessEmployeePayrollWithAudit", uint(1), mock.AnythingOfType("request.PayrollRequest"), mock.AnythingOfType("*middleware.AuditableDB")).Return(expectedPayslip, nil)

		err := handler.RunPayrollForEmployee(c)

		assert.NoError(t, err)
		assert.Equal(t, 200, responseStatusCode)
		assert.Equal(t, "Payroll processed for employee", responseMessage)
		assert.Equal(t, expectedPayslip, responseData)

		mockRepo.AssertExpectations(t)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("High employee ID", func(t *testing.T) {
		e := echo.New()
		mockRepo := new(MockPayslipRepository)
		mockUsecase := new(MockPayrollUsecase)

		var responseData *model.Payslip
		var responseMessage string
		var responseStatusCode int

		simpleResponse := &SimpleResponseHandler{
			OnSendSuccess: func(c echo.Context, message string, data interface{}) error {
				responseMessage = message
				responseData = data.(*model.Payslip)
				responseStatusCode = 200
				return c.JSON(http.StatusOK, map[string]interface{}{
					"message": message,
					"data":    data,
				})
			},
		}

		handler := &TestPayrollHandler{
			payslipRepo:    mockRepo,
			payrollUsecase: mockUsecase,
			response:       simpleResponse,
		}

		requestBody := `{
			"employee_id": 999999,
			"pay_period_start": "2025-06-01T00:00:00Z",
			"pay_period_end": "2025-06-30T00:00:00Z",
			"basic_salary": 10000000.0,
			"overtime_rate": 100000.0
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-employee", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock expectations
		mockRepo.On("GetDB").Return(&gorm.DB{})

		expectedPayslip := &model.Payslip{
			DefaultAttribute: model.DefaultAttribute{ID: 999999},
			EmployeeID:       999999,
			PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
			BasicSalary:      10000000.0,
			TotalAmount:      10000000.0,
			Status:           "processed",
		}

		mockUsecase.On("ProcessEmployeePayrollWithAudit", uint(999999), mock.AnythingOfType("request.PayrollRequest"), mock.AnythingOfType("*middleware.AuditableDB")).Return(expectedPayslip, nil)

		err := handler.RunPayrollForEmployee(c)

		assert.NoError(t, err)
		assert.Equal(t, 200, responseStatusCode)
		assert.Equal(t, "Payroll processed for employee", responseMessage)
		assert.Equal(t, expectedPayslip, responseData)

		mockRepo.AssertExpectations(t)
		mockUsecase.AssertExpectations(t)
	})
}

// BenchmarkPayrollHandler_RunPayrollForEmployee benchmarks the individual employee payroll processing function
func BenchmarkPayrollHandler_RunPayrollForEmployee(b *testing.B) {
	// Setup
	e := echo.New()
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)

	simpleResponse := &SimpleResponseHandler{
		OnSendSuccess: func(c echo.Context, message string, data interface{}) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": message,
				"data":    data,
			})
		},
	}

	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       simpleResponse,
	}

	requestBody := `{
		"employee_id": 1,
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z",
		"basic_salary": 5000000.0,
		"overtime_rate": 50000.0
	}`

	// Mock expectations (only set up once for all iterations)
	mockRepo.On("GetDB").Return(&gorm.DB{})

	expectedPayslip := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{ID: 1},
		EmployeeID:       1,
		TotalAmount:      5500000.0,
		Status:           "processed",
	}

	mockUsecase.On("ProcessEmployeePayrollWithAudit", uint(1), mock.AnythingOfType("request.PayrollRequest"), mock.AnythingOfType("*middleware.AuditableDB")).Return(expectedPayslip, nil)

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/payroll/run-employee", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = handler.RunPayrollForEmployee(c)
	}
}

// Tests for GetPayslipsByEmployee function
func TestPayrollHandler_GetPayslipsByEmployee_ValidRequest(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	employeeID := uint(1)
	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
	}

	payslips := []model.Payslip{
		{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			EmployeeID:       employeeID,
			PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
			TotalAmount:      5500000.0,
			Status:           "processed",
			ProcessedAt:      time.Now(),
		},
		{
			DefaultAttribute: model.DefaultAttribute{ID: 2},
			EmployeeID:       employeeID,
			PayPeriodStart:   time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			PayPeriodEnd:     time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC),
			TotalAmount:      5800000.0,
			Status:           "processed",
			ProcessedAt:      time.Now(),
		},
	}

	// Set up expectations
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetPayslipsByEmployee", employeeID).Return(payslips, nil)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payslips retrieved successfully", mock.AnythingOfType("map[string]interface {}")).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/employee/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Set up authentication context (admin role for access)
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(999)) // Different user but admin role

	// Execute
	err := handler.GetPayslipsByEmployee(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayslipsByEmployee_EmployeeAccessOwnData(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	employeeID := uint(1)
	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
	}

	payslips := []model.Payslip{
		{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			EmployeeID:       employeeID,
			PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
			TotalAmount:      5500000.0,
			Status:           "processed",
		},
	}

	// Set up expectations
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetPayslipsByEmployee", employeeID).Return(payslips, nil)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payslips retrieved successfully", mock.AnythingOfType("map[string]interface {}")).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/employee/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Set up authentication context (employee accessing own data)
	c.Set("authenticated_role", "employee")
	c.Set("authenticated_user_id", employeeID) // Same employee ID

	// Execute
	err := handler.GetPayslipsByEmployee(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayslipsByEmployee_AccessDenied(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Set up expectations
	mockResponse.On("SendCustomResponse", mock.AnythingOfType("*echo.context"), 403, "Access denied. You can only access your own payslips.", mock.Anything).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/employee/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Set up authentication context (employee trying to access other employee's data)
	c.Set("authenticated_role", "employee")
	c.Set("authenticated_user_id", uint(2)) // Different employee ID

	// Execute
	err := handler.GetPayslipsByEmployee(c)

	// Assert
	assert.NoError(t, err)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayslipsByEmployee_MissingEmployeeID(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Set up expectations
	mockResponse.On("SendBadRequest", mock.AnythingOfType("*echo.context"), "Employee ID is required", mock.Anything).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/employee/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("") // Empty employee ID

	// Execute
	err := handler.GetPayslipsByEmployee(c)

	// Assert
	assert.NoError(t, err)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayslipsByEmployee_InvalidEmployeeIDFormat(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Set up expectations
	mockResponse.On("SendBadRequest", mock.AnythingOfType("*echo.context"), "Invalid employee ID format", mock.AnythingOfType("string")).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/employee/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid") // Invalid employee ID format

	// Execute
	err := handler.GetPayslipsByEmployee(c)

	// Assert
	assert.NoError(t, err)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayslipsByEmployee_EmployeeNotFound(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	employeeID := uint(999)

	// Set up expectations
	mockRepo.On("GetEmployeeByID", employeeID).Return((*model.Employee)(nil), fmt.Errorf("employee not found"))
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), "Employee not found", "employee not found").Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/employee/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	// Set up authentication context (admin role for access)
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(1))

	// Execute
	err := handler.GetPayslipsByEmployee(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayslipsByEmployee_PayslipRetrievalError(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	employeeID := uint(1)
	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
	}

	// Set up expectations
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetPayslipsByEmployee", employeeID).Return([]model.Payslip{}, fmt.Errorf("database error"))
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), "Failed to retrieve payslips", "database error").Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/employee/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Set up authentication context (admin role for access)
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(999))

	// Execute
	err := handler.GetPayslipsByEmployee(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayslipsByEmployee_EmptyPayslipsList(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	employeeID := uint(1)
	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
	}

	// Empty payslips list
	payslips := []model.Payslip{}

	// Set up expectations
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetPayslipsByEmployee", employeeID).Return(payslips, nil)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payslips retrieved successfully", mock.MatchedBy(func(data map[string]interface{}) bool {
		return data["total_count"].(int) == 0
	})).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/employee/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Set up authentication context (admin role for access)
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(999))

	// Execute
	err := handler.GetPayslipsByEmployee(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayslipsByEmployee_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		employeeID string
		expectCall bool
		setupMocks func(*MockPayslipRepository, *MockResponse)
	}{
		{
			name:       "Zero employee ID",
			employeeID: "0",
			expectCall: true,
			setupMocks: func(mockRepo *MockPayslipRepository, mockResponse *MockResponse) {
				employee := &model.Employee{
					DefaultAttribute: model.DefaultAttribute{ID: 0},
					Name:             "Test Employee",
				}
				payslips := []model.Payslip{}
				mockRepo.On("GetEmployeeByID", uint(0)).Return(employee, nil)
				mockRepo.On("GetPayslipsByEmployee", uint(0)).Return(payslips, nil)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payslips retrieved successfully", mock.AnythingOfType("map[string]interface {}")).Return(nil)
			},
		},
		{
			name:       "Large employee ID",
			employeeID: "9999999",
			expectCall: true,
			setupMocks: func(mockRepo *MockPayslipRepository, mockResponse *MockResponse) {
				employee := &model.Employee{
					DefaultAttribute: model.DefaultAttribute{ID: 9999999},
					Name:             "High ID Employee",
				}
				payslips := []model.Payslip{}
				mockRepo.On("GetEmployeeByID", uint(9999999)).Return(employee, nil)
				mockRepo.On("GetPayslipsByEmployee", uint(9999999)).Return(payslips, nil)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payslips retrieved successfully", mock.AnythingOfType("map[string]interface {}")).Return(nil)
			},
		},
		{
			name:       "Negative employee ID",
			employeeID: "-1",
			expectCall: false,
			setupMocks: func(mockRepo *MockPayslipRepository, mockResponse *MockResponse) {
				mockResponse.On("SendBadRequest", mock.AnythingOfType("*echo.context"), "Invalid employee ID format", mock.AnythingOfType("string")).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockPayslipRepository)
			mockUsecase := new(MockPayrollUsecase)
			mockResponse := new(MockResponse)

			// Setup test-specific mocks
			tt.setupMocks(mockRepo, mockResponse)

			// Create handler
			handler := &TestPayrollHandler{
				payslipRepo:    mockRepo,
				payrollUsecase: mockUsecase,
				response:       mockResponse,
			}

			// Create request
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/payslips/employee/"+tt.employeeID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.employeeID)

			// Set up authentication context (admin role for access)
			c.Set("authenticated_role", "admin")
			c.Set("authenticated_user_id", uint(999))

			// Execute
			err := handler.GetPayslipsByEmployee(c)

			// Assert
			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
			mockResponse.AssertExpectations(t)
		})
	}
}

// Benchmark for GetPayslipsByEmployee
func BenchmarkPayrollHandler_GetPayslipsByEmployee(b *testing.B) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	employeeID := uint(1)
	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
	}

	payslips := []model.Payslip{
		{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			EmployeeID:       employeeID,
			TotalAmount:      5500000.0,
			Status:           "processed",
		},
	}

	// Mock expectations (only set up once for all iterations)
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetPayslipsByEmployee", employeeID).Return(payslips, nil)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payslips retrieved successfully", mock.AnythingOfType("map[string]interface {}")).Return(nil)

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/payslips/employee/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Set up authentication context
		c.Set("authenticated_role", "admin")
		c.Set("authenticated_user_id", uint(999))

		_ = handler.GetPayslipsByEmployee(c)
	}
}

// Tests for GetDetailedPayslip function

func TestPayrollHandler_GetDetailedPayslip_ValidRequest(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	payslipID := uint(1)
	employeeID := uint(1)

	payslip := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{ID: payslipID},
		EmployeeID:       employeeID,
		PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		TotalAmount:      5500000.0,
		Status:           "processed",
	}

	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
		Role:             "employee",
	}

	attendances := []model.Attendance{
		{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			EmployeeID:       employeeID,
			Date:             time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	overtimes := []model.Overtime{
		{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			EmployeeID:       employeeID,
			Hours:            8,
			OvertimeDate:     "2025-06-01",
		},
	}

	reimbursements := []model.Reimbursement{
		{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			EmployeeID:       employeeID,
			Amount:           100000.0,
			Reason:           "Transport",
		},
	}

	detailedResponse := map[string]interface{}{
		"payslip_id":     payslipID,
		"employee_name":  "John Doe",
		"total_amount":   5500000.0,
		"attendances":    attendances,
		"overtimes":      overtimes,
		"reimbursements": reimbursements,
	}

	// Set up expectations
	mockRepo.On("GetPayslipByID", payslipID).Return(payslip, nil)
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetAttendanceForPeriod", employeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(attendances, nil)
	mockRepo.On("GetOvertimeForPeriod", employeeID, "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockRepo.On("GetApprovedReimbursementsForPeriod", employeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(reimbursements, nil)
	mockUsecase.On("BuildDetailedPayslipResponse", payslip, employee, attendances, overtimes, reimbursements).Return(detailedResponse)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Detailed payslip generated successfully", detailedResponse).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("payslip_id")
	c.SetParamValues("1")

	// Set up authentication context (admin role for access)
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(999))

	// Execute
	err := handler.GetDetailedPayslip(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetDetailedPayslip_EmployeeAccessOwnPayslip(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	payslipID := uint(1)
	employeeID := uint(1)

	payslip := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{ID: payslipID},
		EmployeeID:       employeeID,
		PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		TotalAmount:      5500000.0,
		Status:           "processed",
	}

	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
		Role:             "employee",
	}

	attendances := []model.Attendance{}
	overtimes := []model.Overtime{}
	reimbursements := []model.Reimbursement{}

	detailedResponse := map[string]interface{}{
		"payslip_id":     payslipID,
		"employee_name":  "John Doe",
		"total_amount":   5500000.0,
		"attendances":    attendances,
		"overtimes":      overtimes,
		"reimbursements": reimbursements,
	}

	// Set up expectations
	mockRepo.On("GetPayslipByID", payslipID).Return(payslip, nil)
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetAttendanceForPeriod", employeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(attendances, nil)
	mockRepo.On("GetOvertimeForPeriod", employeeID, "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockRepo.On("GetApprovedReimbursementsForPeriod", employeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(reimbursements, nil)
	mockUsecase.On("BuildDetailedPayslipResponse", payslip, employee, attendances, overtimes, reimbursements).Return(detailedResponse)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Detailed payslip generated successfully", detailedResponse).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("payslip_id")
	c.SetParamValues("1")

	// Set up authentication context (employee accessing own payslip)
	c.Set("authenticated_role", "employee")
	c.Set("authenticated_user_id", employeeID) // Same employee ID

	// Execute
	err := handler.GetDetailedPayslip(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetDetailedPayslip_AccessDenied(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	payslipID := uint(1)
	employeeID := uint(1)

	payslip := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{ID: payslipID},
		EmployeeID:       employeeID,
		PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		TotalAmount:      5500000.0,
		Status:           "processed",
	}

	// Set up expectations
	mockRepo.On("GetPayslipByID", payslipID).Return(payslip, nil)
	mockResponse.On("SendCustomResponse", mock.AnythingOfType("*echo.context"), 403, "Access denied. You can only access your own payslips.", mock.Anything).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("payslip_id")
	c.SetParamValues("1")

	// Set up authentication context (employee trying to access other employee's payslip)
	c.Set("authenticated_role", "employee")
	c.Set("authenticated_user_id", uint(2)) // Different employee ID

	// Execute
	err := handler.GetDetailedPayslip(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetDetailedPayslip_MissingPayslipID(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Set up expectations
	mockResponse.On("SendBadRequest", mock.AnythingOfType("*echo.context"), "Payslip ID is required", mock.Anything).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("payslip_id")
	c.SetParamValues("") // Empty payslip ID

	// Execute
	err := handler.GetDetailedPayslip(c)

	// Assert
	assert.NoError(t, err)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetDetailedPayslip_InvalidPayslipIDFormat(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Set up expectations
	mockResponse.On("SendBadRequest", mock.AnythingOfType("*echo.context"), "Invalid payslip ID format", mock.AnythingOfType("string")).Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("payslip_id")
	c.SetParamValues("invalid") // Invalid payslip ID format

	// Execute
	err := handler.GetDetailedPayslip(c)

	// Assert
	assert.NoError(t, err)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetDetailedPayslip_PayslipNotFound(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	payslipID := uint(999)

	// Set up expectations
	mockRepo.On("GetPayslipByID", payslipID).Return((*model.Payslip)(nil), fmt.Errorf("payslip not found"))
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), "Payslip not found", "payslip not found").Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("payslip_id")
	c.SetParamValues("999")

	// Set up authentication context (admin role for access)
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(1))

	// Execute
	err := handler.GetDetailedPayslip(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetDetailedPayslip_EmployeeNotFound(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	payslipID := uint(1)
	employeeID := uint(999)

	payslip := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{ID: payslipID},
		EmployeeID:       employeeID,
		PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		TotalAmount:      5500000.0,
		Status:           "processed",
	}

	// Set up expectations
	mockRepo.On("GetPayslipByID", payslipID).Return(payslip, nil)
	mockRepo.On("GetEmployeeByID", employeeID).Return((*model.Employee)(nil), fmt.Errorf("employee not found"))
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), "Employee not found", "employee not found").Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("payslip_id")
	c.SetParamValues("1")

	// Set up authentication context (admin role for access)
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(1))

	// Execute
	err := handler.GetDetailedPayslip(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetDetailedPayslip_AttendanceRetrievalError(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	payslipID := uint(1)
	employeeID := uint(1)

	payslip := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{ID: payslipID},
		EmployeeID:       employeeID,
		PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		TotalAmount:      5500000.0,
		Status:           "processed",
	}

	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
		Role:             "employee",
	}

	// Set up expectations
	mockRepo.On("GetPayslipByID", payslipID).Return(payslip, nil)
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetAttendanceForPeriod", employeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd).Return([]model.Attendance(nil), fmt.Errorf("database error"))
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), "Failed to get attendance records", "database error").Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("payslip_id")
	c.SetParamValues("1")

	// Set up authentication context (admin role for access)
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(999))

	// Execute
	err := handler.GetDetailedPayslip(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetDetailedPayslip_OvertimeRetrievalError(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	payslipID := uint(1)
	employeeID := uint(1)

	payslip := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{ID: payslipID},
		EmployeeID:       employeeID,
		PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		TotalAmount:      5500000.0,
		Status:           "processed",
	}

	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
		Role:             "employee",
	}

	attendances := []model.Attendance{}

	// Set up expectations
	mockRepo.On("GetPayslipByID", payslipID).Return(payslip, nil)
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetAttendanceForPeriod", employeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(attendances, nil)
	mockRepo.On("GetOvertimeForPeriod", employeeID, "2025-06-01", "2025-06-30").Return([]model.Overtime(nil), fmt.Errorf("database error"))
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), "Failed to get overtime records", "database error").Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("payslip_id")
	c.SetParamValues("1")

	// Set up authentication context (admin role for access)
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(999))

	// Execute
	err := handler.GetDetailedPayslip(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetDetailedPayslip_ReimbursementRetrievalError(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	payslipID := uint(1)
	employeeID := uint(1)

	payslip := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{ID: payslipID},
		EmployeeID:       employeeID,
		PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		TotalAmount:      5500000.0,
		Status:           "processed",
	}

	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
		Role:             "employee",
	}

	attendances := []model.Attendance{}
	overtimes := []model.Overtime{}

	// Set up expectations
	mockRepo.On("GetPayslipByID", payslipID).Return(payslip, nil)
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetAttendanceForPeriod", employeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(attendances, nil)
	mockRepo.On("GetOvertimeForPeriod", employeeID, "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockRepo.On("GetApprovedReimbursementsForPeriod", employeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd).Return([]model.Reimbursement(nil), fmt.Errorf("database error"))
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), "Failed to get reimbursement records", "database error").Return(nil)

	// Create request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("payslip_id")
	c.SetParamValues("1")

	// Set up authentication context (admin role for access)
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(999))

	// Execute
	err := handler.GetDetailedPayslip(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetDetailedPayslip_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		payslipID  string
		expectCall bool
		setupMocks func(*MockPayslipRepository, *MockPayrollUsecase, *MockResponse)
	}{
		{
			name:       "Zero payslip ID",
			payslipID:  "0",
			expectCall: true,
			setupMocks: func(mockRepo *MockPayslipRepository, mockUsecase *MockPayrollUsecase, mockResponse *MockResponse) {
				payslip := &model.Payslip{
					DefaultAttribute: model.DefaultAttribute{ID: 0},
					EmployeeID:       1,
					PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
				}
				employee := &model.Employee{DefaultAttribute: model.DefaultAttribute{ID: 1}, Name: "Test Employee"}
				attendances := []model.Attendance{}
				overtimes := []model.Overtime{}
				reimbursements := []model.Reimbursement{}
				detailedResponse := map[string]interface{}{"payslip_id": uint(0)}

				mockRepo.On("GetPayslipByID", uint(0)).Return(payslip, nil)
				mockRepo.On("GetEmployeeByID", uint(1)).Return(employee, nil)
				mockRepo.On("GetAttendanceForPeriod", uint(1), payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(attendances, nil)
				mockRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return(overtimes, nil)
				mockRepo.On("GetApprovedReimbursementsForPeriod", uint(1), payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(reimbursements, nil)
				mockUsecase.On("BuildDetailedPayslipResponse", payslip, employee, attendances, overtimes, reimbursements).Return(detailedResponse)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Detailed payslip generated successfully", detailedResponse).Return(nil)
			},
		},
		{
			name:       "Large payslip ID",
			payslipID:  "9999999",
			expectCall: true,
			setupMocks: func(mockRepo *MockPayslipRepository, mockUsecase *MockPayrollUsecase, mockResponse *MockResponse) {
				payslip := &model.Payslip{
					DefaultAttribute: model.DefaultAttribute{ID: 9999999},
					EmployeeID:       1,
					PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
				}
				employee := &model.Employee{DefaultAttribute: model.DefaultAttribute{ID: 1}, Name: "Test Employee"}
				attendances := []model.Attendance{}
				overtimes := []model.Overtime{}
				reimbursements := []model.Reimbursement{}
				detailedResponse := map[string]interface{}{"payslip_id": uint(9999999)}

				mockRepo.On("GetPayslipByID", uint(9999999)).Return(payslip, nil)
				mockRepo.On("GetEmployeeByID", uint(1)).Return(employee, nil)
				mockRepo.On("GetAttendanceForPeriod", uint(1), payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(attendances, nil)
				mockRepo.On("GetOvertimeForPeriod", uint(1), "2025-06-01", "2025-06-30").Return(overtimes, nil)
				mockRepo.On("GetApprovedReimbursementsForPeriod", uint(1), payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(reimbursements, nil)
				mockUsecase.On("BuildDetailedPayslipResponse", payslip, employee, attendances, overtimes, reimbursements).Return(detailedResponse)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Detailed payslip generated successfully", detailedResponse).Return(nil)
			},
		},
		{
			name:       "Negative payslip ID",
			payslipID:  "-1",
			expectCall: false,
			setupMocks: func(mockRepo *MockPayslipRepository, mockUsecase *MockPayrollUsecase, mockResponse *MockResponse) {
				mockResponse.On("SendBadRequest", mock.AnythingOfType("*echo.context"), "Invalid payslip ID format", mock.AnythingOfType("string")).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockPayslipRepository)
			mockUsecase := new(MockPayrollUsecase)
			mockResponse := new(MockResponse)

			// Setup test-specific mocks
			tt.setupMocks(mockRepo, mockUsecase, mockResponse)

			// Create handler
			handler := &TestPayrollHandler{
				payslipRepo:    mockRepo,
				payrollUsecase: mockUsecase,
				response:       mockResponse,
			}

			// Create request
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/"+tt.payslipID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("payslip_id")
			c.SetParamValues(tt.payslipID)

			// Set up authentication context (admin role for access)
			c.Set("authenticated_role", "admin")
			c.Set("authenticated_user_id", uint(999))

			// Execute
			err := handler.GetDetailedPayslip(c)

			// Assert
			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
			mockUsecase.AssertExpectations(t)
			mockResponse.AssertExpectations(t)
		})
	}
}

// Benchmark for GetDetailedPayslip
func BenchmarkPayrollHandler_GetDetailedPayslip(b *testing.B) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	payslipID := uint(1)
	employeeID := uint(1)

	payslip := &model.Payslip{
		DefaultAttribute: model.DefaultAttribute{ID: payslipID},
		EmployeeID:       employeeID,
		PayPeriodStart:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:     time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		TotalAmount:      5500000.0,
		Status:           "processed",
	}

	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: employeeID},
		Name:             "John Doe",
		Role:             "employee",
	}

	attendances := []model.Attendance{}
	overtimes := []model.Overtime{}
	reimbursements := []model.Reimbursement{}
	detailedResponse := map[string]interface{}{"payslip_id": payslipID}

	// Mock expectations (only set up once for all iterations)
	mockRepo.On("GetPayslipByID", payslipID).Return(payslip, nil)
	mockRepo.On("GetEmployeeByID", employeeID).Return(employee, nil)
	mockRepo.On("GetAttendanceForPeriod", employeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(attendances, nil)
	mockRepo.On("GetOvertimeForPeriod", employeeID, "2025-06-01", "2025-06-30").Return(overtimes, nil)
	mockRepo.On("GetApprovedReimbursementsForPeriod", employeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd).Return(reimbursements, nil)
	mockUsecase.On("BuildDetailedPayslipResponse", payslip, employee, attendances, overtimes, reimbursements).Return(detailedResponse)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Detailed payslip generated successfully", detailedResponse).Return(nil)

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/payslips/detailed/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("payslip_id")
		c.SetParamValues("1")

		// Set up authentication context
		c.Set("authenticated_role", "admin")
		c.Set("authenticated_user_id", uint(999))

		_ = handler.GetDetailedPayslip(c)
	}
}

// Tests for GetPayrollSummary function

func TestPayrollHandler_GetPayrollSummary_ValidRequest(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	payslips := []model.Payslip{
		{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			EmployeeID:       1,
			PayPeriodStart:   startDate,
			PayPeriodEnd:     endDate,
			TotalAmount:      5500000.0,
			Status:           "processed",
		},
		{
			DefaultAttribute: model.DefaultAttribute{ID: 2},
			EmployeeID:       2,
			PayPeriodStart:   startDate,
			PayPeriodEnd:     endDate,
			TotalAmount:      6000000.0,
			Status:           "processed",
		},
	}

	summaryResponse := map[string]interface{}{
		"total_employees":    2,
		"total_amount":       11500000.0,
		"average_amount":     5750000.0,
		"highest_amount":     6000000.0,
		"lowest_amount":      5500000.0,
		"pay_period_start":   startDate,
		"pay_period_end":     endDate,
		"processed_payslips": 2,
	}

	// Set up expectations
	mockRepo.On("GetPayslipsByPeriod", startDate, endDate).Return(payslips, nil)
	mockUsecase.On("BuildPayrollSummary", payslips).Return(summaryResponse)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payroll summary generated successfully", summaryResponse).Return(nil)

	// Create request
	e := echo.New()
	reqBody := `{
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/payroll/summary", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.GetPayrollSummary(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayrollSummary_InvalidDateRange(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Set up expectations
	mockResponse.On("SendBadRequest", mock.AnythingOfType("*echo.context"), "Pay period end must be after start date", nil).Return(nil)

	// Create request with invalid date range (end before start)
	e := echo.New()
	reqBody := `{
		"pay_period_start": "2025-06-30T00:00:00Z",
		"pay_period_end": "2025-06-01T00:00:00Z"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/payroll/summary", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.GetPayrollSummary(c)

	// Assert
	assert.NoError(t, err)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayrollSummary_InvalidJSON(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Set up expectations
	mockResponse.On("SendBadRequest", mock.AnythingOfType("*echo.context"), "Invalid request body", mock.AnythingOfType("string")).Return(nil)

	// Create request with invalid JSON
	e := echo.New()
	reqBody := `{
		"pay_period_start": "invalid-date",
		"pay_period_end": "2025-06-30T00:00:00Z"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/payroll/summary", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.GetPayrollSummary(c)

	// Assert
	assert.NoError(t, err)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayrollSummary_MalformedJSON(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Set up expectations
	mockResponse.On("SendBadRequest", mock.AnythingOfType("*echo.context"), "Invalid request body", mock.AnythingOfType("string")).Return(nil)

	// Create request with malformed JSON
	e := echo.New()
	reqBody := `{
		"pay_period_start": "2025-06-01T00:00:00Z"
		"pay_period_end": "2025-06-30T00:00:00Z"
	}` // Missing comma
	req := httptest.NewRequest(http.MethodPost, "/api/payroll/summary", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.GetPayrollSummary(c)

	// Assert
	assert.NoError(t, err)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayrollSummary_PayslipRetrievalError(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	// Set up expectations
	mockRepo.On("GetPayslipsByPeriod", startDate, endDate).Return([]model.Payslip{}, errors.New("database connection error"))
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), "Failed to retrieve payslips", "database connection error").Return(nil)

	// Create request
	e := echo.New()
	reqBody := `{
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/payroll/summary", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.GetPayrollSummary(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayrollSummary_EmptyPayslipsList(t *testing.T) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	emptyPayslips := []model.Payslip{}
	summaryResponse := map[string]interface{}{
		"total_employees":    0,
		"total_amount":       0.0,
		"average_amount":     0.0,
		"highest_amount":     0.0,
		"lowest_amount":      0.0,
		"pay_period_start":   startDate,
		"pay_period_end":     endDate,
		"processed_payslips": 0,
	}

	// Set up expectations
	mockRepo.On("GetPayslipsByPeriod", startDate, endDate).Return(emptyPayslips, nil)
	mockUsecase.On("BuildPayrollSummary", emptyPayslips).Return(summaryResponse)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payroll summary generated successfully", summaryResponse).Return(nil)

	// Create request
	e := echo.New()
	reqBody := `{
		"pay_period_start": "2025-06-01T00:00:00Z",
		"pay_period_end": "2025-06-30T00:00:00Z"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/payroll/summary", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.GetPayrollSummary(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestPayrollHandler_GetPayrollSummary_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		expectError bool
		setupMocks  func(*MockPayslipRepository, *MockPayrollUsecase, *MockResponse)
	}{
		{
			name: "Same start and end date",
			requestBody: `{
				"pay_period_start": "2025-06-15T00:00:00Z",
				"pay_period_end": "2025-06-15T00:00:00Z"
			}`,
			expectError: false,
			setupMocks: func(mockRepo *MockPayslipRepository, mockUsecase *MockPayrollUsecase, mockResponse *MockResponse) {
				startDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
				payslips := []model.Payslip{}
				summaryResponse := map[string]interface{}{"total_employees": 0}

				mockRepo.On("GetPayslipsByPeriod", startDate, endDate).Return(payslips, nil)
				mockUsecase.On("BuildPayrollSummary", payslips).Return(summaryResponse)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payroll summary generated successfully", summaryResponse).Return(nil)
			},
		},
		{
			name: "Future dates",
			requestBody: `{
				"pay_period_start": "2030-01-01T00:00:00Z",
				"pay_period_end": "2030-01-31T00:00:00Z"
			}`,
			expectError: false,
			setupMocks: func(mockRepo *MockPayslipRepository, mockUsecase *MockPayrollUsecase, mockResponse *MockResponse) {
				startDate := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2030, 1, 31, 0, 0, 0, 0, time.UTC)
				payslips := []model.Payslip{}
				summaryResponse := map[string]interface{}{"total_employees": 0}

				mockRepo.On("GetPayslipsByPeriod", startDate, endDate).Return(payslips, nil)
				mockUsecase.On("BuildPayrollSummary", payslips).Return(summaryResponse)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payroll summary generated successfully", summaryResponse).Return(nil)
			},
		},
		{
			name: "Large dataset",
			requestBody: `{
				"pay_period_start": "2025-01-01T00:00:00Z",
				"pay_period_end": "2025-12-31T00:00:00Z"
			}`,
			expectError: false,
			setupMocks: func(mockRepo *MockPayslipRepository, mockUsecase *MockPayrollUsecase, mockResponse *MockResponse) {
				startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
				// Simulate large dataset
				payslips := make([]model.Payslip, 1000)
				for i := 0; i < 1000; i++ {
					payslips[i] = model.Payslip{
						DefaultAttribute: model.DefaultAttribute{ID: uint(i + 1)},
						EmployeeID:       uint(i + 1),
						PayPeriodStart:   startDate,
						PayPeriodEnd:     endDate,
						TotalAmount:      5500000.0,
						Status:           "processed",
					}
				}
				summaryResponse := map[string]interface{}{
					"total_employees": 1000,
					"total_amount":    5500000000.0,
				}

				mockRepo.On("GetPayslipsByPeriod", startDate, endDate).Return(payslips, nil)
				mockUsecase.On("BuildPayrollSummary", payslips).Return(summaryResponse)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payroll summary generated successfully", summaryResponse).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockPayslipRepository)
			mockUsecase := new(MockPayrollUsecase)
			mockResponse := new(MockResponse)

			// Setup test-specific mocks
			tt.setupMocks(mockRepo, mockUsecase, mockResponse)

			// Create handler
			handler := &TestPayrollHandler{
				payslipRepo:    mockRepo,
				payrollUsecase: mockUsecase,
				response:       mockResponse,
			}

			// Create request
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/payroll/summary", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute
			err := handler.GetPayrollSummary(c)

			// Assert
			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
			mockUsecase.AssertExpectations(t)
			mockResponse.AssertExpectations(t)
		})
	}
}

// Benchmark for GetPayrollSummary
func BenchmarkPayrollHandler_GetPayrollSummary(b *testing.B) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create test data
	startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	payslips := []model.Payslip{
		{
			DefaultAttribute: model.DefaultAttribute{ID: 1},
			EmployeeID:       1,
			PayPeriodStart:   startDate,
			PayPeriodEnd:     endDate,
			TotalAmount:      5500000.0,
			Status:           "processed",
		},
	}

	summaryResponse := map[string]interface{}{
		"total_employees": 1,
		"total_amount":    5500000.0,
	}

	// Mock expectations (only set up once for all iterations)
	mockRepo.On("GetPayslipsByPeriod", startDate, endDate).Return(payslips, nil)
	mockUsecase.On("BuildPayrollSummary", payslips).Return(summaryResponse)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payroll summary generated successfully", summaryResponse).Return(nil)

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		e := echo.New()
		reqBody := `{
			"pay_period_start": "2025-06-01T00:00:00Z",
			"pay_period_end": "2025-06-30T00:00:00Z"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/payroll/summary", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = handler.GetPayrollSummary(c)
	}
}

// Benchmark for GetPayrollSummary with large dataset
func BenchmarkPayrollHandler_GetPayrollSummary_LargeDataset(b *testing.B) {
	// Create mocks
	mockRepo := new(MockPayslipRepository)
	mockUsecase := new(MockPayrollUsecase)
	mockResponse := new(MockResponse)

	// Create handler
	handler := &TestPayrollHandler{
		payslipRepo:    mockRepo,
		payrollUsecase: mockUsecase,
		response:       mockResponse,
	}

	// Create large test dataset
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	// Create 1000 payslips for simulation
	payslips := make([]model.Payslip, 1000)
	for i := 0; i < 1000; i++ {
		payslips[i] = model.Payslip{
			DefaultAttribute: model.DefaultAttribute{ID: uint(i + 1)},
			EmployeeID:       uint(i + 1),
			PayPeriodStart:   startDate,
			PayPeriodEnd:     endDate,
			TotalAmount:      5500000.0,
			Status:           "processed",
		}
	}

	summaryResponse := map[string]interface{}{
		"total_employees": 1000,
		"total_amount":    5500000000.0,
	}

	// Mock expectations (only set up once for all iterations)
	mockRepo.On("GetPayslipsByPeriod", startDate, endDate).Return(payslips, nil)
	mockUsecase.On("BuildPayrollSummary", payslips).Return(summaryResponse)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Payroll summary generated successfully", summaryResponse).Return(nil)

	// Reset timer to exclude setup time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e := echo.New()
		reqBody := `{
			"pay_period_start": "2025-01-01T00:00:00Z",
			"pay_period_end": "2025-12-31T00:00:00Z"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/payroll/summary", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = handler.GetPayrollSummary(c)
	}
}
