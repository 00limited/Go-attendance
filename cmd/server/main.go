package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yourname/payslip-system/internal/config"
	"github.com/yourname/payslip-system/internal/database"
	"github.com/yourname/payslip-system/internal/model"
	"github.com/yourname/payslip-system/internal/routes"
	"github.com/yourname/payslip-system/internal/seed"
)

// CustomValidator is a custom validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates structs
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	config.LoadEnv()

	db := database.Connect()
	seed.Run(db)
	db.Debug()
	db.AutoMigrate(&model.Employee{}, &model.Attendance{}, &model.Overtime{}, &model.Reimbursement{}, &model.Payslip{})

	defer database.Close(db)

	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Set custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Register routes
	routes.SetupRoutes(e)

	port := config.GetEnv("PORT", "8080")
	e.Logger.Infof("🚀 Starting JWT-secured payroll server on port %s", port)
	e.Logger.Infof("📋 Default admin credentials: Admin / admin123")
	e.Logger.Infof("📋 Default employee password: password123")
	e.Logger.Infof("📖 API Documentation: See JWT_AUTHENTICATION.md and API_TESTING_GUIDE.md")

	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatalf("Failed to start server: %v", err)
	}
}
