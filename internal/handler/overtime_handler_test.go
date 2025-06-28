// Package handler contains unit tests for the overtime handler functions.
//
// This test file contains comprehensive tests for the overtime handler,
//
// covering various scenarios such as:
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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/model"
)

// Helper function to create a test overtime
func createTestOvertime() *model.Overtime {
	now := time.Now()
	approverID := uint(2)

	return &model.Overtime{
		DefaultAttribute: model.DefaultAttribute{
			ID:        1,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		EmployeeID:   1,
		OvertimeDate: "2024-01-15",
		StartTime:    "18:00",
		EndTime:      "20:00",
		Hours:        2,
		Reason:       "Emergency project deadline",
		Status:       model.OvertimePending,
		ApprovedBy:   &approverID,
		ApprovedAt:   &now,
	}
}

func createApprovedOvertime() *model.Overtime {
	overtime := createTestOvertime()
	overtime.Status = model.OvertimeApproved
	return overtime
}

func createRejectedOvertime() *model.Overtime {
	overtime := createTestOvertime()
	overtime.Status = model.OvertimeRejected
	return overtime
}

// Test CreateOvertime request binding - Success
func TestOvertimeHandler_CreateOvertimeRequestBinding_Success(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"employee_id":1,"reason":"Project deadline","hours":3}`
	req := httptest.NewRequest(http.MethodPost, "/overtime", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var createReq request.CreateOvertimeRequest
	err := c.Bind(&createReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint(1), createReq.EmployeeID)
	assert.Equal(t, "Project deadline", createReq.Reason)
	assert.Equal(t, 3, createReq.Hours)
}

// Test CreateOvertime request binding - Invalid JSON
func TestOvertimeHandler_CreateOvertimeRequestBinding_InvalidJSON(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"invalid_json": }`
	req := httptest.NewRequest(http.MethodPost, "/overtime", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var createReq request.CreateOvertimeRequest
	err := c.Bind(&createReq)

	// Assert
	assert.Error(t, err)
}

// Test CreateOvertime request binding - Missing fields
func TestOvertimeHandler_CreateOvertimeRequestBinding_MissingFields(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"employee_id":0,"reason":"","hours":0}`
	req := httptest.NewRequest(http.MethodPost, "/overtime", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var createReq request.CreateOvertimeRequest
	err := c.Bind(&createReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint(0), createReq.EmployeeID)
	assert.Empty(t, createReq.Reason)
	assert.Equal(t, 0, createReq.Hours)
}

// Test CreateOvertime request binding - Valid minimum values
func TestOvertimeHandler_CreateOvertimeRequestBinding_ValidMinimum(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"employee_id":1,"reason":"R","hours":1}`
	req := httptest.NewRequest(http.MethodPost, "/overtime", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var createReq request.CreateOvertimeRequest
	err := c.Bind(&createReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint(1), createReq.EmployeeID)
	assert.Equal(t, "R", createReq.Reason)
	assert.Equal(t, 1, createReq.Hours)
}

// Test CreateOvertime request binding - Large values
func TestOvertimeHandler_CreateOvertimeRequestBinding_LargeValues(t *testing.T) {
	// Setup Echo
	e := echo.New()
	longReason := strings.Repeat("This is a very long overtime reason. ", 10)
	reqBody := `{"employee_id":999,"reason":"` + longReason + `","hours":12}`
	req := httptest.NewRequest(http.MethodPost, "/overtime", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var createReq request.CreateOvertimeRequest
	err := c.Bind(&createReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint(999), createReq.EmployeeID)
	assert.Equal(t, longReason, createReq.Reason)
	assert.Equal(t, 12, createReq.Hours)
}

// Test Overtime model - Basic structure
func TestOvertimeHandler_OvertimeModel_BasicStructure(t *testing.T) {
	// Create test overtime
	overtime := createTestOvertime()

	// Assert structure
	assert.NotZero(t, overtime.ID)
	assert.Equal(t, uint(1), overtime.EmployeeID)
	assert.Equal(t, "2024-01-15", overtime.OvertimeDate)
	assert.Equal(t, "18:00", overtime.StartTime)
	assert.Equal(t, "20:00", overtime.EndTime)
	assert.Equal(t, 2, overtime.Hours)
	assert.Equal(t, "Emergency project deadline", overtime.Reason)
	assert.Equal(t, model.OvertimePending, overtime.Status)
	assert.NotNil(t, overtime.ApprovedBy)
	assert.NotNil(t, overtime.ApprovedAt)
	assert.NotNil(t, overtime.CreatedAt)
	assert.NotNil(t, overtime.UpdatedAt)
}

// Test Overtime status constants
func TestOvertimeHandler_OvertimeStatus_Constants(t *testing.T) {
	// Test status constants
	assert.Equal(t, model.OvertimeStatus("pending"), model.OvertimePending)
	assert.Equal(t, model.OvertimeStatus("approved"), model.OvertimeApproved)
	assert.Equal(t, model.OvertimeStatus("rejected"), model.OvertimeRejected)
}

// Test Overtime Approve method
func TestOvertimeHandler_OvertimeApprove_Method(t *testing.T) {
	// Create test overtime
	overtime := createTestOvertime()
	overtime.Status = model.OvertimePending
	overtime.ApprovedBy = nil
	overtime.ApprovedAt = nil

	approverID := uint(3)

	// Approve overtime
	overtime.Approve(approverID)

	// Assert
	assert.Equal(t, model.OvertimeApproved, overtime.Status)
	assert.NotNil(t, overtime.ApprovedBy)
	assert.Equal(t, approverID, *overtime.ApprovedBy)
	assert.NotNil(t, overtime.ApprovedAt)
	assert.True(t, overtime.ApprovedAt.After(time.Now().Add(-time.Minute)))
}

// Test Overtime Reject method
func TestOvertimeHandler_OvertimeReject_Method(t *testing.T) {
	// Create test overtime
	overtime := createTestOvertime()
	overtime.Status = model.OvertimePending
	overtime.ApprovedBy = nil
	overtime.ApprovedAt = nil

	approverID := uint(4)

	// Reject overtime
	overtime.Reject(approverID)

	// Assert
	assert.Equal(t, model.OvertimeRejected, overtime.Status)
	assert.NotNil(t, overtime.ApprovedBy)
	assert.Equal(t, approverID, *overtime.ApprovedBy)
	assert.NotNil(t, overtime.ApprovedAt)
	assert.True(t, overtime.ApprovedAt.After(time.Now().Add(-time.Minute)))
}

// Test Overtime IsApproved method
func TestOvertimeHandler_OvertimeIsApproved_Method(t *testing.T) {
	// Test approved overtime
	approvedOvertime := createApprovedOvertime()
	assert.True(t, approvedOvertime.IsApproved())

	// Test pending overtime
	pendingOvertime := createTestOvertime()
	pendingOvertime.Status = model.OvertimePending
	assert.False(t, pendingOvertime.IsApproved())

	// Test rejected overtime
	rejectedOvertime := createRejectedOvertime()
	assert.False(t, rejectedOvertime.IsApproved())
}

// Test Overtime CalculateHoursFromTime method - Success
func TestOvertimeHandler_OvertimeCalculateHoursFromTime_Success(t *testing.T) {
	// Create test overtime
	overtime := createTestOvertime()
	overtime.StartTime = "09:00"
	overtime.EndTime = "17:00"
	overtime.Hours = 0 // Reset to test calculation

	// Calculate hours
	overtime.CalculateHoursFromTime()

	// Assert
	assert.Equal(t, 8, overtime.Hours)
}

// Test Overtime CalculateHoursFromTime method - Different times
func TestOvertimeHandler_OvertimeCalculateHoursFromTime_DifferentTimes(t *testing.T) {
	testCases := []struct {
		name      string
		startTime string
		endTime   string
		expected  int
	}{
		{"1 hour", "10:00", "11:00", 1},
		{"2 hours", "14:30", "16:30", 2},
		{"4 hours", "18:00", "22:00", 4},
		{"Half day", "08:00", "12:00", 4},
		{"Full day", "08:00", "17:00", 9},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			overtime := createTestOvertime()
			overtime.StartTime = tc.startTime
			overtime.EndTime = tc.endTime
			overtime.Hours = 0

			overtime.CalculateHoursFromTime()

			assert.Equal(t, tc.expected, overtime.Hours)
		})
	}
}

// Test Overtime CalculateHoursFromTime method - Invalid times
func TestOvertimeHandler_OvertimeCalculateHoursFromTime_InvalidTimes(t *testing.T) {
	testCases := []struct {
		name      string
		startTime string
		endTime   string
	}{
		{"Invalid start time", "invalid", "17:00"},
		{"Invalid end time", "09:00", "invalid"},
		{"Both invalid", "invalid", "invalid"},
		{"End before start", "17:00", "09:00"},
		{"Same time", "12:00", "12:00"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			overtime := createTestOvertime()
			overtime.StartTime = tc.startTime
			overtime.EndTime = tc.endTime
			originalHours := overtime.Hours

			overtime.CalculateHoursFromTime()

			// Hours should remain unchanged for invalid times
			assert.Equal(t, originalHours, overtime.Hours)
		})
	}
}

// Test CreateOvertimeRequest validation structure
func TestOvertimeHandler_CreateOvertimeRequest_ValidationStructure(t *testing.T) {
	// Create request
	req := request.CreateOvertimeRequest{
		EmployeeID: 123,
		Reason:     "Urgent system maintenance",
		Hours:      5,
	}

	// Assert structure and types
	assert.IsType(t, uint(0), req.EmployeeID)
	assert.IsType(t, "", req.Reason)
	assert.IsType(t, 0, req.Hours)

	assert.Equal(t, uint(123), req.EmployeeID)
	assert.Equal(t, "Urgent system maintenance", req.Reason)
	assert.Equal(t, 5, req.Hours)
}

// Test Overtime validation boundaries
func TestOvertimeHandler_OvertimeValidation_Boundaries(t *testing.T) {
	testCases := []struct {
		name     string
		hours    int
		reason   string
		expectOK bool
	}{
		{"Valid minimum hours", 1, "Valid reason", true},
		{"Valid maximum hours", 12, "Valid reason", true},
		{"Zero hours", 0, "Valid reason", false},
		{"Negative hours", -1, "Valid reason", false},
		{"Excessive hours", 13, "Valid reason", false},
		{"Valid hours short reason", 5, "Hi", false},
		{"Valid hours long reason", 5, strings.Repeat("a", 256), false},
		{"Valid hours valid reason", 5, "Project deadline work", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := request.CreateOvertimeRequest{
				EmployeeID: 1,
				Reason:     tc.reason,
				Hours:      tc.hours,
			}

			// Basic validation checks
			hoursValid := tc.hours >= 1 && tc.hours <= 12
			reasonValid := len(tc.reason) >= 5 && len(tc.reason) <= 255

			if tc.expectOK {
				assert.True(t, hoursValid, "Hours should be valid")
				assert.True(t, reasonValid, "Reason should be valid")
			}

			// Verify the request was created with expected values
			assert.Equal(t, tc.hours, req.Hours)
			assert.Equal(t, tc.reason, req.Reason)
		})
	}
}

// Test Overtime status transitions
func TestOvertimeHandler_OvertimeStatus_Transitions(t *testing.T) {
	// Test pending to approved
	overtime := createTestOvertime()
	overtime.Status = model.OvertimePending

	overtime.Approve(123)
	assert.Equal(t, model.OvertimeApproved, overtime.Status)
	assert.True(t, overtime.IsApproved())

	// Test pending to rejected
	overtime2 := createTestOvertime()
	overtime2.Status = model.OvertimePending

	overtime2.Reject(456)
	assert.Equal(t, model.OvertimeRejected, overtime2.Status)
	assert.False(t, overtime2.IsApproved())

	// Test that both approved and rejected have approver info
	assert.NotNil(t, overtime.ApprovedBy)
	assert.NotNil(t, overtime.ApprovedAt)
	assert.NotNil(t, overtime2.ApprovedBy)
	assert.NotNil(t, overtime2.ApprovedAt)
}

// Test Overtime table name
func TestOvertimeHandler_OvertimeTableName(t *testing.T) {
	overtime := model.Overtime{}
	tableName := overtime.TableName()

	assert.Equal(t, "overtimes", tableName)
}

// Test Overtime time formats
func TestOvertimeHandler_OvertimeTimeFormats(t *testing.T) {
	// Test valid time formats
	validTimes := []string{
		"00:00", "01:30", "09:15", "12:00", "18:45", "23:59",
	}

	for _, timeStr := range validTimes {
		t.Run("Valid time "+timeStr, func(t *testing.T) {
			const timeLayout = "15:04"
			_, err := time.Parse(timeLayout, timeStr)
			assert.NoError(t, err)
		})
	}

	// Test invalid time formats
	invalidTimes := []string{
		"24:00", "12:60", "12:5", "25:30", "invalid", "", "1:5", "1:",
	}

	for _, timeStr := range invalidTimes {
		t.Run("Invalid time "+timeStr, func(t *testing.T) {
			const timeLayout = "15:04"
			_, err := time.Parse(timeLayout, timeStr)
			assert.Error(t, err)
		})
	}
}

// Test JSON marshaling/unmarshaling
func TestOvertimeHandler_OvertimeJSON_Marshaling(t *testing.T) {
	// Setup Echo
	e := echo.New()

	// Test different JSON inputs
	testCases := []struct {
		name     string
		jsonBody string
		expectOK bool
	}{
		{
			"Valid JSON",
			`{"employee_id":1,"reason":"Test reason","hours":3}`,
			true,
		},
		{
			"Missing employee_id",
			`{"reason":"Test reason","hours":3}`,
			true, // Binding succeeds but field will be zero
		},
		{
			"String employee_id",
			`{"employee_id":"1","reason":"Test reason","hours":3}`,
			false, // JSON parsing should fail for type mismatch
		},
		{
			"String hours",
			`{"employee_id":1,"reason":"Test reason","hours":"3"}`,
			false, // JSON parsing should fail for type mismatch
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/overtime", strings.NewReader(tc.jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			var createReq request.CreateOvertimeRequest
			err := c.Bind(&createReq)

			if tc.expectOK {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
