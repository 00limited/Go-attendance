package model

import "time"

// Employee represents an employee in the system.
type Employee struct {
	DefaultAttribute
	Name     string `json:"name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Password string `json:"-" gorm:"not null;size:255"` // Excluded from JSON for security
	Role     string `json:"role" gorm:"not null;size:50;check:role IN ('admin','employee')" validate:"required,oneof=admin employee"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Relationships
	Attendances    []Attendance    `json:"attendances,omitempty" gorm:"foreignKey:EmployeeID"`
	Overtimes      []Overtime      `json:"overtimes,omitempty" gorm:"foreignKey:EmployeeID"`
	Reimbursements []Reimbursement `json:"reimbursements,omitempty" gorm:"foreignKey:EmployeeID"`
	Payslips       []Payslip       `json:"payslips,omitempty" gorm:"foreignKey:EmployeeID"`
}

// TableName returns the table name for the Employee model.
func (Employee) TableName() string {
	return "employees"
}

// SafeEmployee returns employee data without sensitive information
type SafeEmployee struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToSafe converts Employee to SafeEmployee
func (e *Employee) ToSafe() SafeEmployee {
	return SafeEmployee{
		ID:        e.ID,
		Name:      e.Name,
		Role:      e.Role,
		Active:    e.Active,
		CreatedAt: *e.CreatedAt,
		UpdatedAt: *e.UpdatedAt,
	}
}
