package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/handler"
	mymiddleware "github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/repository"
)

func (t *NewRoute) AttendanceRoutes(c *echo.Group) {
	// Add JWT middleware to protect all attendance routes
	c.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  mymiddleware.JWT_SECRET,
		TokenLookup: "header:Authorization:Bearer ",
	}))
	c.Use(mymiddleware.HeaderMiddleware)
	c.Use(mymiddleware.AuditMiddleware()) // Add audit middleware

	h := handler.AttendanceHandler{
		Helper:         t.Helper,
		Response:       t.Response,
		BaseRepo:       repository.NewBaseRepository(t.DB),
		AttendanceRepo: repository.NewAttendanceRepository(t.DB),
	}

	// Employee or Admin routes (employees can manage their own attendance)
	employeeGroup := c.Group("")
	employeeGroup.Use(mymiddleware.EmployeeOrAdmin(t.Response))
	employeeGroup.POST("/check-in", h.CheckinAttendancePeriod)
	employeeGroup.POST("/check-out", h.CheckOutAttendancePeriod)
}
