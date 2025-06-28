package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/helper"
	"github.com/yourname/payslip-system/internal/helper/response"
	"github.com/yourname/payslip-system/internal/repository"
	"github.com/yourname/payslip-system/internal/usecases"
)

type PayrollHandler struct {
	payslipRepo    repository.PayslipRepository
	payrollUsecase *usecases.PayrollUsecase
	response       response.Interface
}

func NewPayrollHandler(payslipRepo repository.PayslipRepository, payrollUsecase *usecases.PayrollUsecase, response response.Interface) *PayrollHandler {
	return &PayrollHandler{
		payslipRepo:    payslipRepo,
		payrollUsecase: payrollUsecase,
		response:       response,
	}
}

// RunPayrollForAllEmployees processes payroll for all active employees
func (h *PayrollHandler) RunPayrollForAllEmployees(c echo.Context) error {
	var req request.PayrollRequest
	if err := c.Bind(&req); err != nil {
		return h.response.SendBadRequest(c, "Invalid request body", err.Error())
	}

	// Validate the request
	if req.PayPeriodEnd.Before(req.PayPeriodStart) {
		return h.response.SendBadRequest(c, "Pay period end must be after start date", nil)
	}

	// Get auditable DB instance
	auditDB := helper.GetAuditableDB(c, h.payslipRepo.GetDB())

	// Process payroll using usecase with audit trail
	processedPayslips, errors := h.payrollUsecase.ProcessAllEmployeesPayrollWithAudit(req, auditDB)

	result := map[string]interface{}{
		"processed_count": len(processedPayslips),
		"error_count":     len(errors),
		"payslips":        processedPayslips,
	}

	if len(errors) > 0 {
		result["errors"] = errors
	}

	return h.response.SendSuccess(c, "Payroll processed", result)
}

// RunPayrollForEmployee processes payroll for a specific employee
func (h *PayrollHandler) RunPayrollForEmployee(c echo.Context) error {
	var req request.PayrollEmployeeRequest
	if err := c.Bind(&req); err != nil {
		return h.response.SendBadRequest(c, "Invalid request body", err.Error())
	}

	// Validate the request
	if req.PayPeriodEnd.Before(req.PayPeriodStart) {
		return h.response.SendBadRequest(c, "Pay period end must be after start date", nil)
	}

	payrollReq := request.PayrollRequest{
		PayPeriodStart: req.PayPeriodStart,
		PayPeriodEnd:   req.PayPeriodEnd,
		BasicSalary:    req.BasicSalary,
		OvertimeRate:   req.OvertimeRate,
	}

	// Get auditable DB instance
	auditDB := helper.GetAuditableDB(c, h.payslipRepo.GetDB())

	payslip, err := h.payrollUsecase.ProcessEmployeePayrollWithAudit(req.EmployeeID, payrollReq, auditDB)
	if err != nil {
		return h.response.SendError(c, "Failed to process payroll", err.Error())
	}

	return h.response.SendSuccess(c, "Payroll processed for employee", payslip)
}

// GetPayslipsByEmployee retrieves payslips for a specific employee
func (h *PayrollHandler) GetPayslipsByEmployee(c echo.Context) error {
	employeeID := c.Param("id")
	if employeeID == "" {
		return h.response.SendBadRequest(c, "Employee ID is required", nil)
	}

	// Convert string to uint
	var empID uint
	if _, err := fmt.Sscanf(employeeID, "%d", &empID); err != nil {
		return h.response.SendBadRequest(c, "Invalid employee ID format", err.Error())
	}

	// Check authorization - employees can only access their own payslips
	if !helper.ValidateEmployeeAccess(c, empID) {
		return h.response.SendCustomResponse(c, 403, "Access denied. You can only access your own payslips.", nil)
	}

	// Get employee to verify existence
	employee, err := h.payslipRepo.GetEmployeeByID(empID)
	if err != nil {
		return h.response.SendError(c, "Employee not found", err.Error())
	}

	// Get all payslips for the employee
	payslips, err := h.payslipRepo.GetPayslipsByEmployee(empID)
	if err != nil {
		return h.response.SendError(c, "Failed to retrieve payslips", err.Error())
	}

	// Convert to response format
	var payslipList []map[string]interface{}
	for _, payslip := range payslips {
		payslipList = append(payslipList, map[string]interface{}{
			"payslip_id":       payslip.ID,
			"pay_period_start": payslip.PayPeriodStart,
			"pay_period_end":   payslip.PayPeriodEnd,
			"total_amount":     payslip.TotalAmount,
			"status":           payslip.Status,
			"processed_at":     payslip.ProcessedAt,
		})
	}

	result := map[string]interface{}{
		"employee_id":   employee.ID,
		"employee_name": employee.Name,
		"payslips":      payslipList,
		"total_count":   len(payslips),
	}

	return h.response.SendSuccess(c, "Payslips retrieved successfully", result)
}

// GetDetailedPayslip generates a detailed payslip with all breakdowns
func (h *PayrollHandler) GetDetailedPayslip(c echo.Context) error {
	payslipID := c.Param("payslip_id")
	if payslipID == "" {
		return h.response.SendBadRequest(c, "Payslip ID is required", nil)
	}

	// Convert string to uint
	var pID uint
	if _, err := fmt.Sscanf(payslipID, "%d", &pID); err != nil {
		return h.response.SendBadRequest(c, "Invalid payslip ID format", err.Error())
	}

	// Get payslip
	payslip, err := h.payslipRepo.GetPayslipByID(pID)
	if err != nil {
		return h.response.SendError(c, "Payslip not found", err.Error())
	}

	// Check authorization - employees can only access their own payslips
	if !helper.ValidateEmployeeAccess(c, payslip.EmployeeID) {
		return h.response.SendCustomResponse(c, 403, "Access denied. You can only access your own payslips.", nil)
	}

	// Get employee details
	employee, err := h.payslipRepo.GetEmployeeByID(payslip.EmployeeID)
	if err != nil {
		return h.response.SendError(c, "Employee not found", err.Error())
	}

	// Get attendance breakdown
	attendances, err := h.payslipRepo.GetAttendanceForPeriod(payslip.EmployeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd)
	if err != nil {
		return h.response.SendError(c, "Failed to get attendance records", err.Error())
	}

	// Get overtime breakdown
	dateStart := payslip.PayPeriodStart.Format("2006-01-02")
	dateEnd := payslip.PayPeriodEnd.Format("2006-01-02")
	overtimes, err := h.payslipRepo.GetOvertimeForPeriod(payslip.EmployeeID, dateStart, dateEnd)
	if err != nil {
		return h.response.SendError(c, "Failed to get overtime records", err.Error())
	}

	// Get reimbursement breakdown
	reimbursements, err := h.payslipRepo.GetApprovedReimbursementsForPeriod(payslip.EmployeeID, payslip.PayPeriodStart, payslip.PayPeriodEnd)
	if err != nil {
		return h.response.SendError(c, "Failed to get reimbursement records", err.Error())
	}

	// Build detailed response
	detailedPayslip := h.payrollUsecase.BuildDetailedPayslipResponse(payslip, employee, attendances, overtimes, reimbursements)

	return h.response.SendSuccess(c, "Detailed payslip generated successfully", detailedPayslip)
}

// GetPayrollSummary generates a summary of all employee payslips for a period
func (h *PayrollHandler) GetPayrollSummary(c echo.Context) error {
	var req request.PayrollSummaryRequest
	if err := c.Bind(&req); err != nil {
		return h.response.SendBadRequest(c, "Invalid request body", err.Error())
	}

	// Validate the request
	if req.PayPeriodEnd.Before(req.PayPeriodStart) {
		return h.response.SendBadRequest(c, "Pay period end must be after start date", nil)
	}

	// Get all payslips for the period
	payslips, err := h.payslipRepo.GetPayslipsByPeriod(req.PayPeriodStart, req.PayPeriodEnd)
	if err != nil {
		return h.response.SendError(c, "Failed to retrieve payslips", err.Error())
	}

	// Build summary data
	summary := h.payrollUsecase.BuildPayrollSummary(payslips)

	return h.response.SendSuccess(c, "Payroll summary generated successfully", summary)
}
