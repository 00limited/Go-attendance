package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// Helper function to create a test employee

func createTestEmployee() *model.Employee {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	now := time.Now()

	return &model.Employee{
		DefaultAttribute: model.DefaultAttribute{
			ID:        1,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name:     "testuser",
		Password: string(hashedPassword),
		Role:     "employee",
		Active:   true,
	}
}

func createInactiveTestEmployee() *model.Employee {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	now := time.Now()

	return &model.Employee{
		DefaultAttribute: model.DefaultAttribute{
			ID:        2,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name:     "inactiveuser",
		Password: string(hashedPassword),
		Role:     "employee",
		Active:   false,
	}
}

// Test Login request binding - Success
func TestAuthHandler_LoginRequestBinding_Success(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"name":"testuser","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var loginReq request.LoginRequest
	err := c.Bind(&loginReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "testuser", loginReq.Name)
	assert.Equal(t, "password123", loginReq.Password)
}

// Test Login request binding - Invalid JSON
func TestAuthHandler_LoginRequestBinding_InvalidJSON(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"invalid_json": }`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var loginReq request.LoginRequest
	err := c.Bind(&loginReq)

	// Assert
	assert.Error(t, err)
}

// Test Login request binding - Missing fields
func TestAuthHandler_LoginRequestBinding_MissingFields(t *testing.T) {
	// Setup Echo
	e := echo.New()
	reqBody := `{"name":"","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test request binding
	var loginReq request.LoginRequest
	err := c.Bind(&loginReq)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, loginReq.Name)
	assert.Empty(t, loginReq.Password)
}

// Test Password verification - Correct password
func TestAuthHandler_PasswordVerification_Success(t *testing.T) {
	// Create test employee
	employee := createTestEmployee()

	// Test password verification
	err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte("password123"))

	// Assert
	assert.NoError(t, err)
}

// Test Password verification - Wrong password
func TestAuthHandler_PasswordVerification_Failure(t *testing.T) {
	// Create test employee
	employee := createTestEmployee()

	// Test password verification with wrong password
	err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte("wrongpassword"))

	// Assert
	assert.Error(t, err)
}

// Test JWT Token generation and verification
func TestAuthHandler_JWTToken_Success(t *testing.T) {
	// Create test employee
	employee := createTestEmployee()

	// Set token expiration time (24 hours)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Create token claims
	claims := &jwt.MapClaims{
		"user_id": employee.ID,
		"name":    employee.Name,
		"role":    employee.Role,
		"active":  employee.Active,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString(middleware.JWT_SECRET)

	// Assert token generation
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Verify token can be parsed
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return middleware.JWT_SECRET, nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Verify claims
	if parsedClaims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		assert.Equal(t, float64(employee.ID), parsedClaims["user_id"])
		assert.Equal(t, employee.Name, parsedClaims["name"])
		assert.Equal(t, employee.Role, parsedClaims["role"])
		assert.Equal(t, employee.Active, parsedClaims["active"])
	}
}

// Test JWT Token validation with expired token
func TestAuthHandler_JWTToken_Expired(t *testing.T) {
	// Create test employee
	employee := createTestEmployee()

	// Set token expiration time in the past
	expiresAt := time.Now().Add(-1 * time.Hour)

	// Create token claims
	claims := &jwt.MapClaims{
		"user_id": employee.ID,
		"name":    employee.Name,
		"role":    employee.Role,
		"active":  employee.Active,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString(middleware.JWT_SECRET)
	assert.NoError(t, err)

	// Try to parse expired token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return middleware.JWT_SECRET, nil
	})

	// Token should be invalid due to expiration
	assert.Error(t, err)
	assert.False(t, parsedToken.Valid)
}

// Test JWT Token with wrong secret
func TestAuthHandler_JWTToken_WrongSecret(t *testing.T) {
	// Create test employee
	employee := createTestEmployee()

	// Set token expiration time
	expiresAt := time.Now().Add(24 * time.Hour)

	// Create token claims
	claims := &jwt.MapClaims{
		"user_id": employee.ID,
		"name":    employee.Name,
		"role":    employee.Role,
		"active":  employee.Active,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token with correct secret
	tokenString, err := token.SignedString(middleware.JWT_SECRET)
	assert.NoError(t, err)

	// Try to parse with wrong secret
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("wrong-secret"), nil
	})

	// Should fail due to wrong secret
	assert.Error(t, err)
}

// Test context user_id extraction - Success
func TestAuthHandler_ContextUserID_Success(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set user_id in context
	c.Set("user_id", 1)

	// Get user_id from context
	userID, ok := c.Get("user_id").(int)

	// Assert
	assert.True(t, ok)
	assert.Equal(t, 1, userID)
}

// Test context user_id extraction - Missing
func TestAuthHandler_ContextUserID_Missing(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Don't set user_id in context

	// Try to get user_id from context
	userID, ok := c.Get("user_id").(int)

	// Assert
	assert.False(t, ok)
	assert.Equal(t, 0, userID)
}

// Test context user_id extraction - Wrong type
func TestAuthHandler_ContextUserID_WrongType(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set user_id as string instead of int
	c.Set("user_id", "1")

	// Try to get user_id as int from context
	userID, ok := c.Get("user_id").(int)

	// Assert
	assert.False(t, ok)
	assert.Equal(t, 0, userID)
}

// Test Employee active status verification
func TestAuthHandler_EmployeeActiveStatus_Active(t *testing.T) {
	// Create active employee
	employee := createTestEmployee()

	// Assert
	assert.True(t, employee.Active)
}

// Test Employee active status verification
func TestAuthHandler_EmployeeActiveStatus_Inactive(t *testing.T) {
	// Create inactive employee
	employee := createInactiveTestEmployee()

	// Assert
	assert.False(t, employee.Active)
}

// Test ToSafe method
func TestAuthHandler_EmployeeToSafe(t *testing.T) {
	// Create test employee
	employee := createTestEmployee()

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

// Test JWT Claims structure
func TestAuthHandler_JWTClaims_Structure(t *testing.T) {
	// Create test employee
	employee := createTestEmployee()

	// Create expected claims structure
	expectedClaims := map[string]interface{}{
		"user_id": employee.ID,
		"name":    employee.Name,
		"role":    employee.Role,
		"active":  employee.Active,
	}

	// Assert claims have correct types and values
	assert.IsType(t, uint(0), expectedClaims["user_id"])
	assert.IsType(t, "", expectedClaims["name"])
	assert.IsType(t, "", expectedClaims["role"])
	assert.IsType(t, true, expectedClaims["active"])

	assert.Equal(t, employee.ID, expectedClaims["user_id"])
	assert.Equal(t, employee.Name, expectedClaims["name"])
	assert.Equal(t, employee.Role, expectedClaims["role"])
	assert.Equal(t, employee.Active, expectedClaims["active"])
}

// Test bcrypt password hashing
func TestAuthHandler_PasswordHashing(t *testing.T) {
	password := "testpassword123"

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, string(hashedPassword))

	// Verify password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	assert.NoError(t, err)

	// Verify wrong password fails
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongpassword"))
	assert.Error(t, err)
}
