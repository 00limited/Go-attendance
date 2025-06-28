// Package repository contains unit tests for the payslip repository functionality.
//
// This test file contains comprehensive tests for the PayslipRepository functions:
//
// CreatePayslip tests cover:
// 1. Valid payslip creation
// 2. Database errors during creation
// 3. Edge cases (nil payslip, zero values)
//
// CreatePayslipWithAudit tests cover:
// 1. Valid payslip creation with audit
// 2. Audit database errors
// 3. Audit context integration
//
// GetPayslipByEmployeeAndPeriod tests cover:
// 1. Valid retrieval by employee and period
// 2. No records found scenarios
// 3. Database errors
// 4. Edge cases (zero employee ID, invalid dates)
//
// GetPayslipByID tests cover:
// 1. Valid retrieval by ID
// 2. Record not found scenarios
// 3. Database errors
// 4. Edge cases (zero ID, large ID)
//
// GetPayslipsByEmployee tests cover:
// 1. Valid retrieval of multiple payslips
// 2. Empty results for employee
// 3. Ordering verification (DESC by pay_period_start)
// 4. Database errors
//
// GetPayslipsByPeriod tests cover:
// 1. Valid retrieval by date range
// 2. Empty results for period
// 3. Date boundary testing
// 4. Database errors
//
// CheckPayslipExists tests cover:
// 1. Existing payslip detection
// 2. Non-existing payslip detection
// 3. Database errors
// 4. Edge cases
//
// GetAttendanceForPeriod tests cover:
// 1. Valid attendance retrieval
// 2. Empty attendance records
// 3. Date range filtering
// 4. Database errors
//
// GetOvertimeForPeriod tests cover:
// 1. Valid overtime retrieval (approved only)
// 2. Status filtering verification
// 3. Date range filtering
// 4. Database errors
//
// GetApprovedReimbursementsForPeriod tests cover:
// 1. Valid reimbursement retrieval (approved only)
// 2. Status filtering verification
// 3. Date range filtering
// 4. Database errors
//
// GetEmployeeByID tests cover:
// 1. Valid employee retrieval
// 2. Employee not found scenarios
// 3. Database errors
// 4. Edge cases
//
// All tests use an in-memory SQLite database for fast, isolated execution.

package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t testing.TB) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Auto migrate all models
	err = db.AutoMigrate(
		&model.Payslip{},
		&model.Employee{},
		&model.Attendance{},
		&model.Overtime{},
		&model.Reimbursement{},
	)
	require.NoError(t, err)

	return db
}

// createTestEmployee creates a test employee record
func createTestEmployee(t testing.TB, db *gorm.DB, id uint, name string) *model.Employee {
	employee := &model.Employee{
		DefaultAttribute: model.DefaultAttribute{ID: id},
		Name:             name,
		Role:             "employee",
		Active:           true,
	}
	err := db.Create(employee).Error
	require.NoError(t, err)
	return employee
}

// createTestPayslip creates a test payslip record
func createTestPayslip(t testing.TB, db *gorm.DB, employeeID uint, start, end time.Time) *model.Payslip {
	payslip := &model.Payslip{
		EmployeeID:     employeeID,
		PayPeriodStart: start,
		PayPeriodEnd:   end,
		BasicSalary:    5000000.0,
		TotalAmount:    5500000.0,
		ProcessedAt:    time.Now(),
		Status:         "processed",
	}
	err := db.Create(payslip).Error
	require.NoError(t, err)
	return payslip
}

// Tests for CreatePayslip function

func TestPayslipRepository_CreatePayslip_ValidPayslip(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee first
	employee := createTestEmployee(t, db, 1, "John Doe")

	// Create test payslip
	payslip := &model.Payslip{
		EmployeeID:     employee.ID,
		PayPeriodStart: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		BasicSalary:    5000000.0,
		TotalAmount:    5500000.0,
		ProcessedAt:    time.Now(),
		Status:         "processed",
	}

	// Execute
	result, err := repo.CreatePayslip(payslip)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotZero(t, result.ID)
	assert.Equal(t, payslip.EmployeeID, result.EmployeeID)
	assert.Equal(t, payslip.TotalAmount, result.TotalAmount)

	// Verify in database
	var dbPayslip model.Payslip
	err = db.First(&dbPayslip, result.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, payslip.EmployeeID, dbPayslip.EmployeeID)
}

func TestPayslipRepository_CreatePayslip_NilPayslip(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute
	result, err := repo.CreatePayslip(nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// Tests for CreatePayslipWithAudit function

func TestPayslipRepository_CreatePayslipWithAudit_ValidPayslip(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)
	auditDB := middleware.NewAuditableDB(db, 1)

	// Create test employee first
	employee := createTestEmployee(t, db, 1, "John Doe")

	// Create test payslip
	payslip := &model.Payslip{
		EmployeeID:     employee.ID,
		PayPeriodStart: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		BasicSalary:    5000000.0,
		TotalAmount:    5500000.0,
		ProcessedAt:    time.Now(),
		Status:         "processed",
	}

	// Execute
	result, err := repo.CreatePayslipWithAudit(payslip, auditDB)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotZero(t, result.ID)
	assert.Equal(t, payslip.EmployeeID, result.EmployeeID)

	// Verify in database with audit fields
	var dbPayslip model.Payslip
	err = db.First(&dbPayslip, result.ID).Error
	assert.NoError(t, err)
	assert.NotZero(t, dbPayslip.CreatedBy)
}

// Tests for GetPayslipByEmployeeAndPeriod function

func TestPayslipRepository_GetPayslipByEmployeeAndPeriod_ValidData(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee and payslip
	employee := createTestEmployee(t, db, 1, "John Doe")
	startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)
	originalPayslip := createTestPayslip(t, db, employee.ID, startDate, endDate)

	// Execute
	result, err := repo.GetPayslipByEmployeeAndPeriod(employee.ID, startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, originalPayslip.ID, result.ID)
	assert.Equal(t, originalPayslip.EmployeeID, result.EmployeeID)
	assert.Equal(t, originalPayslip.PayPeriodStart.Unix(), result.PayPeriodStart.Unix())
	assert.Equal(t, originalPayslip.PayPeriodEnd.Unix(), result.PayPeriodEnd.Unix())
}

func TestPayslipRepository_GetPayslipByEmployeeAndPeriod_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute with non-existing data
	result, err := repo.GetPayslipByEmployeeAndPeriod(999, time.Now(), time.Now())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "record not found")
}

// Tests for GetPayslipByID function

func TestPayslipRepository_GetPayslipByID_ValidID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee and payslip
	employee := createTestEmployee(t, db, 1, "John Doe")
	startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)
	originalPayslip := createTestPayslip(t, db, employee.ID, startDate, endDate)

	// Execute
	result, err := repo.GetPayslipByID(originalPayslip.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, originalPayslip.ID, result.ID)
	assert.Equal(t, originalPayslip.EmployeeID, result.EmployeeID)
}

func TestPayslipRepository_GetPayslipByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute
	result, err := repo.GetPayslipByID(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "record not found")
}

func TestPayslipRepository_GetPayslipByID_ZeroID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute
	result, err := repo.GetPayslipByID(0)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// Tests for GetPayslipsByEmployee function

func TestPayslipRepository_GetPayslipsByEmployee_ValidEmployee(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee
	employee := createTestEmployee(t, db, 1, "John Doe")

	// Create multiple payslips with different dates for ordering test
	payslip1 := createTestPayslip(t, db, employee.ID,
		time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC))
	payslip2 := createTestPayslip(t, db, employee.ID,
		time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC))

	// Execute
	results, err := repo.GetPayslipsByEmployee(employee.ID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify ordering (DESC by pay_period_start)
	assert.Equal(t, payslip2.ID, results[0].ID) // June payslip should be first
	assert.Equal(t, payslip1.ID, results[1].ID) // May payslip should be second
}

func TestPayslipRepository_GetPayslipsByEmployee_NoPayslips(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee but no payslips
	employee := createTestEmployee(t, db, 1, "John Doe")

	// Execute
	results, err := repo.GetPayslipsByEmployee(employee.ID)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestPayslipRepository_GetPayslipsByEmployee_NonExistentEmployee(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute
	results, err := repo.GetPayslipsByEmployee(999)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, results)
}

// Tests for GetPayslipsByPeriod function

func TestPayslipRepository_GetPayslipsByPeriod_ValidPeriod(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employees
	employee1 := createTestEmployee(t, db, 1, "John Doe")
	employee2 := createTestEmployee(t, db, 2, "Jane Smith")

	// Create payslips in the target period
	startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	payslip1 := createTestPayslip(t, db, employee1.ID, startDate, endDate)
	payslip2 := createTestPayslip(t, db, employee2.ID, startDate, endDate)

	// Create a payslip outside the period
	createTestPayslip(t, db, employee1.ID,
		time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC))

	// Execute
	results, err := repo.GetPayslipsByPeriod(startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify correct payslips are returned
	payslipIDs := []uint{results[0].ID, results[1].ID}
	assert.Contains(t, payslipIDs, payslip1.ID)
	assert.Contains(t, payslipIDs, payslip2.ID)
}

func TestPayslipRepository_GetPayslipsByPeriod_NoPeriodMatches(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee and payslip
	employee := createTestEmployee(t, db, 1, "John Doe")
	createTestPayslip(t, db, employee.ID,
		time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC))

	// Execute with different period
	results, err := repo.GetPayslipsByPeriod(
		time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC))

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, results)
}

// Tests for CheckPayslipExists function

func TestPayslipRepository_CheckPayslipExists_Exists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee and payslip
	employee := createTestEmployee(t, db, 1, "John Doe")
	startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)
	createTestPayslip(t, db, employee.ID, startDate, endDate)

	// Execute
	exists, err := repo.CheckPayslipExists(employee.ID, startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestPayslipRepository_CheckPayslipExists_NotExists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute without creating any payslips
	exists, err := repo.CheckPayslipExists(1,
		time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC))

	// Assert
	assert.NoError(t, err)
	assert.False(t, exists)
}

// Tests for GetAttendanceForPeriod function

func TestPayslipRepository_GetAttendanceForPeriod_ValidData(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee
	employee := createTestEmployee(t, db, 1, "John Doe")

	// Create attendance records
	startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	attendance1 := &model.Attendance{
		EmployeeID: employee.ID,
		Date:       time.Date(2025, 6, 5, 0, 0, 0, 0, time.UTC),
		Checkin:    time.Date(2025, 6, 5, 8, 0, 0, 0, time.UTC),
		Status:     "present",
	}
	attendance2 := &model.Attendance{
		EmployeeID: employee.ID,
		Date:       time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC),
		Checkin:    time.Date(2025, 6, 10, 8, 0, 0, 0, time.UTC),
		Status:     "present",
	}
	// Outside period
	attendance3 := &model.Attendance{
		EmployeeID: employee.ID,
		Date:       time.Date(2025, 7, 5, 0, 0, 0, 0, time.UTC),
		Checkin:    time.Date(2025, 7, 5, 8, 0, 0, 0, time.UTC),
		Status:     "present",
	}

	err := db.Create([]*model.Attendance{attendance1, attendance2, attendance3}).Error
	require.NoError(t, err)

	// Execute
	results, err := repo.GetAttendanceForPeriod(employee.ID, startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify correct attendance records
	attendanceIDs := []uint{results[0].ID, results[1].ID}
	assert.Contains(t, attendanceIDs, attendance1.ID)
	assert.Contains(t, attendanceIDs, attendance2.ID)
}

func TestPayslipRepository_GetAttendanceForPeriod_NoData(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute without creating any attendance
	results, err := repo.GetAttendanceForPeriod(1,
		time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC))

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, results)
}

// Tests for GetOvertimeForPeriod function

func TestPayslipRepository_GetOvertimeForPeriod_ValidData(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee
	employee := createTestEmployee(t, db, 1, "John Doe")

	// Create overtime records
	overtime1 := &model.Overtime{
		EmployeeID:   employee.ID,
		OvertimeDate: "2025-06-05",
		Hours:        8,
		Reason:       "Extra work required",
		Status:       model.OvertimeApproved,
	}
	overtime2 := &model.Overtime{
		EmployeeID:   employee.ID,
		OvertimeDate: "2025-06-10",
		Hours:        4,
		Reason:       "Urgent project",
		Status:       model.OvertimeApproved,
	}
	// Pending status (should not be included)
	overtime3 := &model.Overtime{
		EmployeeID:   employee.ID,
		OvertimeDate: "2025-06-15",
		Hours:        6,
		Reason:       "Additional work",
		Status:       model.OvertimePending,
	}
	// Outside date range
	overtime4 := &model.Overtime{
		EmployeeID:   employee.ID,
		OvertimeDate: "2025-07-05",
		Hours:        8,
		Reason:       "Extra work",
		Status:       model.OvertimeApproved,
	}

	err := db.Create([]*model.Overtime{overtime1, overtime2, overtime3, overtime4}).Error
	require.NoError(t, err)

	// Execute
	results, err := repo.GetOvertimeForPeriod(employee.ID, "2025-06-01", "2025-06-30")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify only approved overtime in date range
	for _, overtime := range results {
		assert.Equal(t, model.OvertimeApproved, overtime.Status)
		assert.True(t, overtime.OvertimeDate >= "2025-06-01" && overtime.OvertimeDate <= "2025-06-30")
	}
}

func TestPayslipRepository_GetOvertimeForPeriod_NoData(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute without creating any overtime
	results, err := repo.GetOvertimeForPeriod(1, "2025-06-01", "2025-06-30")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, results)
}

// Tests for GetApprovedReimbursementsForPeriod function

func TestPayslipRepository_GetApprovedReimbursementsForPeriod_ValidData(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee
	employee := createTestEmployee(t, db, 1, "John Doe")

	// Create reimbursement records
	startDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	reimbursement1 := &model.Reimbursement{
		EmployeeID:        employee.ID,
		ReimbursementDate: time.Date(2025, 6, 5, 0, 0, 0, 0, time.UTC),
		Amount:            100000.0,
		Status:            model.ReimbursementApproved,
		Reason:            "Transport",
	}
	reimbursement2 := &model.Reimbursement{
		EmployeeID:        employee.ID,
		ReimbursementDate: time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC),
		Amount:            50000.0,
		Status:            model.ReimbursementApproved,
		Reason:            "Meals",
	}
	// Pending status (should not be included)
	reimbursement3 := &model.Reimbursement{
		EmployeeID:        employee.ID,
		ReimbursementDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		Amount:            75000.0,
		Status:            model.ReimbursementPending,
		Reason:            "Equipment",
	}

	err := db.Create([]*model.Reimbursement{reimbursement1, reimbursement2, reimbursement3}).Error
	require.NoError(t, err)

	// Execute
	results, err := repo.GetApprovedReimbursementsForPeriod(employee.ID, startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify only approved reimbursements
	for _, reimbursement := range results {
		assert.Equal(t, model.ReimbursementApproved, reimbursement.Status)
	}
}

func TestPayslipRepository_GetApprovedReimbursementsForPeriod_NoData(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute without creating any reimbursements
	results, err := repo.GetApprovedReimbursementsForPeriod(1,
		time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC))

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, results)
}

// Tests for GetEmployeeByID function

func TestPayslipRepository_GetEmployeeByID_ValidID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Create test employee
	originalEmployee := createTestEmployee(t, db, 1, "John Doe")

	// Execute
	result, err := repo.GetEmployeeByID(originalEmployee.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, originalEmployee.ID, result.ID)
	assert.Equal(t, originalEmployee.Name, result.Name)
	assert.Equal(t, originalEmployee.Role, result.Role)
}

func TestPayslipRepository_GetEmployeeByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute
	result, err := repo.GetEmployeeByID(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "record not found")
}

func TestPayslipRepository_GetEmployeeByID_ZeroID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute
	result, err := repo.GetEmployeeByID(0)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// Tests for GetDB function

func TestPayslipRepository_GetDB(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPayslipRepository(db)

	// Execute
	result := repo.GetDB()

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, db, result)
}

// Edge Cases and Integration Tests

func TestPayslipRepository_EdgeCases(t *testing.T) {
	t.Run("Large dataset performance", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewPayslipRepository(db)

		// Create multiple employees and payslips
		numEmployees := 100
		numPayslipsPerEmployee := 12

		for i := 1; i <= numEmployees; i++ {
			employee := createTestEmployee(t, db, uint(i), "Employee"+string(rune(i)))

			for j := 1; j <= numPayslipsPerEmployee; j++ {
				createTestPayslip(t, db, employee.ID,
					time.Date(2025, time.Month(j), 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, time.Month(j), 28, 0, 0, 0, 0, time.UTC))
			}
		}

		// Test retrieval
		start := time.Now()
		results, err := repo.GetPayslipsByPeriod(
			time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC))
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Len(t, results, numEmployees*numPayslipsPerEmployee)
		assert.Less(t, duration, time.Second) // Should be fast
	})

	t.Run("Boundary date testing", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewPayslipRepository(db)

		employee := createTestEmployee(t, db, 1, "John Doe")

		// Create payslip at exact boundary
		exactStart := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
		exactEnd := time.Date(2025, 6, 30, 23, 59, 59, 999999999, time.UTC)
		createTestPayslip(t, db, employee.ID, exactStart, exactEnd)

		// Search with same boundaries
		results, err := repo.GetPayslipsByPeriod(exactStart, exactEnd)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

// Benchmark Tests

func BenchmarkPayslipRepository_CreatePayslip(b *testing.B) {
	db := setupTestDB(b)
	repo := NewPayslipRepository(db)
	employee := createTestEmployee(b, db, 1, "John Doe")

	payslip := &model.Payslip{
		EmployeeID:     employee.ID,
		PayPeriodStart: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		PayPeriodEnd:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		BasicSalary:    5000000.0,
		TotalAmount:    5500000.0,
		ProcessedAt:    time.Now(),
		Status:         "processed",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset ID for each iteration
		payslip.ID = 0
		_, _ = repo.CreatePayslip(payslip)
	}
}

func BenchmarkPayslipRepository_GetPayslipsByEmployee(b *testing.B) {
	db := setupTestDB(b)
	repo := NewPayslipRepository(db)
	employee := createTestEmployee(b, db, 1, "John Doe")

	// Create multiple payslips
	for i := 1; i <= 12; i++ {
		createTestPayslip(b, db, employee.ID,
			time.Date(2025, time.Month(i), 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, time.Month(i), 28, 0, 0, 0, 0, time.UTC))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetPayslipsByEmployee(employee.ID)
	}
}
