package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/handler"
	mymiddleware "github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/repository"
)

func (t *NewRoute) ReimbusementRoutes(c *echo.Group) {
	// Add JWT middleware to protect all reimbursement routes
	c.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  mymiddleware.JWT_SECRET,
		TokenLookup: "header:Authorization:Bearer ",
	}))
	c.Use(mymiddleware.HeaderMiddleware)
	c.Use(mymiddleware.AuditMiddleware()) // Add audit middleware

	h := handler.ReimbusementHandler{
		Helper:           t.Helper,
		Response:         t.Response,
		BaseRepo:         repository.NewBaseRepository(t.DB),
		ReimbusementRepo: repository.NewReimbusementRepository(t.DB),
	}

	// Employee or Admin routes (employees can create their own reimbursements)
	employeeGroup := c.Group("")
	employeeGroup.Use(mymiddleware.EmployeeOrAdmin(t.Response))
	employeeGroup.POST("/create", h.CreateReimbusement)
}
