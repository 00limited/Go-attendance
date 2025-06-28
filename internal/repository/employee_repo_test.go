package repository

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupEmployeeTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&model.Employee{})
	if err != nil {
		panic("failed to migrate: " + err.Error())
	}

	return db
}

func createEmployeeTestData() request.CreateEmployeeRequest {
	return request.CreateEmployeeRequest{
		Name:     "John Doe",
		Password: "password123",
		Role:     "admin",
		Active:   true,
	}
}

func createUpdateEmployeeTestData() request.UpdateEmployeeRequest {
	return request.UpdateEmployeeRequest{
		Name:     "Jane Doe",
		Password: "newpassword123",
		Role:     "admin",
		Active:   false,
	}
}

func seedEmployeesForTest(db *gorm.DB, count int) []model.Employee {
	employees := make([]model.Employee, count)
	for i := 0; i < count; i++ {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"+strconv.Itoa(i)), bcrypt.DefaultCost)
		employee := model.Employee{
			Name:     "Employee " + strconv.Itoa(i+1),
			Password: string(hashedPassword),
			Role:     "admin",
			Active:   i%2 == 0, // Alternate between active and inactive
		}
		db.Create(&employee)
		employees[i] = employee
	}
	return employees
}

// Test CreateEmployee - Valid Case
func TestCreateEmployee_Success(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	req := createEmployeeTestData()
	employee, err := repo.CreateEmployee(req)

	assert.NoError(t, err)
	assert.NotNil(t, employee)
	assert.Equal(t, req.Name, employee.Name)
	assert.Equal(t, req.Role, employee.Role)
	assert.Equal(t, req.Active, employee.Active)
	assert.NotEmpty(t, employee.Password)
	assert.NotEqual(t, req.Password, employee.Password) // Password should be hashed

	// Verify password is hashed correctly
	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(req.Password))
	assert.NoError(t, err)

	// Verify employee is saved in database
	var savedEmployee model.Employee
	err = db.First(&savedEmployee, employee.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, employee.Name, savedEmployee.Name)
}

// Test CreateEmployee - Password Hashing Error
func TestCreateEmployee_PasswordHashingError(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	req := request.CreateEmployeeRequest{
		Name:     "Test User",
		Password: "", // Empty password might cause issues
		Role:     "admin",
		Active:   true,
	}

	employee, err := repo.CreateEmployee(req)

	// Should succeed even with empty password
	assert.NoError(t, err)
	assert.NotNil(t, employee)
}

// Test CreateEmployee - Database Error
func TestCreateEmployee_DatabaseError(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Drop the table to simulate database error
	db.Migrator().DropTable(&model.Employee{})

	req := createEmployeeTestData()
	employee, err := repo.CreateEmployee(req)

	assert.Error(t, err)
	assert.Nil(t, employee)
}

// Test CreateEmployee - Duplicate Name
func TestCreateEmployee_DuplicateName(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Add unique constraint on name for this test
	db.Exec("CREATE UNIQUE INDEX idx_employees_name ON employees(name)")

	req := createEmployeeTestData()

	// Create first employee
	employee1, err := repo.CreateEmployee(req)
	assert.NoError(t, err)
	assert.NotNil(t, employee1)

	// Try to create second employee with same name
	employee2, err := repo.CreateEmployee(req)
	assert.Error(t, err)
	assert.Nil(t, employee2)
}

// Test GetAllEmployees - Success
func TestGetAllEmployees_Success(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Seed some employees
	seedEmployeesForTest(db, 5)

	employees, err := repo.GetAllEmployees()

	assert.NoError(t, err)
	assert.Len(t, employees, 5)
}

// Test GetAllEmployees - Empty Database
func TestGetAllEmployees_EmptyDatabase(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	employees, err := repo.GetAllEmployees()

	assert.NoError(t, err)
	assert.Len(t, employees, 0)
}

// Test GetAllEmployees - Database Error
func TestGetAllEmployees_DatabaseError(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Drop the table to simulate database error
	db.Migrator().DropTable(&model.Employee{})

	employees, err := repo.GetAllEmployees()

	assert.Error(t, err)
	assert.Nil(t, employees)
}

// Test GetAllActiveEmployees - Success
func TestGetAllActiveEmployees_Success(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Seed employees (alternating active/inactive)
	seedEmployeesForTest(db, 6)

	employees, err := repo.GetAllActiveEmployees()

	assert.NoError(t, err)
	assert.Len(t, employees, 3) // Half should be active
	for _, emp := range employees {
		assert.True(t, emp.Active)
	}
}

// Test GetAllActiveEmployees - No Active Employees
func TestGetAllActiveEmployees_NoActiveEmployees(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create only inactive employees
	for i := 0; i < 3; i++ {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		employee := model.Employee{
			Name:     "Inactive Employee " + strconv.Itoa(i+1),
			Password: string(hashedPassword),
			Role:     "admin",
			Active:   false,
		}
		db.Create(&employee)
	}

	employees, err := repo.GetAllActiveEmployees()

	assert.NoError(t, err)
	assert.Len(t, employees, 0)
}

// Test GetAllActiveEmployees - Database Error
func TestGetAllActiveEmployees_DatabaseError(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Drop the table to simulate database error
	db.Migrator().DropTable(&model.Employee{})

	employees, err := repo.GetAllActiveEmployees()

	assert.Error(t, err)
	assert.Nil(t, employees)
}

// Test UpdateEmployee - Success
func TestUpdateEmployee_Success(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create initial employee
	createReq := createEmployeeTestData()
	createdEmployee, err := repo.CreateEmployee(createReq)
	require.NoError(t, err)

	// Update employee
	updateReq := createUpdateEmployeeTestData()
	updatedEmployee, err := repo.UpdateEmployee(strconv.Itoa(int(createdEmployee.ID)), updateReq)

	assert.NoError(t, err)
	assert.NotNil(t, updatedEmployee)
	assert.Equal(t, updateReq.Name, updatedEmployee.Name)
	assert.Equal(t, updateReq.Role, updatedEmployee.Role)
	assert.Equal(t, updateReq.Active, updatedEmployee.Active)
	assert.NotEqual(t, updateReq.Password, updatedEmployee.Password) // Password should be hashed

	// Verify password is hashed correctly
	err = bcrypt.CompareHashAndPassword([]byte(updatedEmployee.Password), []byte(updateReq.Password))
	assert.NoError(t, err)
}

// Test UpdateEmployee - Employee Not Found
func TestUpdateEmployee_EmployeeNotFound(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	updateReq := createUpdateEmployeeTestData()
	updatedEmployee, err := repo.UpdateEmployee("999", updateReq)

	assert.Error(t, err)
	assert.Nil(t, updatedEmployee)
	assert.Contains(t, err.Error(), "record not found")
}

// Test UpdateEmployee - Invalid Employee ID
func TestUpdateEmployee_InvalidEmployeeID(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	updateReq := createUpdateEmployeeTestData()
	updatedEmployee, err := repo.UpdateEmployee("invalid", updateReq)

	assert.Error(t, err)
	assert.Nil(t, updatedEmployee)
}

// Test UpdateEmployee - Password Hashing Error
func TestUpdateEmployee_PasswordHashingError(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create initial employee
	createReq := createEmployeeTestData()
	createdEmployee, err := repo.CreateEmployee(createReq)
	require.NoError(t, err)

	// Update with empty password (should still work)
	updateReq := request.UpdateEmployeeRequest{
		Name:     "Updated Name",
		Password: "",
		Role:     "admin",
		Active:   false,
	}

	updatedEmployee, err := repo.UpdateEmployee(strconv.Itoa(int(createdEmployee.ID)), updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updatedEmployee)
}

// Test DeleteEmployee - Success
func TestDeleteEmployee_Success(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create employee
	createReq := createEmployeeTestData()
	createdEmployee, err := repo.CreateEmployee(createReq)
	require.NoError(t, err)

	// Delete employee
	err = repo.DeleteEmployee(strconv.Itoa(int(createdEmployee.ID)))

	assert.NoError(t, err)

	// Verify employee is soft deleted
	var deletedEmployee model.Employee
	err = db.Unscoped().First(&deletedEmployee, createdEmployee.ID).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedEmployee.DeletedAt)
}

// Test DeleteEmployee - Employee Not Found
func TestDeleteEmployee_EmployeeNotFound(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	err := repo.DeleteEmployee("999")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

// Test DeleteEmployee - Invalid Employee ID
func TestDeleteEmployee_InvalidEmployeeID(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	err := repo.DeleteEmployee("invalid")

	assert.Error(t, err)
}

// Test GetEmployeeByID - Success
func TestGetEmployeeByID_Success(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create employee
	createReq := createEmployeeTestData()
	createdEmployee, err := repo.CreateEmployee(createReq)
	require.NoError(t, err)

	// Get employee by ID
	foundEmployee, err := repo.GetEmployeeByID(createdEmployee.ID)

	assert.NoError(t, err)
	assert.NotNil(t, foundEmployee)
	assert.Equal(t, createdEmployee.ID, foundEmployee.ID)
	assert.Equal(t, createdEmployee.Name, foundEmployee.Name)
	assert.Equal(t, createdEmployee.Role, foundEmployee.Role)
	assert.Equal(t, createdEmployee.Active, foundEmployee.Active)
}

// Test GetEmployeeByID - Employee Not Found
func TestGetEmployeeByID_EmployeeNotFound(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	foundEmployee, err := repo.GetEmployeeByID(999)

	assert.Error(t, err)
	assert.Nil(t, foundEmployee)
	assert.Contains(t, err.Error(), "record not found")
}

// Test GetEmployeeByID - Database Error
func TestGetEmployeeByID_DatabaseError(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Drop the table to simulate database error
	db.Migrator().DropTable(&model.Employee{})

	foundEmployee, err := repo.GetEmployeeByID(1)

	assert.Error(t, err)
	assert.Nil(t, foundEmployee)
}

// Test GetEmployeeByName - Success
func TestGetEmployeeByName_Success(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create employee
	createReq := createEmployeeTestData()
	createdEmployee, err := repo.CreateEmployee(createReq)
	require.NoError(t, err)

	// Get employee by name
	foundEmployee, err := repo.GetEmployeeByName(createdEmployee.Name)

	assert.NoError(t, err)
	assert.NotNil(t, foundEmployee)
	assert.Equal(t, createdEmployee.ID, foundEmployee.ID)
	assert.Equal(t, createdEmployee.Name, foundEmployee.Name)
	assert.Equal(t, createdEmployee.Role, foundEmployee.Role)
	assert.Equal(t, createdEmployee.Active, foundEmployee.Active)
}

// Test GetEmployeeByName - Employee Not Found
func TestGetEmployeeByName_EmployeeNotFound(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	foundEmployee, err := repo.GetEmployeeByName("NonExistent User")

	assert.Error(t, err)
	assert.Nil(t, foundEmployee)
	assert.Contains(t, err.Error(), "record not found")
}

// Test GetEmployeeByName - Empty Name
func TestGetEmployeeByName_EmptyName(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	foundEmployee, err := repo.GetEmployeeByName("")

	assert.Error(t, err)
	assert.Nil(t, foundEmployee)
}

// Test GetEmployeeByName - Database Error
func TestGetEmployeeByName_DatabaseError(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Drop the table to simulate database error
	db.Migrator().DropTable(&model.Employee{})

	foundEmployee, err := repo.GetEmployeeByName("Test User")

	assert.Error(t, err)
	assert.Nil(t, foundEmployee)
}

// Test CreateEmployeeWithAudit - Success
func TestCreateEmployeeWithAudit_Success(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	auditDB := middleware.NewAuditableDB(db, 1)
	req := createEmployeeTestData()
	employee, err := repo.CreateEmployeeWithAudit(req, auditDB)

	assert.NoError(t, err)
	assert.NotNil(t, employee)
	assert.Equal(t, req.Name, employee.Name)
	assert.Equal(t, req.Role, employee.Role)
	assert.Equal(t, req.Active, employee.Active)
	assert.NotEmpty(t, employee.Password)
	assert.NotEqual(t, req.Password, employee.Password) // Password should be hashed

	// Verify audit fields are set
	assert.NotNil(t, employee.CreatedBy)
	assert.Equal(t, uint(1), *employee.CreatedBy)

	// Verify password is hashed correctly
	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(req.Password))
	assert.NoError(t, err)
}

// Test CreateEmployeeWithAudit - Password Hashing Error
func TestCreateEmployeeWithAudit_PasswordHashingError(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	auditDB := middleware.NewAuditableDB(db, 1)
	req := request.CreateEmployeeRequest{
		Name:     "Test User",
		Password: "", // Empty password might cause issues
		Role:     "admin",
		Active:   true,
	}

	employee, err := repo.CreateEmployeeWithAudit(req, auditDB)

	// Should succeed even with empty password
	assert.NoError(t, err)
	assert.NotNil(t, employee)
}

// Test CreateEmployeeWithAudit - Database Error
func TestCreateEmployeeWithAudit_DatabaseError(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Drop the table to simulate database error
	db.Migrator().DropTable(&model.Employee{})

	auditDB := middleware.NewAuditableDB(db, 1)
	req := createEmployeeTestData()
	employee, err := repo.CreateEmployeeWithAudit(req, auditDB)

	assert.Error(t, err)
	assert.Nil(t, employee)
}

// Test UpdateEmployeeWithAudit - Success
func TestUpdateEmployeeWithAudit_Success(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create initial employee
	createReq := createEmployeeTestData()
	createdEmployee, err := repo.CreateEmployee(createReq)
	require.NoError(t, err)

	// Update employee with audit
	auditDB := middleware.NewAuditableDB(db, 2)
	updateReq := createUpdateEmployeeTestData()
	updatedEmployee, err := repo.UpdateEmployeeWithAudit(strconv.Itoa(int(createdEmployee.ID)), updateReq, auditDB)

	assert.NoError(t, err)
	assert.NotNil(t, updatedEmployee)
	assert.Equal(t, updateReq.Name, updatedEmployee.Name)
	assert.Equal(t, updateReq.Role, updatedEmployee.Role)
	assert.Equal(t, updateReq.Active, updatedEmployee.Active)
	assert.NotEqual(t, updateReq.Password, updatedEmployee.Password) // Password should be hashed

	// Verify audit fields are set
	assert.NotNil(t, updatedEmployee.UpdatedBy)
	assert.Equal(t, uint(2), *updatedEmployee.UpdatedBy)

	// Verify password is hashed correctly
	err = bcrypt.CompareHashAndPassword([]byte(updatedEmployee.Password), []byte(updateReq.Password))
	assert.NoError(t, err)
}

// Test UpdateEmployeeWithAudit - Employee Not Found
func TestUpdateEmployeeWithAudit_EmployeeNotFound(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	auditDB := middleware.NewAuditableDB(db, 1)
	updateReq := createUpdateEmployeeTestData()
	updatedEmployee, err := repo.UpdateEmployeeWithAudit("999", updateReq, auditDB)

	assert.Error(t, err)
	assert.Nil(t, updatedEmployee)
	assert.Contains(t, err.Error(), "record not found")
}

// Test UpdateEmployeeWithAudit - Invalid Employee ID
func TestUpdateEmployeeWithAudit_InvalidEmployeeID(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	auditDB := middleware.NewAuditableDB(db, 1)
	updateReq := createUpdateEmployeeTestData()
	updatedEmployee, err := repo.UpdateEmployeeWithAudit("invalid", updateReq, auditDB)

	assert.Error(t, err)
	assert.Nil(t, updatedEmployee)
}

// Test UpdateEmployeeWithAudit - Password Hashing Error
func TestUpdateEmployeeWithAudit_PasswordHashingError(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create initial employee
	createReq := createEmployeeTestData()
	createdEmployee, err := repo.CreateEmployee(createReq)
	require.NoError(t, err)

	// Update with empty password (should still work)
	auditDB := middleware.NewAuditableDB(db, 1)
	updateReq := request.UpdateEmployeeRequest{
		Name:     "Updated Name",
		Password: "",
		Role:     "admin",
		Active:   false,
	}

	updatedEmployee, err := repo.UpdateEmployeeWithAudit(strconv.Itoa(int(createdEmployee.ID)), updateReq, auditDB)
	assert.NoError(t, err)
	assert.NotNil(t, updatedEmployee)
}

// Test DeleteEmployeeWithAudit - Success
func TestDeleteEmployeeWithAudit_Success(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create employee
	createReq := createEmployeeTestData()
	createdEmployee, err := repo.CreateEmployee(createReq)
	require.NoError(t, err)

	// Delete employee with audit
	auditDB := middleware.NewAuditableDB(db, 3)
	err = repo.DeleteEmployeeWithAudit(strconv.Itoa(int(createdEmployee.ID)), auditDB)

	assert.NoError(t, err)

	// Verify employee is soft deleted
	var deletedEmployee model.Employee
	err = db.Unscoped().First(&deletedEmployee, createdEmployee.ID).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedEmployee.DeletedAt)

	// Verify audit fields are set
	assert.NotNil(t, deletedEmployee.DeletedBy)
	assert.Equal(t, uint(3), *deletedEmployee.DeletedBy)
}

// Test DeleteEmployeeWithAudit - Employee Not Found
func TestDeleteEmployeeWithAudit_EmployeeNotFound(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	auditDB := middleware.NewAuditableDB(db, 1)
	err := repo.DeleteEmployeeWithAudit("999", auditDB)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

// Test DeleteEmployeeWithAudit - Invalid Employee ID
func TestDeleteEmployeeWithAudit_InvalidEmployeeID(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	auditDB := middleware.NewAuditableDB(db, 1)
	err := repo.DeleteEmployeeWithAudit("invalid", auditDB)

	assert.Error(t, err)
}

// Edge Case: Test Multiple Employees with Same Name (if allowed)
func TestCreateEmployee_MultipleEmployeesWithSameName(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	req := createEmployeeTestData()

	// Create first employee
	employee1, err := repo.CreateEmployee(req)
	assert.NoError(t, err)
	assert.NotNil(t, employee1)

	// Create second employee with same name (should succeed if no unique constraint)
	employee2, err := repo.CreateEmployee(req)
	assert.NoError(t, err)
	assert.NotNil(t, employee2)
	assert.NotEqual(t, employee1.ID, employee2.ID)
}

// Edge Case: Test Role Validation
func TestCreateEmployee_DifferentRoles(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	roles := []string{"admin", "employee"}

	for _, role := range roles {
		req := request.CreateEmployeeRequest{
			Name:     "User " + role,
			Password: "password123",
			Role:     role,
			Active:   true,
		}

		employee, err := repo.CreateEmployee(req)
		assert.NoError(t, err)
		assert.NotNil(t, employee)
		assert.Equal(t, role, employee.Role)
	}
}

// Edge Case: Test Active/Inactive States
func TestCreateEmployee_ActiveInactiveStates(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	states := []bool{true, false}

	for i, active := range states {
		req := request.CreateEmployeeRequest{
			Name:     "User " + strconv.Itoa(i),
			Password: "password123",
			Role:     "admin",
			Active:   active,
		}

		employee, err := repo.CreateEmployee(req)
		assert.NoError(t, err)
		assert.NotNil(t, employee)
		assert.Equal(t, active, employee.Active)
	}
}

// Performance Benchmark: CreateEmployee
func BenchmarkCreateEmployee(b *testing.B) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	req := createEmployeeTestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Name = "Benchmark User " + strconv.Itoa(i)
		_, err := repo.CreateEmployee(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Performance Benchmark: GetAllEmployees
func BenchmarkGetAllEmployees(b *testing.B) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Seed data
	seedEmployeesForTest(db, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetAllEmployees()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Performance Benchmark: GetAllActiveEmployees
func BenchmarkGetAllActiveEmployees(b *testing.B) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Seed data
	seedEmployeesForTest(db, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetAllActiveEmployees()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Performance Benchmark: GetEmployeeByID
func BenchmarkGetEmployeeByID(b *testing.B) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Seed data
	employees := seedEmployeesForTest(db, 100)
	targetID := employees[50].ID

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetEmployeeByID(targetID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Performance Benchmark: GetEmployeeByName
func BenchmarkGetEmployeeByName(b *testing.B) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Seed data
	employees := seedEmployeesForTest(db, 100)
	targetName := employees[50].Name

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetEmployeeByName(targetName)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Performance Benchmark: UpdateEmployee
func BenchmarkUpdateEmployee(b *testing.B) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create initial employee
	createReq := createEmployeeTestData()
	employee, err := repo.CreateEmployee(createReq)
	if err != nil {
		b.Fatal(err)
	}

	updateReq := createUpdateEmployeeTestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updateReq.Name = "Updated User " + strconv.Itoa(i)
		_, err := repo.UpdateEmployee(strconv.Itoa(int(employee.ID)), updateReq)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Performance Benchmark: DeleteEmployee
func BenchmarkDeleteEmployee(b *testing.B) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create employee for each iteration
		createReq := createEmployeeTestData()
		createReq.Name = "Delete Test " + strconv.Itoa(i)
		employee, err := repo.CreateEmployee(createReq)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()

		err = repo.DeleteEmployee(strconv.Itoa(int(employee.ID)))
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Concurrent Access Test
func TestEmployeeRepository_ConcurrentAccess(t *testing.T) {
	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	const numGoroutines = 10
	const numOperations = 5

	done := make(chan bool, numGoroutines)

	// Test concurrent creation
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			for j := 0; j < numOperations; j++ {
				req := request.CreateEmployeeRequest{
					Name:     "Concurrent User " + strconv.Itoa(id) + "-" + strconv.Itoa(j),
					Password: "password123",
					Role:     "admin",
					Active:   true,
				}

				_, err := repo.CreateEmployee(req)
				assert.NoError(t, err)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all employees were created
	employees, err := repo.GetAllEmployees()
	assert.NoError(t, err)
	assert.Len(t, employees, numGoroutines*numOperations)
}

// Large Dataset Test
func TestEmployeeRepository_LargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	db := setupEmployeeTestDB()
	repo := NewEmployeeRepository(db)

	// Create 1000 employees
	const numEmployees = 1000

	start := time.Now()
	seedEmployeesForTest(db, numEmployees)
	createDuration := time.Since(start)

	// Test retrieval performance
	start = time.Now()
	employees, err := repo.GetAllEmployees()
	retrievalDuration := time.Since(start)

	assert.NoError(t, err)
	assert.Len(t, employees, numEmployees)

	// Performance assertions (adjust thresholds as needed)
	t.Logf("Created %d employees in %v", numEmployees, createDuration)
	t.Logf("Retrieved %d employees in %v", numEmployees, retrievalDuration)

	assert.Less(t, createDuration, 10*time.Second, "Creation should complete within 10 seconds")
	assert.Less(t, retrievalDuration, 1*time.Second, "Retrieval should complete within 1 second")
}
