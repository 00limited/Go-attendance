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
