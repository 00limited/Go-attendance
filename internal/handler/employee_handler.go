package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/helper"
	"github.com/yourname/payslip-system/internal/helper/response"
	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/repository"
	"gorm.io/gorm"
)

// EmployeeHandler handles employee-related requests.
type EmployeeHandler struct {
	Helper   helper.NewHelper
	DB       *gorm.DB
	Response response.Interface

	BaseRepo     repository.BaseRepositoryInterface
	EmployeeRepo repository.EmployeeRepository
}

// NewEmployeeHandler creates a new instance of EmployeeHandler.
func (h *EmployeeHandler) CreateEmployee(c echo.Context) error {
	req := request.CreateEmployeeRequest{}
	if err := c.Bind(&req); err != nil {
		return h.Response.SendError(c, err.Error(), "Invalid request data")
	}

	// Get auditable database instance
	auditDB := helper.GetAuditableDB(c, h.BaseRepo.GetDB())

	_, err := h.EmployeeRepo.CreateEmployeeWithAudit(req, auditDB)
	if err != nil {
		return h.Response.SendError(c, err.Error(), "Failed to create employee")
	}

	return h.Response.SendSuccess(c, "Employee created successfully", nil)
}

func (h *EmployeeHandler) GetAllEmployees(c echo.Context) error {
	employees, err := h.EmployeeRepo.GetAllEmployees()
	if err != nil {
		return h.Response.SendError(c, err.Error(), "Failed to retrieve employees")
	}
	return h.Response.SendSuccess(c, "Employees retrieved successfully", employees)
}

// EditEmployee updates an employee with audit tracking
func (h *EmployeeHandler) EditEmployee(c echo.Context) error {
	// Implementation for editing an employee
	req := request.UpdateEmployeeRequest{}
	if err := c.Bind(&req); err != nil {
		return h.Response.SendError(c, err.Error(), "Invalid request data")
	}

	employeeID := c.Param("id")

	// Get auditable database instance
	auditDB := helper.GetAuditableDB(c, h.BaseRepo.GetDB())

	_, err := h.EmployeeRepo.UpdateEmployeeWithAudit(employeeID, req, auditDB)
	if err != nil {
		return h.Response.SendError(c, err.Error(), "Failed to update employee")
	}
	return h.Response.SendSuccess(c, "Employee updated successfully", nil)
}

// DeleteEmployee deletes an employee with audit tracking
func (h *EmployeeHandler) DeleteEmployee(c echo.Context) error {
	// Implementation for deleting an employee
	employeeID := c.Param("id")

	// Get auditable database instance
	auditDB := helper.GetAuditableDB(c, h.BaseRepo.GetDB())

	err := h.EmployeeRepo.DeleteEmployeeWithAudit(employeeID, auditDB)
	if err != nil {
		return h.Response.SendError(c, err.Error(), "Failed to delete employee")
	}
	return h.Response.SendSuccess(c, "Employee deleted successfully", nil)
}

// GetEmployeeByID retrieves an employee by ID with access control
func (h *EmployeeHandler) GetEmployeeByID(c echo.Context) error {
	employeeIDStr := c.Param("id")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		return h.Response.SendBadRequest(c, "Invalid employee ID", err.Error())
	}

	// Check if user has access to this employee's data
	if !middleware.ValidateEmployeeAccess(c, uint(employeeID)) {
		return h.Response.SendCustomResponse(c, 403, "Access denied. You can only view your own profile.", nil)
	}

	employee, err := h.EmployeeRepo.GetEmployeeByID(uint(employeeID))
	if err != nil {
		return h.Response.SendError(c, "Employee not found", err.Error())
	}

	// Return safe employee data (without password)
	return h.Response.SendSuccess(c, "Employee retrieved successfully", employee.ToSafe())
}
