CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    created_by INT NOT NULL,
    updated_by INT NOT NULL,
    deleted_by INT NULL DEFAULT NULL
);

CREATE TABLE attendances (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL,
    checkin TIMESTAMP NOT NULL,
    checkout TIMESTAMP NULL,
    hours_worked INT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('present', 'absent', 'leave', 'holiday')),
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    updated_by INT NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by INT NULL DEFAULT NULL,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE overtimes (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL,
    overtime_date DATE NOT NULL,
    start_time TIME NULL,
    end_time TIME NULL,
    hours INT NOT NULL,
    reason VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT NULL,
    updated_by INT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by INT NULL DEFAULT NULL,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE reimbursements (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL,
    reimbursement_date DATE NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    reason VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    category VARCHAR(50) NULL,
    approved_by INT NULL DEFAULT NULL,
    approved_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    updated_by INT NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by INT NULL DEFAULT NULL,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE payslips (
    id SERIAL PRIMARY KEY,
    payPeriodStart DATE NOT NULL,
    payPeriodEnd DATE NOT NULL,
    employee_id INT NOT NULL,
    basic_salary NUMERIC(12, 2) NOT NULL,
    overtime_hours INT NOT NULL,
    overtime_amount NUMERIC(12, 2) NOT NULL,
    reimbursement_amount NUMERIC(12, 2) NOT NULL,
    total_amount NUMERIC(12, 2) NOT NULL,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processed', 'paid')),
    attendance_days INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INT NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by INT NULL DEFAULT NULL,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);