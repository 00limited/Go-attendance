package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/handler"
	mymiddleware "github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/repository"
	"github.com/yourname/payslip-system/internal/usecases"
)

// PayrollRoutes sets up the payroll-related routes
func (t *NewRoute) PayrollRoutes(c *echo.Group) {
	// Add JWT middleware to protect all payroll routes
	c.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  mymiddleware.JWT_SECRET,
		TokenLookup: "header:Authorization:Bearer ",
	}))
	c.Use(mymiddleware.HeaderMiddleware)
	c.Use(mymiddleware.AuditMiddleware()) // Add audit middleware

	payslipRepo := repository.NewPayslipRepository(t.DB)
	employeeRepo := repository.NewEmployeeRepository(t.DB)
	payrollUsecase := usecases.NewPayrollUsecase(payslipRepo, employeeRepo)
	h := handler.NewPayrollHandler(
		payslipRepo,
		payrollUsecase,
		t.Response,
	)

	// Admin-only payroll management routes
	adminGroup := c.Group("")
	adminGroup.Use(mymiddleware.AdminOnly(t.Response))

	// Run payroll for all employees (Admin only)
	adminGroup.POST("/run", h.RunPayrollForAllEmployees)

	// Run payroll for specific employee (Admin only)
	adminGroup.POST("/run/employee", h.RunPayrollForEmployee)

	// Get payroll summary for admin overview (Admin only)
	adminGroup.POST("/summary", h.GetPayrollSummary)

	// Employee and Admin accessible routes
	employeeGroup := c.Group("")
	employeeGroup.Use(mymiddleware.EmployeeOrAdmin(t.Response))

	// Get list of payslips for an employee (Employee can access own, Admin can access any)
	employeeGroup.GET("/employee/:id/payslips", h.GetPayslipsByEmployee)

	// Get detailed payslip with full breakdown (Employee can access own, Admin can access any)
	employeeGroup.GET("/payslip/:payslip_id/details", h.GetDetailedPayslip)
}
