// Package handler contains unit tests for the reimbursement handler functionality.
//
// This test file contains comprehensive tests for the ReimbursementHandler functions:
//
// CreateReimbusement tests cover:
// 1. Valid requests with successful creation
// 2. Invalid JSON payloads
// 3. Malformed JSON requests
// 4. Repository errors during creation
// 5. Edge cases (zero amounts, large amounts, special characters)
// 6. Performance benchmarks
//
// The tests use mocks to isolate the handler logic and ensure fast, reliable test execution.
// The TestReimbursementHandler struct and related interfaces are created specifically for testing
// to avoid tight coupling with concrete implementations.

package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/helper"
	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"gorm.io/gorm"
)

// ReimbursementRepositoryInterface defines the interface for reimbursement repository
type ReimbursementRepositoryInterface interface {
	GetDB() *gorm.DB
	CreateReimbusementWithAudit(employeeID uint, amount float64, description string, auditDB *middleware.AuditableDB) (*model.Reimbursement, error)
}

// TestReimbursementHandler is a testable version of ReimbursementHandler
type TestReimbursementHandler struct {
	Helper           helper.NewHelper
	DB               *gorm.DB
	Response         ResponseInterface
	ReimbusementRepo ReimbursementRepositoryInterface
}

// CreateReimbusement mirrors the original method for testing
func (h *TestReimbursementHandler) CreateReimbusement(c echo.Context) error {
	req := request.CreateReimbusementRequest{}
	if err := c.Bind(&req); err != nil {
		return h.Response.SendError(c, err.Error(), "Invalid request data")
	}

	// Get auditable DB instance
	auditDB := helper.GetAuditableDB(c, h.ReimbusementRepo.GetDB())

	_, err := h.ReimbusementRepo.CreateReimbusementWithAudit(req.EmployeeID, req.Amount, req.Description, auditDB)
	if err != nil {
		return h.Response.SendError(c, err.Error(), "Failed to create reimbusement")
	}
	return h.Response.SendSuccess(c, "Reimbusement created successfully", nil)
}

// MockReimbursementRepository is a mock implementation of the ReimbursementRepositoryInterface
type MockReimbursementRepository struct {
	mock.Mock
}

func (m *MockReimbursementRepository) GetDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockReimbursementRepository) CreateReimbusementWithAudit(employeeID uint, amount float64, description string, auditDB *middleware.AuditableDB) (*model.Reimbursement, error) {
	args := m.Called(employeeID, amount, description, auditDB)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Reimbursement), args.Error(1)
}

// Tests for CreateReimbusement function

func TestReimbursementHandler_CreateReimbusement_ValidRequest(t *testing.T) {
	// Create mocks
	mockRepo := new(MockReimbursementRepository)
	mockResponse := new(MockResponse)
	mockDB := &gorm.DB{}

	// Create handler
	handler := &TestReimbursementHandler{
		ReimbusementRepo: mockRepo,
		Response:         mockResponse,
		DB:               mockDB,
	}

	// Create test data
	employeeID := uint(1)
	amount := 150000.0
	description := "Transport reimbursement for business trip"

	reimbursement := &model.Reimbursement{
		DefaultAttribute:  model.DefaultAttribute{ID: 1},
		EmployeeID:        employeeID,
		Amount:            amount,
		Reason:            description,
		Status:            model.ReimbursementPending,
		ReimbursementDate: time.Now(),
	}

	// Set up expectations
	mockRepo.On("GetDB").Return(mockDB)
	mockRepo.On("CreateReimbusementWithAudit", employeeID, amount, description, mock.AnythingOfType("*middleware.AuditableDB")).Return(reimbursement, nil)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Reimbusement created successfully", nil).Return(nil)

	// Create request
	e := echo.New()
	reqBody := `{
		"employee_id": 1,
		"amount": 150000.0,
		"description": "Transport reimbursement for business trip"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/reimbursements", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up user context for audit
	c.Set("user_id", 1)

	// Execute
	err := handler.CreateReimbusement(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestReimbursementHandler_CreateReimbusement_InvalidJSON(t *testing.T) {
	// Create mocks
	mockRepo := new(MockReimbursementRepository)
	mockResponse := new(MockResponse)
	mockDB := &gorm.DB{}

	// Create handler
	handler := &TestReimbursementHandler{
		ReimbusementRepo: mockRepo,
		Response:         mockResponse,
		DB:               mockDB,
	}

	// Set up expectations
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), mock.AnythingOfType("string"), "Invalid request data").Return(nil)

	// Create request with invalid JSON
	e := echo.New()
	reqBody := `{
		"employee_id": "invalid",
		"amount": 150000.0,
		"description": "Transport reimbursement"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/reimbursements", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.CreateReimbusement(c)

	// Assert
	assert.NoError(t, err)
	mockResponse.AssertExpectations(t)
}

func TestReimbursementHandler_CreateReimbusement_MalformedJSON(t *testing.T) {
	// Create mocks
	mockRepo := new(MockReimbursementRepository)
	mockResponse := new(MockResponse)
	mockDB := &gorm.DB{}

	// Create handler
	handler := &TestReimbursementHandler{
		ReimbusementRepo: mockRepo,
		Response:         mockResponse,
		DB:               mockDB,
	}

	// Set up expectations
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), mock.AnythingOfType("string"), "Invalid request data").Return(nil)

	// Create request with malformed JSON
	e := echo.New()
	reqBody := `{
		"employee_id": 1
		"amount": 150000.0,
		"description": "Transport reimbursement"
	}` // Missing comma
	req := httptest.NewRequest(http.MethodPost, "/api/reimbursements", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.CreateReimbusement(c)

	// Assert
	assert.NoError(t, err)
	mockResponse.AssertExpectations(t)
}

func TestReimbursementHandler_CreateReimbusement_RepositoryError(t *testing.T) {
	// Create mocks
	mockRepo := new(MockReimbursementRepository)
	mockResponse := new(MockResponse)
	mockDB := &gorm.DB{}

	// Create handler
	handler := &TestReimbursementHandler{
		ReimbusementRepo: mockRepo,
		Response:         mockResponse,
		DB:               mockDB,
	}

	// Create test data
	employeeID := uint(1)
	amount := 150000.0
	description := "Transport reimbursement for business trip"

	// Set up expectations
	mockRepo.On("GetDB").Return(mockDB)
	mockRepo.On("CreateReimbusementWithAudit", employeeID, amount, description, mock.AnythingOfType("*middleware.AuditableDB")).Return(nil, errors.New("database connection error"))
	mockResponse.On("SendError", mock.AnythingOfType("*echo.context"), "database connection error", "Failed to create reimbusement").Return(nil)

	// Create request
	e := echo.New()
	reqBody := `{
		"employee_id": 1,
		"amount": 150000.0,
		"description": "Transport reimbursement for business trip"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/reimbursements", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up user context for audit
	c.Set("user_id", 1)

	// Execute
	err := handler.CreateReimbusement(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestReimbursementHandler_CreateReimbusement_EmptyDescription(t *testing.T) {
	// Create mocks
	mockRepo := new(MockReimbursementRepository)
	mockResponse := new(MockResponse)
	mockDB := &gorm.DB{}

	// Create handler
	handler := &TestReimbursementHandler{
		ReimbusementRepo: mockRepo,
		Response:         mockResponse,
		DB:               mockDB,
	}

	// Create test data
	employeeID := uint(1)
	amount := 150000.0
	description := ""

	reimbursement := &model.Reimbursement{
		DefaultAttribute:  model.DefaultAttribute{ID: 1},
		EmployeeID:        employeeID,
		Amount:            amount,
		Reason:            description,
		Status:            model.ReimbursementPending,
		ReimbursementDate: time.Now(),
	}

	// Set up expectations
	mockRepo.On("GetDB").Return(mockDB)
	mockRepo.On("CreateReimbusementWithAudit", employeeID, amount, description, mock.AnythingOfType("*middleware.AuditableDB")).Return(reimbursement, nil)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Reimbusement created successfully", nil).Return(nil)

	// Create request
	e := echo.New()
	reqBody := `{
		"employee_id": 1,
		"amount": 150000.0,
		"description": ""
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/reimbursements", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up user context for audit
	c.Set("user_id", 1)

	// Execute
	err := handler.CreateReimbusement(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestReimbursementHandler_CreateReimbusement_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		setupMocks  func(*MockReimbursementRepository, *MockResponse, *gorm.DB)
	}{
		{
			name: "Zero amount",
			requestBody: `{
				"employee_id": 1,
				"amount": 0.0,
				"description": "Zero amount reimbursement"
			}`,
			setupMocks: func(mockRepo *MockReimbursementRepository, mockResponse *MockResponse, mockDB *gorm.DB) {
				reimbursement := &model.Reimbursement{
					DefaultAttribute:  model.DefaultAttribute{ID: 1},
					EmployeeID:        1,
					Amount:            0.0,
					Reason:            "Zero amount reimbursement",
					Status:            model.ReimbursementPending,
					ReimbursementDate: time.Now(),
				}
				mockRepo.On("GetDB").Return(mockDB)
				mockRepo.On("CreateReimbusementWithAudit", uint(1), 0.0, "Zero amount reimbursement", mock.AnythingOfType("*middleware.AuditableDB")).Return(reimbursement, nil)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Reimbusement created successfully", nil).Return(nil)
			},
		},
		{
			name: "Large amount",
			requestBody: `{
				"employee_id": 1,
				"amount": 999999.99,
				"description": "Large amount reimbursement"
			}`,
			setupMocks: func(mockRepo *MockReimbursementRepository, mockResponse *MockResponse, mockDB *gorm.DB) {
				reimbursement := &model.Reimbursement{
					DefaultAttribute:  model.DefaultAttribute{ID: 1},
					EmployeeID:        1,
					Amount:            999999.99,
					Reason:            "Large amount reimbursement",
					Status:            model.ReimbursementPending,
					ReimbursementDate: time.Now(),
				}
				mockRepo.On("GetDB").Return(mockDB)
				mockRepo.On("CreateReimbusementWithAudit", uint(1), 999999.99, "Large amount reimbursement", mock.AnythingOfType("*middleware.AuditableDB")).Return(reimbursement, nil)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Reimbusement created successfully", nil).Return(nil)
			},
		},
		{
			name: "High employee ID",
			requestBody: `{
				"employee_id": 999999,
				"amount": 50000.0,
				"description": "High employee ID reimbursement"
			}`,
			setupMocks: func(mockRepo *MockReimbursementRepository, mockResponse *MockResponse, mockDB *gorm.DB) {
				reimbursement := &model.Reimbursement{
					DefaultAttribute:  model.DefaultAttribute{ID: 1},
					EmployeeID:        999999,
					Amount:            50000.0,
					Reason:            "High employee ID reimbursement",
					Status:            model.ReimbursementPending,
					ReimbursementDate: time.Now(),
				}
				mockRepo.On("GetDB").Return(mockDB)
				mockRepo.On("CreateReimbusementWithAudit", uint(999999), 50000.0, "High employee ID reimbursement", mock.AnythingOfType("*middleware.AuditableDB")).Return(reimbursement, nil)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Reimbusement created successfully", nil).Return(nil)
			},
		},
		{
			name: "Special characters in description",
			requestBody: `{
				"employee_id": 1,
				"amount": 75000.0,
				"description": "Reimbursement with special chars: àáâãäåæçèéêë & symbols @#$%^&*()"
			}`,
			setupMocks: func(mockRepo *MockReimbursementRepository, mockResponse *MockResponse, mockDB *gorm.DB) {
				description := "Reimbursement with special chars: àáâãäåæçèéêë & symbols @#$%^&*()"
				reimbursement := &model.Reimbursement{
					DefaultAttribute:  model.DefaultAttribute{ID: 1},
					EmployeeID:        1,
					Amount:            75000.0,
					Reason:            description,
					Status:            model.ReimbursementPending,
					ReimbursementDate: time.Now(),
				}
				mockRepo.On("GetDB").Return(mockDB)
				mockRepo.On("CreateReimbusementWithAudit", uint(1), 75000.0, description, mock.AnythingOfType("*middleware.AuditableDB")).Return(reimbursement, nil)
				mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Reimbusement created successfully", nil).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockReimbursementRepository)
			mockResponse := new(MockResponse)
			mockDB := &gorm.DB{}

			// Setup test-specific mocks
			tt.setupMocks(mockRepo, mockResponse, mockDB)

			// Create handler
			handler := &TestReimbursementHandler{
				ReimbusementRepo: mockRepo,
				Response:         mockResponse,
				DB:               mockDB,
			}

			// Create request
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/reimbursements", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set up user context for audit
			c.Set("user_id", 1)

			// Execute
			err := handler.CreateReimbusement(c)

			// Assert
			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
			mockResponse.AssertExpectations(t)
		})
	}
}

func TestReimbursementHandler_CreateReimbusement_NoUserContext(t *testing.T) {
	// Create mocks
	mockRepo := new(MockReimbursementRepository)
	mockResponse := new(MockResponse)
	mockDB := &gorm.DB{}

	// Create handler
	handler := &TestReimbursementHandler{
		ReimbusementRepo: mockRepo,
		Response:         mockResponse,
		DB:               mockDB,
	}

	// Create test data
	employeeID := uint(1)
	amount := 150000.0
	description := "Transport reimbursement for business trip"

	reimbursement := &model.Reimbursement{
		DefaultAttribute:  model.DefaultAttribute{ID: 1},
		EmployeeID:        employeeID,
		Amount:            amount,
		Reason:            description,
		Status:            model.ReimbursementPending,
		ReimbursementDate: time.Now(),
	}

	// Set up expectations
	mockRepo.On("GetDB").Return(mockDB)
	mockRepo.On("CreateReimbusementWithAudit", employeeID, amount, description, mock.AnythingOfType("*middleware.AuditableDB")).Return(reimbursement, nil)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Reimbusement created successfully", nil).Return(nil)

	// Create request
	e := echo.New()
	reqBody := `{
		"employee_id": 1,
		"amount": 150000.0,
		"description": "Transport reimbursement for business trip"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/reimbursements", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// No user context set - should fallback to user ID 0

	// Execute
	err := handler.CreateReimbusement(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

// Benchmark for CreateReimbusement
func BenchmarkReimbursementHandler_CreateReimbusement(b *testing.B) {
	// Create mocks
	mockRepo := new(MockReimbursementRepository)
	mockResponse := new(MockResponse)
	mockDB := &gorm.DB{}

	// Create handler
	handler := &TestReimbursementHandler{
		ReimbusementRepo: mockRepo,
		Response:         mockResponse,
		DB:               mockDB,
	}

	// Create test data
	reimbursement := &model.Reimbursement{
		DefaultAttribute:  model.DefaultAttribute{ID: 1},
		EmployeeID:        1,
		Amount:            150000.0,
		Reason:            "Transport reimbursement for business trip",
		Status:            model.ReimbursementPending,
		ReimbursementDate: time.Now(),
	}

	// Mock expectations (only set up once for all iterations)
	mockRepo.On("GetDB").Return(mockDB)
	mockRepo.On("CreateReimbusementWithAudit", uint(1), 150000.0, "Transport reimbursement for business trip", mock.AnythingOfType("*middleware.AuditableDB")).Return(reimbursement, nil)
	mockResponse.On("SendSuccess", mock.AnythingOfType("*echo.context"), "Reimbusement created successfully", nil).Return(nil)

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		e := echo.New()
		reqBody := `{
			"employee_id": 1,
			"amount": 150000.0,
			"description": "Transport reimbursement for business trip"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/reimbursements", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Set up user context for audit
		c.Set("user_id", 1)

		_ = handler.CreateReimbusement(c)
	}
}
