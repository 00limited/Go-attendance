package handler

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/dto/request"
	dto_response "github.com/yourname/payslip-system/internal/dto/response"
	"github.com/yourname/payslip-system/internal/helper/response"
	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"github.com/yourname/payslip-system/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	employeeRepo repository.EmployeeRepository
	response     response.Interface
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(employeeRepo repository.EmployeeRepository, response response.Interface) *AuthHandler {
	return &AuthHandler{
		employeeRepo: employeeRepo,
		response:     response,
	}
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c echo.Context) error {
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return h.response.SendBadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return h.response.SendBadRequest(c, "Validation failed", err.Error())
	}

	// Find employee by name
	employee, err := h.employeeRepo.GetEmployeeByName(req.Name)
	if err != nil {
		return h.response.SendUnauthorized(c, "Invalid credentials", nil)
	}

	// Check if employee is active
	if !employee.Active {
		return h.response.SendUnauthorized(c, "Account is deactivated", nil)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(req.Password)); err != nil {
		return h.response.SendUnauthorized(c, "Invalid credentials", nil)
	}

	// Generate JWT token
	token, expiresAt, err := h.generateJWTToken(employee)
	if err != nil {
		return h.response.SendError(c, "Failed to generate token", err.Error())
	}

	// Prepare response
	loginResponse := dto_response.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: expiresAt,
		User: dto_response.UserInfo{
			ID:     employee.ID,
			Name:   employee.Name,
			Role:   employee.Role,
			Active: employee.Active,
		},
	}

	return h.response.SendSuccess(c, "Login successful", loginResponse)
}

// generateJWTToken creates a new JWT token for the authenticated user
func (h *AuthHandler) generateJWTToken(employee *model.Employee) (string, time.Time, error) {
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
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// GetProfile returns the current user's profile information
func (h *AuthHandler) GetProfile(c echo.Context) error {
	// Get user information from context (set by HeaderMiddleware)
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return h.response.SendUnauthorized(c, "Invalid token", nil)
	}

	// Get employee details
	employee, err := h.employeeRepo.GetEmployeeByID(uint(userID))
	if err != nil {
		return h.response.SendError(c, "User not found", err.Error())
	}

	// Return safe employee data
	return h.response.SendSuccess(c, "Profile retrieved successfully", employee.ToSafe())
}

// RefreshToken generates a new token for the authenticated user
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// Get user information from context
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return h.response.SendUnauthorized(c, "Invalid token", nil)
	}

	// Get employee details
	employee, err := h.employeeRepo.GetEmployeeByID(uint(userID))
	if err != nil {
		return h.response.SendError(c, "User not found", err.Error())
	}

	// Check if employee is still active
	if !employee.Active {
		return h.response.SendUnauthorized(c, "Account is deactivated", nil)
	}

	// Generate new token
	token, expiresAt, err := h.generateJWTToken(employee)
	if err != nil {
		return h.response.SendError(c, "Failed to refresh token", err.Error())
	}

	// Prepare response
	refreshResponse := dto_response.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: expiresAt,
		User: dto_response.UserInfo{
			ID:     employee.ID,
			Name:   employee.Name,
			Role:   employee.Role,
			Active: employee.Active,
		},
	}

	return h.response.SendSuccess(c, "Token refreshed successfully", refreshResponse)
}
