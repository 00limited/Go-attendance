package usecases

import (
	"fmt"
	"time"

	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"github.com/yourname/payslip-system/internal/repository"
)

type PayrollUsecase struct {
	payslipRepo  repository.PayslipRepository
	employeeRepo repository.EmployeeRepository
}

func NewPayrollUsecase(payslipRepo repository.PayslipRepository, employeeRepo repository.EmployeeRepository) *PayrollUsecase {
	return &PayrollUsecase{
		payslipRepo:  payslipRepo,
		employeeRepo: employeeRepo,
	}
}

// ProcessEmployeePayroll handles the payroll calculation for a single employee
func (uc *PayrollUsecase) ProcessEmployeePayroll(employeeID uint, req request.PayrollRequest) (*model.Payslip, error) {
	// Check if payslip already exists for this period
	exists, err := uc.payslipRepo.CheckPayslipExists(employeeID, req.PayPeriodStart, req.PayPeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing payslip: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("payslip already exists for this period")
	}

	// Get attendance records for the period
	attendances, err := uc.payslipRepo.GetAttendanceForPeriod(employeeID, req.PayPeriodStart, req.PayPeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance records: %v", err)
	}

	// Get overtime records for the period
	dateStart := req.PayPeriodStart.Format("2006-01-02")
	dateEnd := req.PayPeriodEnd.Format("2006-01-02")
	overtimes, err := uc.payslipRepo.GetOvertimeForPeriod(employeeID, dateStart, dateEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get overtime records: %v", err)
	}

	// Get approved reimbursements for the period
	reimbursements, err := uc.payslipRepo.GetApprovedReimbursementsForPeriod(employeeID, req.PayPeriodStart, req.PayPeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get reimbursement records: %v", err)
	}

	// Calculate totals
	attendanceDays := len(attendances)
	totalOvertimeHours := uc.calculateTotalOvertimeHours(overtimes)
	totalReimbursementAmount := uc.calculateTotalReimbursementAmount(reimbursements)

	// Calculate amounts
	overtimeAmount := float64(totalOvertimeHours) * req.OvertimeRate
	totalAmount := req.BasicSalary + overtimeAmount + totalReimbursementAmount

	// Create payslip
	payslip := &model.Payslip{
		EmployeeID:          employeeID,
		PayPeriodStart:      req.PayPeriodStart,
		PayPeriodEnd:        req.PayPeriodEnd,
		BasicSalary:         req.BasicSalary,
		OvertimeHours:       totalOvertimeHours,
		OvertimeAmount:      overtimeAmount,
		ReimbursementAmount: totalReimbursementAmount,
		TotalAmount:         totalAmount,
		ProcessedAt:         time.Now(),
		Status:              "processed",
		AttendanceDays:      attendanceDays,
	}

	return uc.payslipRepo.CreatePayslip(payslip)
}

// ProcessEmployeePayrollWithAudit handles the payroll calculation for a single employee with audit trail
func (uc *PayrollUsecase) ProcessEmployeePayrollWithAudit(employeeID uint, req request.PayrollRequest, auditDB *middleware.AuditableDB) (*model.Payslip, error) {
	// Check if payslip already exists for this period
	exists, err := uc.payslipRepo.CheckPayslipExists(employeeID, req.PayPeriodStart, req.PayPeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing payslip: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("payslip already exists for this period")
	}

	// Get attendance records for the period
	attendances, err := uc.payslipRepo.GetAttendanceForPeriod(employeeID, req.PayPeriodStart, req.PayPeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance records: %v", err)
	}

	// Get overtime records for the period
	dateStart := req.PayPeriodStart.Format("2006-01-02")
	dateEnd := req.PayPeriodEnd.Format("2006-01-02")
	overtimes, err := uc.payslipRepo.GetOvertimeForPeriod(employeeID, dateStart, dateEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get overtime records: %v", err)
	}

	// Get approved reimbursements for the period
	reimbursements, err := uc.payslipRepo.GetApprovedReimbursementsForPeriod(employeeID, req.PayPeriodStart, req.PayPeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get reimbursement records: %v", err)
	}

	// Calculate totals
	attendanceDays := len(attendances)
	totalOvertimeHours := uc.calculateTotalOvertimeHours(overtimes)
	totalReimbursementAmount := uc.calculateTotalReimbursementAmount(reimbursements)

	// Calculate amounts
	overtimeAmount := float64(totalOvertimeHours) * req.OvertimeRate
	totalAmount := req.BasicSalary + overtimeAmount + totalReimbursementAmount

	// Create payslip with audit trail
	payslip := &model.Payslip{
		EmployeeID:          employeeID,
		PayPeriodStart:      req.PayPeriodStart,
		PayPeriodEnd:        req.PayPeriodEnd,
		BasicSalary:         req.BasicSalary,
		OvertimeHours:       totalOvertimeHours,
		OvertimeAmount:      overtimeAmount,
		ReimbursementAmount: totalReimbursementAmount,
		TotalAmount:         totalAmount,
		ProcessedAt:         time.Now(),
		Status:              "processed",
		AttendanceDays:      attendanceDays,
	}

	return uc.payslipRepo.CreatePayslipWithAudit(payslip, auditDB)
}

// ProcessAllEmployeesPayroll processes payroll for all active employees
func (uc *PayrollUsecase) ProcessAllEmployeesPayroll(req request.PayrollRequest) ([]model.Payslip, []string) {
	// Get all active employees
	employees, err := uc.employeeRepo.GetAllActiveEmployees()
	if err != nil {
		return nil, []string{fmt.Sprintf("Failed to get employees: %v", err)}
	}

	var processedPayslips []model.Payslip
	var errors []string

	for _, employee := range employees {
		payslip, err := uc.ProcessEmployeePayroll(employee.ID, req)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Employee %d: %s", employee.ID, err.Error()))
			continue
		}
		processedPayslips = append(processedPayslips, *payslip)
	}

	return processedPayslips, errors
}

// ProcessAllEmployeesPayrollWithAudit processes payroll for all active employees with audit trail
func (uc *PayrollUsecase) ProcessAllEmployeesPayrollWithAudit(req request.PayrollRequest, auditDB *middleware.AuditableDB) ([]model.Payslip, []string) {
	// Get all active employees
	employees, err := uc.employeeRepo.GetAllActiveEmployees()
	if err != nil {
		return nil, []string{fmt.Sprintf("Failed to get employees: %v", err)}
	}

	var processedPayslips []model.Payslip
	var errors []string

	for _, employee := range employees {
		payslip, err := uc.ProcessEmployeePayrollWithAudit(employee.ID, req, auditDB)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Employee %d: %s", employee.ID, err.Error()))
			continue
		}
		processedPayslips = append(processedPayslips, *payslip)
	}

	return processedPayslips, errors
}

// BuildDetailedPayslipResponse constructs the detailed payslip response
func (uc *PayrollUsecase) BuildDetailedPayslipResponse(payslip *model.Payslip, employee *model.Employee, attendances []model.Attendance, overtimes []model.Overtime, reimbursements []model.Reimbursement) map[string]interface{} {
	// Build attendance breakdown
	attendanceBreakdown := uc.buildAttendanceBreakdown(attendances)

	// Build overtime breakdown with calculated amounts
	overtimeBreakdown := uc.buildOvertimeBreakdown(overtimes, payslip)

	// Build reimbursement breakdown
	reimbursementBreakdown := uc.buildReimbursementBreakdown(reimbursements)

	// Build summary
	summary := map[string]interface{}{
		"basic_salary":          payslip.BasicSalary,
		"total_attendance_days": payslip.AttendanceDays,
		"total_overtime_hours":  payslip.OvertimeHours,
		"overtime_amount":       payslip.OvertimeAmount,
		"reimbursement_amount":  payslip.ReimbursementAmount,
		"total_take_home_pay":   payslip.TotalAmount,
	}

	return map[string]interface{}{
		"payslip_id":              payslip.ID,
		"employee_id":             payslip.EmployeeID,
		"employee_name":           employee.Name,
		"pay_period_start":        payslip.PayPeriodStart,
		"pay_period_end":          payslip.PayPeriodEnd,
		"processed_at":            payslip.ProcessedAt,
		"status":                  payslip.Status,
		"summary":                 summary,
		"attendance_breakdown":    attendanceBreakdown,
		"overtime_breakdown":      overtimeBreakdown,
		"reimbursement_breakdown": reimbursementBreakdown,
	}
}

// BuildPayrollSummary constructs the payroll summary response
func (uc *PayrollUsecase) BuildPayrollSummary(payslips []model.Payslip) map[string]interface{} {
	var employeeSummaries []map[string]interface{}
	var totalTakeHomePay float64
	var totalBasicSalary float64
	var totalOvertimeAmount float64
	var totalReimbursementAmount float64
	var totalAttendanceDays int
	var totalOvertimeHours int

	// Group payslips by employee to get employee totals
	employeePayslips := make(map[uint][]model.Payslip)
	employeeNames := make(map[uint]string)

	for _, payslip := range payslips {
		employeePayslips[payslip.EmployeeID] = append(employeePayslips[payslip.EmployeeID], payslip)

		// Get employee name if we don't have it yet
		if _, exists := employeeNames[payslip.EmployeeID]; !exists {
			employee, err := uc.payslipRepo.GetEmployeeByID(payslip.EmployeeID)
			if err == nil {
				employeeNames[payslip.EmployeeID] = employee.Name
			} else {
				employeeNames[payslip.EmployeeID] = "Unknown Employee"
			}
		}
	}

	// Calculate totals for each employee
	for employeeID, empPayslips := range employeePayslips {
		employeeSummary := uc.calculateEmployeeSummary(employeeID, empPayslips, employeeNames[employeeID])
		employeeSummaries = append(employeeSummaries, employeeSummary)

		// Add to overall totals
		totalTakeHomePay += employeeSummary["total_take_home_pay"].(float64)
		totalBasicSalary += employeeSummary["total_basic_salary"].(float64)
		totalOvertimeAmount += employeeSummary["total_overtime_amount"].(float64)
		totalReimbursementAmount += employeeSummary["total_reimbursement"].(float64)
		totalAttendanceDays += employeeSummary["total_attendance_days"].(int)
		totalOvertimeHours += employeeSummary["total_overtime_hours"].(int)
	}

	// Calculate averages
	summaryTotals := uc.calculateSummaryTotals(
		employeeSummaries,
		payslips,
		totalTakeHomePay,
		totalBasicSalary,
		totalOvertimeAmount,
		totalReimbursementAmount,
		totalAttendanceDays,
		totalOvertimeHours,
	)

	return map[string]interface{}{
		"summary_totals":     summaryTotals,
		"employee_summaries": employeeSummaries,
	}
}

// Helper functions for calculations and data building

func (uc *PayrollUsecase) calculateTotalOvertimeHours(overtimes []model.Overtime) int {
	totalOvertimeHours := 0
	for _, overtime := range overtimes {
		totalOvertimeHours += overtime.Hours
	}
	return totalOvertimeHours
}

func (uc *PayrollUsecase) calculateTotalReimbursementAmount(reimbursements []model.Reimbursement) float64 {
	totalReimbursementAmount := 0.0
	for _, reimbursement := range reimbursements {
		totalReimbursementAmount += reimbursement.Amount
	}
	return totalReimbursementAmount
}

func (uc *PayrollUsecase) buildAttendanceBreakdown(attendances []model.Attendance) []map[string]interface{} {
	var attendanceBreakdown []map[string]interface{}
	for _, attendance := range attendances {
		attendanceBreakdown = append(attendanceBreakdown, map[string]interface{}{
			"date":         attendance.Date,
			"check_in":     attendance.Checkin,
			"check_out":    attendance.Checkout,
			"hours_worked": attendance.HoursWorked,
			"status":       attendance.Status,
		})
	}
	return attendanceBreakdown
}

func (uc *PayrollUsecase) buildOvertimeBreakdown(overtimes []model.Overtime, payslip *model.Payslip) []map[string]interface{} {
	var overtimeBreakdown []map[string]interface{}
	overtimeRate := 0.0
	if payslip.OvertimeHours > 0 {
		overtimeRate = payslip.OvertimeAmount / float64(payslip.OvertimeHours)
	}

	for _, overtime := range overtimes {
		amount := float64(overtime.Hours) * overtimeRate
		overtimeBreakdown = append(overtimeBreakdown, map[string]interface{}{
			"date":   overtime.OvertimeDate,
			"hours":  overtime.Hours,
			"rate":   overtimeRate,
			"amount": amount,
			"reason": overtime.Reason,
		})
	}
	return overtimeBreakdown
}

func (uc *PayrollUsecase) buildReimbursementBreakdown(reimbursements []model.Reimbursement) []map[string]interface{} {
	var reimbursementBreakdown []map[string]interface{}
	for _, reimbursement := range reimbursements {
		reimbursementBreakdown = append(reimbursementBreakdown, map[string]interface{}{
			"date":   reimbursement.ReimbursementDate,
			"amount": reimbursement.Amount,
			"reason": reimbursement.Reason,
			"status": reimbursement.Status,
		})
	}
	return reimbursementBreakdown
}

func (uc *PayrollUsecase) calculateEmployeeSummary(employeeID uint, empPayslips []model.Payslip, employeeName string) map[string]interface{} {
	var empTotalTakeHome float64
	var empTotalBasic float64
	var empTotalOvertime float64
	var empTotalReimbursement float64
	var empTotalAttendanceDays int
	var empTotalOvertimeHours int
	var payslipCount int

	for _, payslip := range empPayslips {
		empTotalTakeHome += payslip.TotalAmount
		empTotalBasic += payslip.BasicSalary
		empTotalOvertime += payslip.OvertimeAmount
		empTotalReimbursement += payslip.ReimbursementAmount
		empTotalAttendanceDays += payslip.AttendanceDays
		empTotalOvertimeHours += payslip.OvertimeHours
		payslipCount++
	}

	return map[string]interface{}{
		"employee_id":           employeeID,
		"employee_name":         employeeName,
		"payslip_count":         payslipCount,
		"total_take_home_pay":   empTotalTakeHome,
		"total_basic_salary":    empTotalBasic,
		"total_overtime_amount": empTotalOvertime,
		"total_reimbursement":   empTotalReimbursement,
		"total_attendance_days": empTotalAttendanceDays,
		"total_overtime_hours":  empTotalOvertimeHours,
	}
}

func (uc *PayrollUsecase) calculateSummaryTotals(employeeSummaries []map[string]interface{}, payslips []model.Payslip, totalTakeHomePay, totalBasicSalary, totalOvertimeAmount, totalReimbursementAmount float64, totalAttendanceDays, totalOvertimeHours int) map[string]interface{} {
	employeeCount := len(employeeSummaries)
	avgTakeHomePay := 0.0
	avgBasicSalary := 0.0
	avgOvertimeAmount := 0.0
	avgReimbursementAmount := 0.0

	if employeeCount > 0 {
		avgTakeHomePay = totalTakeHomePay / float64(employeeCount)
		avgBasicSalary = totalBasicSalary / float64(employeeCount)
		avgOvertimeAmount = totalOvertimeAmount / float64(employeeCount)
		avgReimbursementAmount = totalReimbursementAmount / float64(employeeCount)
	}

	return map[string]interface{}{
		"total_employees":            employeeCount,
		"total_payslips":             len(payslips),
		"total_take_home_pay":        totalTakeHomePay,
		"total_basic_salary":         totalBasicSalary,
		"total_overtime_amount":      totalOvertimeAmount,
		"total_reimbursement_amount": totalReimbursementAmount,
		"total_attendance_days":      totalAttendanceDays,
		"total_overtime_hours":       totalOvertimeHours,
		"average_take_home_pay":      avgTakeHomePay,
		"average_basic_salary":       avgBasicSalary,
		"average_overtime_amount":    avgOvertimeAmount,
		"average_reimbursement":      avgReimbursementAmount,
	}
}
