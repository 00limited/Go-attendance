package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/helper/response"
)

var JWT_SECRET = []byte("super-secret-key")

// HeaderMiddleware is a middleware that adds a custom header to the response.
func HeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("role", claims["role"].(string))

		return next(c)
	}
}

// AdminOnly ensures only admin users can access the route
func AdminOnly(response response.Interface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok || role != "admin" {
				return response.SendCustomResponse(c, 403, "Access denied. Admin privileges required.", nil)
			}
			return next(c)
		}
	}
}

// IsAdmin is a legacy function, kept for backward compatibility
func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != "admin" {
			return echo.NewHTTPError(403, "Forbidden: You do not have permission to access this resource")
		}
		return next(c)
	}
}

// EmployeeOrAdmin allows both employees and admins to access the route
// Employees can only access their own data, admins can access any data
func EmployeeOrAdmin(response response.Interface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok {
				return response.SendUnauthorized(c, "Authentication required", nil)
			}

			if role != "employee" && role != "admin" {
				return response.SendCustomResponse(c, 403, "Access denied. Employee or admin privileges required.", nil)
			}

			// Set additional context for authorization checks in handlers
			userID, ok := c.Get("user_id").(int)
			if ok {
				c.Set("authenticated_user_id", uint(userID))
			}
			c.Set("authenticated_role", role)

			return next(c)
		}
	}
}
