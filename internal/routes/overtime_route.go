package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/handler"
	mymiddleware "github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/repository"
)

func (t *NewRoute) OvertimeRoutes(c *echo.Group) {
	// Add JWT middleware to protect all overtime routes
	c.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  mymiddleware.JWT_SECRET,
		TokenLookup: "header:Authorization:Bearer ",
	}))
	c.Use(mymiddleware.HeaderMiddleware)
	c.Use(mymiddleware.AuditMiddleware()) // Add audit middleware

	h := handler.OvertimeHandler{
		Helper:       t.Helper,
		Response:     t.Response,
		BaseRepo:     repository.NewBaseRepository(t.DB),
		OvertimeRepo: repository.NewOvertimeRepository(t.DB),
	}

	// Employee or Admin routes (employees can create their own overtime)
	employeeGroup := c.Group("")
	employeeGroup.Use(mymiddleware.EmployeeOrAdmin(t.Response))
	employeeGroup.POST("/create", h.CreateOvertime)
}
