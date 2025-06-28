# Employee Attendance & Payroll Management System

A comprehensive Go-based REST API system for managing employee attendance, overtime, reimbursements, and payroll processing with JWT authentication and role-based authorization.

## 🚀 Features

- **Authentication & Authorization**: JWT-based authentication with role-based access control (Admin/Employee)
- **Employee Management**: CRUD operations for employee records (Admin only)
- **Attendance Tracking**: Check-in/check-out functionality with time tracking
- **Overtime Management**: Employee overtime request submission and tracking
- **Reimbursement System**: Employee expense reimbursement requests
- **Payroll Processing**: Automated payroll calculation and payslip generation
- **Audit Logging**: Comprehensive audit trail for all operations
- **Database Migration**: Automated database schema management

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Testing Guide](#testing-guide)
- [Database Schema](#database-schema)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)

## 🏗️ Architecture Overview

### System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client Apps   │    │   Web Browser   │    │   Mobile App    │
│  (Postman/curl) │    │                 │    │                 │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │      REST API Server      │
                    │     (Echo Framework)      │
                    │    Port: 8080 (HTTP)      │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │     Business Logic        │
                    │   (Handlers & UseCases)   │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │    Data Access Layer      │
                    │    (Repository Pattern)   │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │      PostgreSQL DB        │
                    │   (Primary Data Store)    │
                    └───────────────────────────┘
```

### Component Architecture

The application follows Clean Architecture principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                        Presentation Layer                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Routes    │  │  Handlers   │  │ Middleware  │         │
│  │             │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                     Business Logic Layer                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  Use Cases  │  │   Services  │  │  Validators │         │
│  │             │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                    Data Access Layer                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │Repositories │  │   Models    │  │   Database  │         │
│  │             │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└───────────────────────────────────────────────────────────────┘
```

## ✅ Prerequisites

### System Requirements

- **Go**: Version 1.19 or higher
- **PostgreSQL**: Version 12 or higher
- **Git**: For version control
- **Make**: (Optional) For build automation

### Development Tools (Recommended)

- **Visual Studio Code** with Go extension
- **Postman** or **Insomnia** for API testing
- **pgAdmin** or **DBeaver** for database management
- **Docker** (Optional) for containerized development

## 🛠️ Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd Project\ Attendance
```

### 2. Install Dependencies

```bash
go mod download
go mod verify
```

### 3. Setup PostgreSQL Database

#### Option A: Local PostgreSQL Installation

```bash
# Create database
createdb attendance_db

# Create user (optional)
psql -c "CREATE USER attendance_user WITH PASSWORD 'your_password';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE attendance_db TO attendance_user;"
```

#### Option B: Using Docker

```bash
docker run --name postgres-attendance \
  -e POSTGRES_DB=attendance_db \
  -e POSTGRES_USER=attendance_user \
  -e POSTGRES_PASSWORD=your_password \
  -p 5432:5432 \
  -d postgres:13
```

### 4. Environment Configuration

Create a `.env` file in the project root:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5434
DB_USER=attendance_user
DB_PASSWORD=your_password
DB_NAME=attendance_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRY_HOURS=24

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Application Configuration
APP_ENV=development
LOG_LEVEL=debug
```

### 5. Database Migration

```bash
# Run migrations
go run cmd/server/main.go migrate

# Or manually run migrations
migrate -path migrations -database "postgresql://username:password@localhost/attendance_db?sslmode=disable" up
```

### 6. Seed Initial Data

```bash
# Run the application with seed flag
go run cmd/server/main.go --seed
```

This will create:

- Default admin user: `Admin` / `admin123`
- Sample employee users with password: `password123`

## ⚙️ Configuration

### Database Configuration

The application uses GORM for database operations. Configure your database connection in `internal/config/config.go`:

```go
type Config struct {
    Database struct {
        Host     string
        Port     string
        User     string
        Password string
        Name     string
        SSLMode  string
    }
    JWT struct {
        Secret      string
        ExpiryHours int
    }
    Server struct {
        Port string
        Host string
    }
}
```

### JWT Configuration

- **Secret Key**: Use a strong, random secret key for production
- **Expiry**: Default token expiry is 24 hours
- **Refresh**: Tokens can be refreshed using the `/auth/refresh` endpoint

## 🚀 Running the Application

### Development Mode

```bash
# Run with hot reload (if using air)
air

# Or run directly
go run cmd/server/main.go

# Run with specific flags
go run cmd/server/main.go --port=8080 --env=development
```

### Production Mode

```bash
# Build the application
go build -o bin/server cmd/server/main.go

# Run the built binary
./bin/server
```

### Using Docker Compose

```bash
# Start all services (app + database)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

The application will be available at: `http://localhost:8080`

## 📖 API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Available Endpoints

| Method | Endpoint                         | Description              | Access Level   |
| ------ | -------------------------------- | ------------------------ | -------------- |
| GET    | `/health`                        | Health check             | Public         |
| POST   | `/auth/login`                    | User login               | Public         |
| GET    | `/auth/profile`                  | Get user profile         | Authenticated  |
| POST   | `/auth/refresh`                  | Refresh token            | Authenticated  |
| GET    | `/employee/get-all-employee`     | Get all employees        | Admin          |
| POST   | `/employee/create`               | Create employee          | Admin          |
| GET    | `/employee/profile/:id`          | Get employee profile     | Employee/Admin |
| PUT    | `/employee/edit/:id`             | Update employee          | Admin          |
| DELETE | `/employee/delete/:id`           | Delete employee          | Admin          |
| POST   | `/attendance/check-in`           | Check in attendance      | Employee/Admin |
| POST   | `/attendance/check-out`          | Check out attendance     | Employee/Admin |
| POST   | `/overtime/create`               | Create overtime request  | Employee/Admin |
| POST   | `/reimbursement/create`          | Create reimbursement     | Employee/Admin |
| POST   | `/payroll/run`                   | Run payroll for all      | Admin          |
| POST   | `/payroll/run/employee`          | Run payroll for employee | Admin          |
| POST   | `/payroll/summary`               | Get payroll summary      | Admin          |
| GET    | `/payroll/employee/:id/payslips` | Get employee payslips    | Employee/Admin |
| GET    | `/payroll/payslip/:id/details`   | Get payslip details      | Employee/Admin |

For detailed API examples with request/response formats, see [API_TESTING_GUIDE.md](./API_TESTING_GUIDE.md).

## 🧪 Testing Guide

### Automated Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/handler/...

# Run tests with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Manual Testing

#### 1. Setup Test Environment

```bash
# Start the server
go run cmd/server/main.go

# Verify health endpoint
curl http://localhost:8080/health
```

#### 2. Authentication Flow

```bash
# Admin login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"name": "Admin", "password": "admin123"}'

# Save the token from response for subsequent requests
export ADMIN_TOKEN="your-admin-token-here"
```

#### 3. Test Employee Management

```bash
# Get all employees
curl -X GET http://localhost:8080/api/v1/employee/get-all-employee \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Create new employee
curl -X POST http://localhost:8080/api/v1/employee/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name": "Test Employee", "password": "password123", "role": "employee", "active": true}'
```

#### 4. Test Attendance Flow

```bash
# Employee login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Employee", "password": "password123"}'

export EMPLOYEE_TOKEN="your-employee-token-here"

# Check in
curl -X POST http://localhost:8080/api/v1/attendance/check-in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" \
  -d '{"employee_id": 1, "date": "2025-06-28", "check_in_time": "09:00:00"}'
```

### Load Testing

```bash
# Install hey (load testing tool)
go install github.com/rakyll/hey@latest

# Test login endpoint
hey -n 1000 -c 10 -m POST -H "Content-Type: application/json" \
  -d '{"name": "Admin", "password": "admin123"}' \
  http://localhost:8080/api/v1/auth/login
```

## 🗄️ Database Schema

### Entity Relationship Diagram

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    employees    │────▶│   attendances   │     │    overtimes    │
│                 │     │                 │     │                 │
│ - id (PK)       │     │ - id (PK)       │     │ - id (PK)       │
│ - name          │     │ - employee_id   │────▶│ - employee_id   │
│ - password      │     │ - date          │     │ - date          │
│ - role          │     │ - check_in_time │     │ - hours         │
│ - active        │     │ - check_out_time│     │ - description   │
│ - created_at    │     │ - created_at    │     │ - status        │
│ - updated_at    │     │ - updated_at    │     │ - created_at    │
└─────────────────┘     └─────────────────┘     │ - updated_at    │
         │                                      └─────────────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐
│ reimbursements  │     │    payslips     │
│                 │     │                 │
│ - id (PK)       │     │ - id (PK)       │
│ - employee_id   │────▶│ - employee_id   │
│ - amount        │     │ - period_start  │
│ - category      │     │ - period_end    │
│ - description   │     │ - basic_salary  │
│ - receipt_url   │     │ - overtime_pay  │
│ - status        │     │ - reimbursements│
│ - created_at    │     |                 |
│ - updated_at    │     │ - deductions    │
└─────────────────┘     │                 │
                        │ - status        │
                        │ - created_at    │
                        │ - updated_at    │
                        └─────────────────┘
```

### Key Tables

#### employees

- Primary entity for user management
- Stores authentication and role information
- Supports soft deletion

#### attendances

- Daily attendance records
- Tracks check-in/check-out times
- Calculates total working hours

#### overtimes

- Employee overtime requests
- Includes approval workflow
- Links to payroll calculations

#### reimbursements

- Employee expense claims
- Categorized expenses
- Receipt URL storage

#### payslips

- Generated payroll records
- Comprehensive salary breakdown
- Historical payroll data

## 📁 Project Structure

```
Project Attendance/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── database/
│   │   └── database.go             # Database connection setup
│   ├── dto/
│   │   ├── request/                # Request DTOs
│   │   └── response/               # Response DTOs
│   ├── handler/                    # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── employee_handler.go
│   │   ├── attendance_handler.go
│   │   ├── overtime_handler.go
│   │   ├── reimbursement_handler.go
│   │   └── payroll_handler.go
│   ├── helper/                     # Utility functions
│   │   ├── common.go
│   │   ├── helper.go
│   │   └── response/
│   ├── middleware/                 # Custom middleware
│   │   ├── auth_middleware.go
│   │   ├── audit_middleware.go
│   │   └── header_middleware.go
│   ├── model/                      # Data models
│   │   ├── employee_model.go
│   │   ├── attendance_model.go
│   │   ├── overtime_model.go
│   │   ├── reimbursement_model.go
│   │   └── payslip_model.go
│   ├── repository/                 # Data access layer
│   │   ├── base.go
│   │   ├── employee_repo.go
│   │   ├── attendance_repo.go
│   │   ├── overtime_repo.go
│   │   ├── reimbursement_repo.go
│   │   └── payslip_repo.go
│   ├── routes/                     # Route definitions
│   │   ├── route.go
│   │   ├── auth_route.go
│   │   ├── employee_route.go
│   │   ├── attendance_route.go
│   │   ├── overtime_route.go
│   │   ├── reimbursement_route.go
│   │   └── payroll_route.go
│   ├── seed/
│   │   └── seed.go                 # Database seeding
│   └── usecases/                   # Business logic
│       └── payroll_usecase.go
├── migrations/                     # Database migrations
├── bin/                           # Compiled binaries
├── docker-compose.yaml            # Docker configuration
├── go.mod                         # Go module definition
├── go.sum                         # Go module checksums
├── API_TESTING_GUIDE.md          # API testing documentation
├── JWT_AUTHENTICATION.md         # JWT implementation guide
├── MIDDLEWARE_AUTHORIZATION.md    # Authorization guide
├── PAYROLL_API.md                # Payroll API documentation
└── README.md                     # This file
```

### Package Responsibilities

- **cmd/**: Application entry points and CLI commands
- **internal/config/**: Configuration management and environment variables
- **internal/database/**: Database connection and initialization
- **internal/dto/**: Data Transfer Objects for API requests/responses
- **internal/handler/**: HTTP request handlers (controllers)
- **internal/helper/**: Utility functions and common operations
- **internal/middleware/**: Custom middleware for authentication, logging, etc.
- **internal/model/**: Database models and entity definitions
- **internal/repository/**: Data access layer with repository pattern
- **internal/routes/**: HTTP route definitions and grouping
- **internal/seed/**: Database seeding for initial data
- **internal/usecases/**: Business logic and use case implementations

## 🔄 Development Workflow

### Code Style and Standards

#### Go Code Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for code formatting
- Run `golint` for style checking
- Use meaningful variable and function names
- Write comprehensive tests for all handlers

#### Git Workflow

```bash
# Create feature branch
git checkout -b feature/new-feature

# Make changes and commit
git add .
git commit -m "feat: add new attendance validation"

# Push and create pull request
git push origin feature/new-feature
```

#### Commit Message Convention

```
feat: add new feature
fix: bug fix
docs: documentation changes
style: formatting changes
refactor: code refactoring
test: add or update tests
chore: maintenance tasks
```

### Adding New Features

#### 1. Create Model

```go
// internal/model/new_model.go
type NewModel struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### 2. Create Repository

```go
// internal/repository/new_repo.go
type NewRepository interface {
    Create(model *NewModel) error
    GetByID(id uint) (*NewModel, error)
    Update(model *NewModel) error
    Delete(id uint) error
}
```

#### 3. Create Handler

```go
// internal/handler/new_handler.go
type NewHandler struct {
    repo     repository.NewRepository
    response responseHelper.Interface
}

func (h *NewHandler) Create(c echo.Context) error {
    // Implementation
}
```

#### 4. Add Routes

```go
// internal/routes/new_route.go
func (nr *NewRoute) NewRoutes(group *echo.Group) {
    group.POST("/create", handler.Create)
    group.GET("/:id", handler.GetByID)
}
```

#### 5. Add Tests

```go
// internal/handler/new_handler_test.go
func TestNewHandler_Create(t *testing.T) {
    // Test implementation
}
```

## 🚀 Deployment

### Environment Setup

#### Development

```bash
export APP_ENV=development
export LOG_LEVEL=debug
go run cmd/server/main.go
```

#### Staging

```bash
export APP_ENV=staging
export LOG_LEVEL=info
./bin/server
```

#### Production

```bash
export APP_ENV=production
export LOG_LEVEL=warn
./bin/server
```

### Docker Deployment

#### Build Docker Image

```bash
# Build image
docker build -t attendance-system:latest .

# Run container
docker run -d \
  --name attendance-app \
  -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_USER=attendance_user \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=attendance_db \
  -e JWT_SECRET=your-secret-key \
  attendance-system:latest
```

#### Docker Compose Deployment

```bash
# Production deployment
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f app
```

### Common Issues

#### Database Connection Issues

```bash
# Check PostgreSQL status
systemctl status postgresql

# Test connection
psql -h localhost -U attendance_user -d attendance_db

# Check logs
tail -f /var/log/postgresql/postgresql-13-main.log
```

#### JWT Token Issues

```bash
# Verify token (use online JWT decoder)
# Check expiry time
# Verify secret key matches

# Generate new secret
openssl rand -base64 32
```

#### Migration Issues

```bash
# Check migration status
migrate -path migrations -database "postgresql://..." version

# Force migration
migrate -path migrations -database "postgresql://..." force 1

# Reset database
migrate -path migrations -database "postgresql://..." drop
```

#### Performance Issues

```bash
# Check database queries
# Enable query logging in PostgreSQL
# Use EXPLAIN ANALYZE for slow queries

# Monitor goroutines
go tool pprof http://localhost:8080/debug/pprof/goroutine

# Memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap
```

### Error Codes

| Code | Description           | Solution                           |
| ---- | --------------------- | ---------------------------------- |
| 401  | Unauthorized          | Check JWT token validity           |
| 403  | Forbidden             | Verify user role permissions       |
| 404  | Not Found             | Check endpoint URL and resource ID |
| 422  | Validation Error      | Check request payload format       |
| 500  | Internal Server Error | Check server logs for details      |

### Logging

#### Enable Debug Logging

```bash
export LOG_LEVEL=debug
go run cmd/server/main.go
```

#### Log Locations

- Application logs: stdout/stderr
- Database logs: PostgreSQL log directory
- Access logs: Configure in middleware

#### Memory Usage

```bash
curl http://localhost:8080/debug/pprof/heap
```

## 📞 Support

### Documentation

- [API Testing Guide](./API_TESTING_GUIDE.md) - Comprehensive API examples
- [JWT Authentication](./JWT_AUTHENTICATION.md) - Authentication implementation
- [Authorization Middleware](./MIDDLEWARE_AUTHORIZATION.md) - Authorization guide
- [Payroll API](./PAYROLL_API.md) - Payroll system documentation

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

---

**Built with ❤️ using Go, Echo Framework, and PostgreSQL**
