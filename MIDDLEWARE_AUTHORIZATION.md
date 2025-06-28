# Middleware Authorization Documentation

## Overview

This document outlines the authorization middleware implemented for the payroll/attendance system to ensure proper access control for different user roles.

## Middleware Functions

### 1. AdminOnly

**Purpose**: Restricts access to admin users only.

**Usage**:

```go
adminGroup.Use(middleware.AdminOnly(t.Response))
```

**Authorization Rules**:

- ✅ **Admin**: Full access
- ❌ **Employee**: Access denied
- ❌ **Unauthenticated**: Access denied

**Applied to Routes**:

- `POST /payroll/run` - Run payroll for all employees
- `POST /payroll/run/employee` - Run payroll for specific employee
- `POST /payroll/summary` - Get payroll summary

### 2. EmployeeOrAdmin

**Purpose**: Allows both employees and admins to access routes, but with different permission levels.

**Usage**:

```go
employeeGroup.Use(middleware.EmployeeOrAdmin(t.Response))
```

**Authorization Rules**:

- ✅ **Admin**: Can access any employee's data
- ✅ **Employee**: Can only access their own data
- ❌ **Unauthenticated**: Access denied

**Applied to Routes**:

- `GET /payroll/employee/:id/payslips` - Get employee payslips
- `GET /payroll/payslip/:payslip_id/details` - Get detailed payslip

### 3. ValidateEmployeeAccess (Helper Function)

**Purpose**: Helper function to validate if the current user can access specific employee data.

**Parameters**:

- `c echo.Context`: The Echo context containing user information
- `targetEmployeeID uint`: The employee ID being accessed

**Returns**: `bool` - true if access is allowed, false otherwise

**Logic**:

```go
// Admins can access any employee's data
if role == "admin" {
    return true
}

// Employees can only access their own data
if role == "employee" && userID == targetEmployeeID {
    return true
}

return false
```

## Implementation Details

### Context Variables Set by Middleware

| Variable                | Type     | Description                                |
| ----------------------- | -------- | ------------------------------------------ |
| `user_id`               | `int`    | User ID from JWT token                     |
| `role`                  | `string` | User role (`admin` or `employee`)          |
| `authenticated_user_id` | `uint`   | Converted user ID for authorization checks |
| `authenticated_role`    | `string` | User role for authorization checks         |

### Response Codes

| Scenario                        | HTTP Code | Message                                                 |
| ------------------------------- | --------- | ------------------------------------------------------- |
| Admin access required           | 403       | "Access denied. Admin privileges required."             |
| Employee/Admin access required  | 403       | "Access denied. Employee or admin privileges required." |
| Authentication required         | 401       | "Authentication required"                               |
| Employee accessing other's data | 403       | "Access denied. You can only access your own payslips." |

## Route Protection Summary

### Admin-Only Routes

```
POST /payroll/run                    # Process payroll for all employees
POST /payroll/run/employee          # Process payroll for specific employee
POST /payroll/summary               # Get payroll summary for management
```

### Employee + Admin Routes

```
GET /payroll/employee/:id/payslips           # Get employee payslips (with data ownership check)
GET /payroll/payslip/:payslip_id/details    # Get detailed payslip (with data ownership check)
```

## Security Features

### 1. **Role-Based Access Control (RBAC)**

- Different middleware for different permission levels
- Clear separation between admin and employee capabilities

### 2. **Data Ownership Validation**

- Employees can only access their own data
- Admins can access any employee's data
- Validation happens at both middleware and handler levels

### 3. **JWT Token Integration**

- Uses existing JWT authentication system
- Extracts user ID and role from token claims
- Maintains session context throughout request lifecycle

### 4. **Consistent Error Responses**

- Standardized error messages and HTTP status codes
- Uses application's response helper for consistent formatting

## Usage Examples

### Basic Admin Route Protection

```go
adminGroup := c.Group("")
adminGroup.Use(middleware.AdminOnly(t.Response))
adminGroup.POST("/sensitive-admin-action", handler.AdminAction)
```

### Employee Data Access with Ownership Check

```go
employeeGroup := c.Group("")
employeeGroup.Use(middleware.EmployeeOrAdmin(t.Response))
employeeGroup.GET("/employee/:id/data", handler.GetEmployeeData)

// In handler:
func (h *Handler) GetEmployeeData(c echo.Context) error {
    empID := extractEmployeeIDFromParams(c)

    if !middleware.ValidateEmployeeAccess(c, empID) {
        return h.response.SendCustomResponse(c, 403, "Access denied", nil)
    }

    // Proceed with data retrieval
}
```

## Testing Authorization

### Test Cases to Verify

1. **Admin accessing admin-only routes** → ✅ Success
2. **Employee accessing admin-only routes** → ❌ 403 Forbidden
3. **Admin accessing employee data** → ✅ Success (any employee)
4. **Employee accessing own data** → ✅ Success
5. **Employee accessing other employee's data** → ❌ 403 Forbidden
6. **Unauthenticated user accessing protected routes** → ❌ 401 Unauthorized

This middleware system ensures robust authorization while maintaining clean separation of concerns and consistent error handling throughout the application.
