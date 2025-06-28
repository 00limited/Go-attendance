# API Testing Examples for JWT Authentication

This document provides example API calls to test the JWT authentication system.

## Prerequisites

1. Start the server: `go run cmd/server/main.go`
2. The server will run on `http://localhost:8080`
3. Default admin credentials: `Admin` / `admin123`
4. Default employee credentials: any seeded employee name / `password123`

## 0. Health Check (No Authentication Required)

```bash
curl -X GET http://localhost:8080/health
```

Expected Response:

````json
{
  "message": "Server is up and running"
}
```## 20. Get Employee Payslips

```bash
curl -## 21. Get Detailed Payslip

```bash
curl -X GET http://localhost:8080/api/v1/payroll/payslip/PAYSLIP_ID/details \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
````

Expected Response:

```json
{
  "status": "success",
  "message": "Detailed payslip retrieved successfully",
  "data": {
    "id": 15,
    "employee_id": 1,
    "employee_name": "John Doe",
    "period_start": "2025-06-01",
    "period_end": "2025-06-30",
    "basic_salary": 5000000,
    "overtime_pay": 500000,
    "reimbursements": 150000,
    "gross_salary": 5650000,
    "tax_deduction": 1130000,
    "insurance_deduction": 270000,
    "total_deductions": 1400000,
    "net_salary": 4250000,
    "attendance_summary": {
      "total_days": 22,
      "present_days": 20,
      "overtime_hours": 10
    },
    "payment_status": "paid",
    "payment_date": "2025-06-30T10:00:00Z",
    "created_at": "2025-06-28T16:30:00Z"
  }
}
```

_Note: Employees can only access their own payslip details, while admins can access any payslip details._

## 1. Notes for Testing://localhost:8080/api/v1/payroll/employee/EMPLOYEE_ID/payslips \

-H "Authorization: Bearer YOUR_TOKEN_HERE"

````

Expected Response:

```json
{
  "status": "success",
  "message": "Payslips retrieved successfully",
  "data": [
    {
      "id": 15,
      "employee_id": 1,
      "period_start": "2025-06-01",
      "period_end": "2025-06-30",
      "gross_salary": 5000000,
      "net_salary": 4250000,
      "status": "paid",
      "created_at": "2025-06-28T16:30:00Z"
    },
    {
      "id": 10,
      "employee_id": 1,
      "period_start": "2025-05-01",
      "period_end": "2025-05-31",
      "gross_salary": 5000000,
      "net_salary": 4250000,
      "status": "paid",
      "created_at": "2025-05-31T16:30:00Z"
    }
  ]
}
````

_Note: Employees can only access their own payslips, while admins can access any employee's payslips._

## 2. Get Detailed PayslipAdmin Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin",
    "password": "admin123"
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_at": "2025-06-29T10:30:00Z",
    "user": {
      "id": 101,
      "name": "Admin",
      "role": "admin",
      "active": true
    }
  }
}
```

## 3. Get Admin Profile

```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE"
```

Expected Response:

```json
{
  "status": "success",
  "message": "Profile retrieved successfully",
  "data": {
    "id": 101,
    "name": "Admin",
    "role": "admin",
    "active": true
  }
}
```

## 4. Admin-Only: Get All Employees

```bash
curl -X GET http://localhost:8080/api/v1/employee/get-all-employee \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE"
```

Expected Response:

```json
{
  "status": "success",
  "message": "Employees retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "role": "employee",
      "active": true,
      "created_at": "2025-06-01T00:00:00Z",
      "updated_at": "2025-06-01T00:00:00Z"
    },
    {
      "id": 2,
      "name": "Jane Smith",
      "role": "employee",
      "active": true,
      "created_at": "2025-06-01T00:00:00Z",
      "updated_at": "2025-06-01T00:00:00Z"
    }
  ]
}
```

## 5. Admin-Only: Create New Employee

```bash
curl -X POST http://localhost:8080/api/v1/employee/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE" \
  -d '{
    "name": "New Employee",
    "password": "newpassword123",
    "role": "employee",
    "active": true
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Employee created successfully",
  "data": {
    "id": 3,
    "name": "New Employee",
    "role": "employee",
    "active": true,
    "created_at": "2025-06-28T10:30:00Z",
    "updated_at": "2025-06-28T10:30:00Z"
  }
}
```

## 6. Employee Login (Use Actual Employee Name from Database)

First, get an employee name from the database, then:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ACTUAL_EMPLOYEE_NAME",
    "password": "password123"
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_at": "2025-06-29T10:30:00Z",
    "user": {
      "id": 1,
      "name": "John Doe",
      "role": "employee",
      "active": true
    }
  }
}
```

## 7. Employee: Check Own Profile

```bash
curl -X GET http://localhost:8080/api/v1/employee/profile/EMPLOYEE_ID \
  -H "Authorization: Bearer YOUR_EMPLOYEE_TOKEN_HERE"
```

Expected Response:

```json
{
  "status": "success",
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "role": "employee",
    "active": true,
    "created_at": "2025-06-01T00:00:00Z",
    "updated_at": "2025-06-01T00:00:00Z"
  }
}
```

## 8. Employee: Check-in Attendance

```bash
curl -X POST http://localhost:8080/api/v1/attendance/check-in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_EMPLOYEE_TOKEN_HERE" \
  -d '{
    "employee_id": EMPLOYEE_ID,
    "date": "2025-06-28",
    "check_in_time": "09:00:00"
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Check-in recorded successfully",
  "data": {
    "id": 1,
    "employee_id": 1,
    "date": "2025-06-28",
    "check_in_time": "09:00:00",
    "check_out_time": null,
    "created_at": "2025-06-28T09:00:00Z",
    "updated_at": "2025-06-28T09:00:00Z"
  }
}
```

## 9. Employee: Create Overtime Request

```bash
curl -X POST http://localhost:8080/api/v1/overtime/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_EMPLOYEE_TOKEN_HERE" \
  -d '{
    "employee_id": EMPLOYEE_ID,
    "date": "2025-06-28",
    "hours": 2,
    "description": "Project deadline work"
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Overtime request created successfully",
  "data": {
    "id": 1,
    "employee_id": 1,
    "date": "2025-06-28",
    "hours": 2,
    "description": "Project deadline work",
    "status": "pending",
    "created_at": "2025-06-28T17:30:00Z",
    "updated_at": "2025-06-28T17:30:00Z"
  }
}
```

## 10. Employee: Create Reimbursement Request

```bash
curl -X POST http://localhost:8080/api/v1/reimbursement/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_EMPLOYEE_TOKEN_HERE" \
  -d '{
    "employee_id": EMPLOYEE_ID,
    "amount": 50000,
    "category": "transport",
    "description": "Business trip transportation",
    "receipt_url": "https://example.com/receipt.jpg"
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Reimbursement request created successfully",
  "data": {
    "id": 1,
    "employee_id": 1,
    "amount": 50000,
    "category": "transport",
    "description": "Business trip transportation",
    "receipt_url": "https://example.com/receipt.jpg",
    "status": "pending",
    "created_at": "2025-06-28T14:30:00Z",
    "updated_at": "2025-06-28T14:30:00Z"
  }
}
```

## 11. Admin-Only: Run Payroll

```bash
curl -X POST http://localhost:8080/api/v1/payroll/run \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE" \
  -d '{
    "period_start": "2025-06-01",
    "period_end": "2025-06-30"
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Payroll processed successfully for all employees",
  "data": {
    "processed_employees": 15,
    "total_payroll_amount": 75000000,
    "period_start": "2025-06-01",
    "period_end": "2025-06-30",
    "processing_date": "2025-06-28T15:30:00Z"
  }
}
```

## 12. Token Refresh

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

Expected Response:

```json
{
  "status": "success",
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_at": "2025-06-29T15:30:00Z",
    "user": {
      "id": 1,
      "name": "John Doe",
      "role": "employee",
      "active": true
    }
  }
}
```

## 13. Testing Access Denied (Employee Trying Admin Endpoint)

```bash
curl -X GET http://localhost:8080/api/v1/employee/get-all-employee \
  -H "Authorization: Bearer YOUR_EMPLOYEE_TOKEN_HERE"
```

Expected Response:

```json
{
  "status": "error",
  "message": "Access denied. Admin privileges required."
}
```

## 14. Testing Unauthorized Access (No Token)

```bash
curl -X GET http://localhost:8080/api/v1/employee/get-all-employee
```

Expected Response:

```json
{
  "status": "error",
  "message": "missing or malformed jwt"
}
```

## Security Testing Scenarios

### Invalid Token

```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer invalid.token.here"
```

### Expired Token

Use a token that has been expired (after 24 hours)

### Wrong Credentials

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin",
    "password": "wrongpassword"
  }'
```

### Employee Accessing Another Employee's Data

```bash
curl -X GET http://localhost:8080/api/v1/employee/profile/DIFFERENT_EMPLOYEE_ID \
  -H "Authorization: Bearer YOUR_EMPLOYEE_TOKEN_HERE"
```

Expected Response:

```json
{
  "status": "error",
  "message": "Access denied. You can only view your own profile."
}
```

## 15. Admin-Only: Edit Employee

```bash
curl -X PUT http://localhost:8080/api/v1/employee/edit/EMPLOYEE_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE" \
  -d '{
    "name": "Updated Employee Name",
    "password": "newpassword123",
    "role": "employee",
    "active": true
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Employee updated successfully",
  "data": {
    "id": 1,
    "name": "Updated Employee Name",
    "role": "employee",
    "active": true,
    "created_at": "2025-06-01T00:00:00Z",
    "updated_at": "2025-06-28T16:00:00Z"
  }
}
```

## 16. Admin-Only: Delete Employee

```bash
curl -X DELETE http://localhost:8080/api/v1/employee/delete/EMPLOYEE_ID \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE"
```

Expected Response:

```json
{
  "status": "success",
  "message": "Employee deleted successfully",
  "data": {
    "deleted_employee_id": 1,
    "deleted_at": "2025-06-28T16:30:00Z"
  }
}
```

## 17. Employee: Check-out Attendance

```bash
curl -X POST http://localhost:8080/api/v1/attendance/check-out \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_EMPLOYEE_TOKEN_HERE" \
  -d '{
    "employee_id": EMPLOYEE_ID,
    "date": "2025-06-28",
    "check_out_time": "17:00:00"
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Check-out recorded successfully",
  "data": {
    "id": 1,
    "employee_id": 1,
    "date": "2025-06-28",
    "check_in_time": "09:00:00",
    "check_out_time": "17:00:00",
    "total_hours": 8,
    "created_at": "2025-06-28T09:00:00Z",
    "updated_at": "2025-06-28T17:00:00Z"
  }
}
```

## 18. Admin-Only: Run Payroll for Specific Employee

```bash
curl -X POST http://localhost:8080/api/v1/payroll/run/employee \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE" \
  -d '{
    "employee_id": EMPLOYEE_ID,
    "period_start": "2025-06-01",
    "period_end": "2025-06-30"
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Payroll processed successfully for employee",
  "data": {
    "employee_id": 1,
    "employee_name": "John Doe",
    "payslip_id": 15,
    "gross_salary": 5000000,
    "net_salary": 4250000,
    "period_start": "2025-06-01",
    "period_end": "2025-06-30",
    "processing_date": "2025-06-28T16:30:00Z"
  }
}
```

## 19. Admin-Only: Get Payroll Summary

```bash
curl -X POST http://localhost:8080/api/v1/payroll/summary \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE" \
  -d '{
    "period_start": "2025-06-01",
    "period_end": "2025-06-30"
  }'
```

Expected Response:

```json
{
  "status": "success",
  "message": "Payroll summary retrieved successfully",
  "data": {
    "period_start": "2025-06-01",
    "period_end": "2025-06-30",
    "total_employees": 15,
    "total_gross_salary": 75000000,
    "total_net_salary": 63750000,
    "total_deductions": 11250000,
    "total_overtime_pay": 5000000,
    "total_reimbursements": 2500000,
    "summary_generated_at": "2025-06-28T17:00:00Z"
  }
}
```

## 20. Get Employee Payslips

```bash
curl -X GET http://localhost:8080/api/v1/payroll/employee/EMPLOYEE_ID/payslips \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

_Note: Employees can only access their own payslips, while admins can access any employee's payslips._

## 21. Get Detailed Payslip

```bash
curl -X GET http://localhost:8080/api/v1/payroll/payslip/PAYSLIP_ID/details \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

_Note: Employees can only access their own payslip details, while admins can access any payslip details._

## Notes for Testing

1. Replace `YOUR_ADMIN_TOKEN_HERE` and `YOUR_EMPLOYEE_TOKEN_HERE` with actual tokens from login responses
2. Replace `EMPLOYEE_ID` with actual employee IDs from the database
3. Replace `ACTUAL_EMPLOYEE_NAME` with names from seeded employees
4. All timestamps should be in appropriate format for your timezone
5. The server must be running for these tests to work
6. Use tools like Postman for easier testing with GUI
