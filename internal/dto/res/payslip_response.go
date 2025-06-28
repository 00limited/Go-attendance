package res

import "time"

// AttendanceDetail represents attendance breakdown in payslip
type AttendanceDetail struct {
	Date        string `json:"date"`
	CheckIn     string `json:"check_in"`
	CheckOut    string `json:"check_out"`
	HoursWorked int    `json:"hours_worked"`
	Status      string `json:"status"`
}

// OvertimeDetail represents overtime breakdown in payslip
type OvertimeDetail struct {
	Date   string  `json:"date"`
	Hours  int     `json:"hours"`
	Rate   float64 `json:"rate"`
	Amount float64 `json:"amount"`
	Reason string  `json:"reason"`
}

// ReimbursementDetail represents reimbursement breakdown in payslip
type ReimbursementDetail struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Reason string  `json:"reason"`
	Status string  `json:"status"`
}

// PayslipSummary represents the salary calculation summary
type PayslipSummary struct {
	BasicSalary         float64 `json:"basic_salary"`
	TotalAttendanceDays int     `json:"total_attendance_days"`
	TotalOvertimeHours  int     `json:"total_overtime_hours"`
	OvertimeAmount      float64 `json:"overtime_amount"`
	ReimbursementAmount float64 `json:"reimbursement_amount"`
	TotalTakeHomePay    float64 `json:"total_take_home_pay"`
}

// DetailedPayslipResponse represents the complete payslip with breakdowns
type DetailedPayslipResponse struct {
	PayslipID              uint                  `json:"payslip_id"`
	EmployeeID             uint                  `json:"employee_id"`
	EmployeeName           string                `json:"employee_name"`
	PayPeriodStart         time.Time             `json:"pay_period_start"`
	PayPeriodEnd           time.Time             `json:"pay_period_end"`
	ProcessedAt            time.Time             `json:"processed_at"`
	Status                 string                `json:"status"`
	Summary                PayslipSummary        `json:"summary"`
	AttendanceBreakdown    []AttendanceDetail    `json:"attendance_breakdown"`
	OvertimeBreakdown      []OvertimeDetail      `json:"overtime_breakdown"`
	ReimbursementBreakdown []ReimbursementDetail `json:"reimbursement_breakdown"`
}

// PayslipListResponse represents a list of payslips for an employee
type PayslipListResponse struct {
	PayslipID      uint      `json:"payslip_id"`
	PayPeriodStart time.Time `json:"pay_period_start"`
	PayPeriodEnd   time.Time `json:"pay_period_end"`
	TotalAmount    float64   `json:"total_amount"`
	Status         string    `json:"status"`
	ProcessedAt    time.Time `json:"processed_at"`
}
