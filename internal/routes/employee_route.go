package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/handler"
	mymiddleware "github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/repository"
)

// EmployeeRoutes initializes the routes for employee management
func (t *NewRoute) EmployeeRoutes(c *echo.Group) {
	// Add JWT middleware to protect all employee routes
	c.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  mymiddleware.JWT_SECRET,
		TokenLookup: "header:Authorization:Bearer ",
	}))
	c.Use(mymiddleware.HeaderMiddleware)
	c.Use(mymiddleware.AuditMiddleware()) // Add audit middleware after JWT validation

	h := handler.EmployeeHandler{
		Helper:       t.Helper,
		Response:     t.Response,
		BaseRepo:     repository.NewBaseRepository(t.DB),
		EmployeeRepo: repository.NewEmployeeRepository(t.DB),
	}

	// Admin-only routes
	adminGroup := c.Group("")
	adminGroup.Use(mymiddleware.AdminOnly(t.Response))
	adminGroup.POST("/create", h.CreateEmployee)
	adminGroup.GET("/get-all-employee", h.GetAllEmployees)
	adminGroup.PUT("/edit/:id", h.EditEmployee)
	adminGroup.DELETE("/delete/:id", h.DeleteEmployee)

	// Employee or Admin routes (employees can view their own data)
	employeeGroup := c.Group("")
	employeeGroup.Use(mymiddleware.EmployeeOrAdmin(t.Response))
	employeeGroup.GET("/profile/:id", h.GetEmployeeByID) // Use existing method
}
