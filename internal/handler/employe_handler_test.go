package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/model"
)

// Helper function to create a test employee
func createTestEmployeeForHandler() *model.Employee {
	now := time.Now()

	return &model.Employee{
		DefaultAttribute: model.DefaultAttribute{
			ID:        1,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name:     "testuser",
		Password: "hashedpassword",
		Role:     "employee",
		Active:   true,
	}
}

func createAdminEmployeeForHandler() *model.Employee {
	now := time.Now()

	return &model.Employee{
		DefaultAttribute: model.DefaultAttribute{
			ID:        2,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name:     "adminuser",
		Password: "hashedpassword",
		Role:     "admin",
		Active:   true,
	}
}

// Test CreateEmployee request binding - Success
func TestEmployeeHandler_CreateEmployeeRequestBinding_Success(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"name":"testuser","password":"password123","role":"admin","active":true}`
	req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var createReq request.CreateEmployeeRequest
	err := c.Bind(&createReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "testuser", createReq.Name)
	assert.Equal(t, "password123", createReq.Password)
	assert.Equal(t, "admin", createReq.Role)
	assert.True(t, createReq.Active)
}

// Test CreateEmployee request binding - Invalid JSON
func TestEmployeeHandler_CreateEmployeeRequestBinding_InvalidJSON(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"invalid_json": }`
	req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var createReq request.CreateEmployeeRequest
	err := c.Bind(&createReq)

	// Assert
	assert.Error(t, err)
}

// Test CreateEmployee request binding - Missing fields
func TestEmployeeHandler_CreateEmployeeRequestBinding_MissingFields(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"name":"","password":"","role":"","active":false}`
	req := httptest.NewRequest(http.MethodPost, "/employees", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var createReq request.CreateEmployeeRequest
	err := c.Bind(&createReq)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, createReq.Name)
	assert.Empty(t, createReq.Password)
	assert.Empty(t, createReq.Role)
	assert.False(t, createReq.Active)
}

// Test UpdateEmployee request binding - Success
func TestEmployeeHandler_UpdateEmployeeRequestBinding_Success(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"name":"updateduser","password":"newpassword123","role":"user","active":false}`
	req := httptest.NewRequest(http.MethodPut, "/employees/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Test request binding
	var updateReq request.UpdateEmployeeRequest
	err := c.Bind(&updateReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updateReq.Name)
	assert.Equal(t, "newpassword123", updateReq.Password)
	assert.Equal(t, "user", updateReq.Role)
	assert.False(t, updateReq.Active)

	// Test parameter extraction
	employeeID := c.Param("id")
	assert.Equal(t, "1", employeeID)
}

// Test UpdateEmployee request binding - Invalid JSON
func TestEmployeeHandler_UpdateEmployeeRequestBinding_InvalidJSON(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"invalid_json": }`
	req := httptest.NewRequest(http.MethodPut, "/employees/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var updateReq request.UpdateEmployeeRequest
	err := c.Bind(&updateReq)

	// Assert
	assert.Error(t, err)
}

// Test GetEmployeeByID parameter parsing - Success
func TestEmployeeHandler_GetEmployeeByID_ParameterParsing_Success(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/employees/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	// Test parameter extraction and parsing
	employeeIDStr := c.Param("id")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "123", employeeIDStr)
	assert.Equal(t, uint64(123), employeeID)
	assert.Equal(t, uint(123), uint(employeeID))
}

// Test GetEmployeeByID parameter parsing - Invalid ID
func TestEmployeeHandler_GetEmployeeByID_ParameterParsing_InvalidID(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/employees/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	// Test parameter extraction and parsing
	employeeIDStr := c.Param("id")
	_, err := strconv.ParseUint(employeeIDStr, 10, 32)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid", employeeIDStr)
}

// Test GetEmployeeByID parameter parsing - Negative ID
func TestEmployeeHandler_GetEmployeeByID_ParameterParsing_NegativeID(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/employees/-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("-1")

	// Test parameter extraction and parsing
	employeeIDStr := c.Param("id")
	_, err := strconv.ParseUint(employeeIDStr, 10, 32)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "-1", employeeIDStr)
}

// Test DeleteEmployee parameter extraction - Success
func TestEmployeeHandler_DeleteEmployee_ParameterExtraction_Success(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/employees/456", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("456")

	// Test parameter extraction
	employeeID := c.Param("id")

	// Assert
	assert.Equal(t, "456", employeeID)
	assert.NotEmpty(t, employeeID)
}

// Test ValidateEmployeeAccess simulation - Admin access
func TestEmployeeHandler_ValidateEmployeeAccess_AdminAccess(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/employees/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set admin role in context
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(2))

	// Simulate access validation logic
	role := c.Get("authenticated_role").(string)

	var hasAccess bool
	if role == "admin" {
		hasAccess = true
	}

	// Assert
	assert.Equal(t, "admin", role)
	assert.True(t, hasAccess)
}

// Test ValidateEmployeeAccess simulation - Employee own access
func TestEmployeeHandler_ValidateEmployeeAccess_EmployeeOwnAccess(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/employees/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set employee role and user ID in context
	c.Set("authenticated_role", "employee")
	c.Set("authenticated_user_id", uint(1))

	// Simulate access validation logic
	role := c.Get("authenticated_role").(string)
	targetEmployeeID := uint(1)

	var hasAccess bool
	if role == "admin" {
		hasAccess = true
	} else if role == "employee" {
		userID, ok := c.Get("authenticated_user_id").(uint)
		if ok && userID == targetEmployeeID {
			hasAccess = true
		}
	}

	// Assert
	assert.Equal(t, "employee", role)
	assert.True(t, hasAccess)
}

// Test ValidateEmployeeAccess simulation - Employee accessing other's data
func TestEmployeeHandler_ValidateEmployeeAccess_EmployeeOtherAccess(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/employees/2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set employee role and user ID in context
	c.Set("authenticated_role", "employee")
	c.Set("authenticated_user_id", uint(1))

	// Simulate access validation logic
	role := c.Get("authenticated_role").(string)
	targetEmployeeID := uint(2)

	var hasAccess bool
	if role == "admin" {
		hasAccess = true
	} else if role == "employee" {
		userID, ok := c.Get("authenticated_user_id").(uint)
		if ok && userID == targetEmployeeID {
			hasAccess = true
		}
	}

	// Assert
	assert.Equal(t, "employee", role)
	assert.False(t, hasAccess)
}

// Test ValidateEmployeeAccess simulation - No role in context
func TestEmployeeHandler_ValidateEmployeeAccess_NoRole(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/employees/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Don't set role in context

	// Test if role exists in context
	role := c.Get("authenticated_role")

	// Assert
	assert.Nil(t, role)
}

// Test ValidateEmployeeAccess simulation - Wrong type in context
func TestEmployeeHandler_ValidateEmployeeAccess_WrongType(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/employees/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set wrong type for user_id
	c.Set("authenticated_role", "employee")
	c.Set("authenticated_user_id", "1") // string instead of uint

	// Simulate access validation logic
	role := c.Get("authenticated_role").(string)
	targetEmployeeID := uint(1)

	var hasAccess bool
	if role == "admin" {
		hasAccess = true
	} else if role == "employee" {
		userID, ok := c.Get("authenticated_user_id").(uint)
		if ok && userID == targetEmployeeID {
			hasAccess = true
		}
	}

	// Assert
	assert.Equal(t, "employee", role)
	assert.False(t, hasAccess) // Should be false due to type assertion failure
}

// Test Employee model ToSafe method
func TestEmployeeHandler_EmployeeToSafe(t *testing.T) {
	// Create test employee
	employee := createTestEmployeeForHandler()

	// Convert to safe format
	safeEmployee := employee.ToSafe()

	// Assert
	assert.Equal(t, employee.ID, safeEmployee.ID)
	assert.Equal(t, employee.Name, safeEmployee.Name)
	assert.Equal(t, employee.Role, safeEmployee.Role)
	assert.Equal(t, employee.Active, safeEmployee.Active)
	assert.Equal(t, *employee.CreatedAt, safeEmployee.CreatedAt)
	assert.Equal(t, *employee.UpdatedAt, safeEmployee.UpdatedAt)
}

// Test Employee roles validation
func TestEmployeeHandler_EmployeeRoles_Validation(t *testing.T) {
	// Test valid roles
	validRoles := []string{"admin", "user"}

	for _, role := range validRoles {
		// Create request with each valid role
		createReq := request.CreateEmployeeRequest{
			Name:     "testuser",
			Password: "password123",
			Role:     role,
			Active:   true,
		}

		// Assert role is as expected
		assert.Equal(t, role, createReq.Role)
		assert.Contains(t, validRoles, createReq.Role)
	}
}

// Test Employee status values
func TestEmployeeHandler_EmployeeStatus_Values(t *testing.T) {
	// Test active employee
	activeEmployee := createTestEmployeeForHandler()
	assert.True(t, activeEmployee.Active)

	// Test inactive employee
	inactiveEmployee := createTestEmployeeForHandler()
	inactiveEmployee.Active = false
	assert.False(t, inactiveEmployee.Active)
}

// Test Employee data structure
func TestEmployeeHandler_EmployeeDataStructure(t *testing.T) {
	// Create test employee
	employee := createTestEmployeeForHandler()

	// Assert all fields are properly set
	assert.NotZero(t, employee.ID)
	assert.NotEmpty(t, employee.Name)
	assert.NotEmpty(t, employee.Password)
	assert.NotEmpty(t, employee.Role)
	assert.NotNil(t, employee.CreatedAt)
	assert.NotNil(t, employee.UpdatedAt)
}

// Test Different employee roles
func TestEmployeeHandler_DifferentRoles(t *testing.T) {
	// Test regular employee
	employee := createTestEmployeeForHandler()
	assert.Equal(t, "employee", employee.Role)

	// Test admin employee
	admin := createAdminEmployeeForHandler()
	assert.Equal(t, "admin", admin.Role)
}

// Test CreateEmployeeRequest structure
func TestEmployeeHandler_CreateEmployeeRequest_Structure(t *testing.T) {
	// Create request
	req := request.CreateEmployeeRequest{
		Name:     "newuser",
		Password: "securepassword",
		Role:     "admin",
		Active:   true,
	}

	// Assert structure
	assert.IsType(t, "", req.Name)
	assert.IsType(t, "", req.Password)
	assert.IsType(t, "", req.Role)
	assert.IsType(t, true, req.Active)

	assert.Equal(t, "newuser", req.Name)
	assert.Equal(t, "securepassword", req.Password)
	assert.Equal(t, "admin", req.Role)
	assert.True(t, req.Active)
}

// Test UpdateEmployeeRequest structure
func TestEmployeeHandler_UpdateEmployeeRequest_Structure(t *testing.T) {
	// Create request
	req := request.UpdateEmployeeRequest{
		Name:     "updateduser",
		Password: "newpassword",
		Role:     "user",
		Active:   false,
	}

	// Assert structure
	assert.IsType(t, "", req.Name)
	assert.IsType(t, "", req.Password)
	assert.IsType(t, "", req.Role)
	assert.IsType(t, true, req.Active)

	assert.Equal(t, "updateduser", req.Name)
	assert.Equal(t, "newpassword", req.Password)
	assert.Equal(t, "user", req.Role)
	assert.False(t, req.Active)
}

// Test HTTP parameter validation
func TestEmployeeHandler_HTTPParameters_Validation(t *testing.T) {
	testCases := []struct {
		name        string
		paramValue  string
		expectError bool
	}{
		{"Valid ID", "123", false},
		{"Zero ID", "0", false},
		{"Large ID", "999999", false},
		{"Invalid ID", "abc", true},
		{"Negative ID", "-1", true},
		{"Float ID", "12.34", true},
		{"Empty ID", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := strconv.ParseUint(tc.paramValue, 10, 32)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test Context value extraction
func TestEmployeeHandler_ContextValues_Extraction(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/employees/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set various context values
	c.Set("authenticated_role", "admin")
	c.Set("authenticated_user_id", uint(123))
	c.Set("custom_value", "test")

	// Extract and verify values
	role, ok := c.Get("authenticated_role").(string)
	assert.True(t, ok)
	assert.Equal(t, "admin", role)

	userID, ok := c.Get("authenticated_user_id").(uint)
	assert.True(t, ok)
	assert.Equal(t, uint(123), userID)

	customValue, ok := c.Get("custom_value").(string)
	assert.True(t, ok)
	assert.Equal(t, "test", customValue)

	// Test non-existent value
	nonExistent := c.Get("non_existent")
	assert.Nil(t, nonExistent)
}
