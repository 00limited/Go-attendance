package repository

import (
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/middleware"
	"github.com/yourname/payslip-system/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type employee struct {
	db *gorm.DB
}

// NewEmployeeRepository creates a new instance of employee repository.
func NewEmployeeRepository(db *gorm.DB) *employee {
	return &employee{db: db}
}

type EmployeeRepository interface {
	CreateEmployee(req request.CreateEmployeeRequest) (*model.Employee, error)
	GetAllEmployees() ([]model.Employee, error)
	GetAllActiveEmployees() ([]model.Employee, error)
	UpdateEmployee(employeeID string, req request.UpdateEmployeeRequest) (*model.Employee, error)
	DeleteEmployee(employeeID string) error
	GetEmployeeByID(id uint) (*model.Employee, error)
	GetEmployeeByName(name string) (*model.Employee, error)
	CreateEmployeeWithAudit(req request.CreateEmployeeRequest, auditDB *middleware.AuditableDB) (*model.Employee, error)
	UpdateEmployeeWithAudit(employeeID string, req request.UpdateEmployeeRequest, auditDB *middleware.AuditableDB) (*model.Employee, error)
	DeleteEmployeeWithAudit(employeeID string, auditDB *middleware.AuditableDB) error
}

// CreateEmployee creates a new employee record in the database.
func (e *employee) CreateEmployee(req request.CreateEmployeeRequest) (*model.Employee, error) {
	// Hash the password before saving
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	emp := model.Employee{
		Name:     req.Name,
		Password: hashedPassword,
		Role:     req.Role,
		Active:   req.Active,
	}

	err = e.db.Create(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (e *employee) GetAllEmployees() ([]model.Employee, error) {
	var emps []model.Employee
	err := e.db.Debug().Find(&emps).Error
	if err != nil {
		return nil, err
	}
	return emps, nil
}

func (e *employee) GetAllActiveEmployees() ([]model.Employee, error) {
	var emps []model.Employee
	err := e.db.Debug().Where("active = ?", true).Find(&emps).Error
	if err != nil {
		return nil, err
	}
	return emps, nil
}

func (e *employee) UpdateEmployee(employeeID string, req request.UpdateEmployeeRequest) (*model.Employee, error) {
	var emp model.Employee
	err := e.db.Debug().Where("id = ?", employeeID).First(&emp).Error
	if err != nil {
		return nil, err
	}
	// Hash the password before updating
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	emp.Name = req.Name
	emp.Password = hashedPassword
	emp.Role = req.Role
	emp.Active = req.Active
	err = e.db.Save(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (e *employee) DeleteEmployee(employeeID string) error {
	var emp model.Employee
	err := e.db.Debug().Where("id = ?", employeeID).First(&emp).Error
	if err != nil {
		return err
	}
	err = e.db.Delete(&emp).Error
	if err != nil {
		return err
	}
	return nil
}

// GetEmployeeByID retrieves an employee by their ID
func (e *employee) GetEmployeeByID(id uint) (*model.Employee, error) {
	var emp model.Employee
	err := e.db.Debug().Where("id = ?", id).First(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// GetEmployeeByName retrieves an employee by their name (for login)
func (e *employee) GetEmployeeByName(name string) (*model.Employee, error) {
	var emp model.Employee
	err := e.db.Debug().Where("name = ?", name).First(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// CreateEmployeeWithAudit creates a new employee record with audit fields
func (e *employee) CreateEmployeeWithAudit(req request.CreateEmployeeRequest, auditDB *middleware.AuditableDB) (*model.Employee, error) {
	// Hash the password before saving
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	emp := model.Employee{
		Name:     req.Name,
		Password: hashedPassword,
		Role:     req.Role,
		Active:   req.Active,
	}

	err = auditDB.Create(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// UpdateEmployeeWithAudit updates an employee record with audit fields
func (e *employee) UpdateEmployeeWithAudit(employeeID string, req request.UpdateEmployeeRequest, auditDB *middleware.AuditableDB) (*model.Employee, error) {
	var emp model.Employee
	err := e.db.Debug().Where("id = ?", employeeID).First(&emp).Error
	if err != nil {
		return nil, err
	}

	// Hash the password before updating
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	emp.Name = req.Name
	emp.Password = hashedPassword
	emp.Role = req.Role
	emp.Active = req.Active

	err = auditDB.Save(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// DeleteEmployeeWithAudit soft deletes an employee with audit fields
func (e *employee) DeleteEmployeeWithAudit(employeeID string, auditDB *middleware.AuditableDB) error {
	var emp model.Employee
	err := e.db.Debug().Where("id = ?", employeeID).First(&emp).Error
	if err != nil {
		return err
	}

	err = auditDB.Delete(&emp).Error
	if err != nil {
		return err
	}
	return nil
}

// hashPassword hashes a plain password using bcrypt.
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
