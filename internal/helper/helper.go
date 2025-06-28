package helper

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/helper/postgre"
	"github.com/yourname/payslip-system/internal/helper/response"
	"github.com/yourname/payslip-system/internal/middleware"
	"gorm.io/gorm"
)

// NewHelper ...
type NewHelper struct {
	Response     response.Interface
	TimeLocation *time.Location
	DB           postgre.Database
}

// GetAuditableDB returns a database instance with audit fields automatically set
func GetAuditableDB(c echo.Context, db *gorm.DB) *middleware.AuditableDB {
	userID, ok := c.Get("user_id").(int)
	if !ok || userID <= 0 {
		// Fallback to regular DB if no user context
		return middleware.NewAuditableDB(db, 0)
	}

	return middleware.NewAuditableDB(db, uint(userID))
}

// ValidateEmployeeAccess checks if the current user can access employee-specific data
func ValidateEmployeeAccess(c echo.Context, targetEmployeeID uint) bool {
	role := c.Get("authenticated_role").(string)

	// Admins can access any employee's data
	if role == "admin" {
		return true
	}

	// Employees can only access their own data
	if role == "employee" {
		userID, ok := c.Get("authenticated_user_id").(uint)
		if ok && userID == targetEmployeeID {
			return true
		}
	}

	return false
}
