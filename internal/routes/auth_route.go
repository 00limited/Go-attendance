package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/handler"
	mymiddleware "github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/repository"
)

// AuthRoutes sets up authentication routes
func (nr *NewRoute) AuthRoutes(group *echo.Group) {
	// Initialize repositories
	employeeRepo := repository.NewEmployeeRepository(nr.DB)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(employeeRepo, nr.Response)

	// Public routes (no authentication required)
	group.POST("/login", authHandler.Login)

	// Protected routes (authentication required)
	protected := group.Group("")
	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  mymiddleware.JWT_SECRET,
		TokenLookup: "header:Authorization:Bearer ",
	}))
	protected.Use(mymiddleware.HeaderMiddleware)

	// Profile and token management routes
	protected.GET("/profile", authHandler.GetProfile)
	protected.POST("/refresh", authHandler.RefreshToken)
}
