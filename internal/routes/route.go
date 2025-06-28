package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/database"
	"github.com/yourname/payslip-system/internal/helper"
	responseHelper "github.com/yourname/payslip-system/internal/helper/response"
	"gorm.io/gorm"
)

// NewRoute Handler
type NewRoute struct {
	Echo     *echo.Echo
	Response responseHelper.Interface
	Helper   helper.NewHelper
	DB       *gorm.DB
}

func SetupRoutes(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Server is up and running",
		})
	})

	api := e.Group("/api/v1")

	// Initialize NewRoute
	// You need to provide a valid *gorm.DB instance here, e.g., from a parameter or global variable
	newRoute := &NewRoute{
		Echo:     e,
		Response: responseHelper.NewResponse(),
		Helper:   helper.NewHelper{},
		DB:       database.DB,
	}

	// Authentication Routes (public)
	authGroup := api.Group("/auth")
	newRoute.AuthRoutes(authGroup)

	// Employee Routes (protected)
	employeeGroup := api.Group("/employee")
	newRoute.EmployeeRoutes(employeeGroup)
	attendanceGroup := api.Group("/attendance")
	// Register Attendance Routes
	newRoute.AttendanceRoutes(attendanceGroup)
	// Add other route groups as needed

	// Overtime Routes
	overtimeGroup := api.Group("/overtime")
	newRoute.OvertimeRoutes(overtimeGroup)

	// Reimbusement Routes
	reimbusementGroup := api.Group("/reimbusement")
	newRoute.ReimbusementRoutes(reimbusementGroup)

	// Payroll Routes
	payrollGroup := api.Group("/payroll")
	newRoute.PayrollRoutes(payrollGroup)
}
