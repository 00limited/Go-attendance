package middleware

import (
	"reflect"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// AuditMiddleware automatically sets created_by, updated_by, and deleted_by fields
// based on the authenticated user from JWT token
func AuditMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the authenticated user ID from context (set by HeaderMiddleware)
			userID, ok := c.Get("user_id").(int)
			if ok && userID > 0 {
				// Store the user ID for use in database operations
				c.Set("audit_user_id", uint(userID))
			}

			return next(c)
		}
	}
}

// SetAuditFields is a helper function to set audit fields on models
func SetAuditFields(model interface{}, operation string, userID uint) {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	switch operation {
	case "create":
		// Set created_by field
		if field := v.FieldByName("CreatedBy"); field.IsValid() && field.CanSet() {
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.ValueOf(&userID))
			} else if field.Kind() == reflect.Uint {
				field.SetUint(uint64(userID))
			}
		}

	case "update":
		// Set updated_by field
		if field := v.FieldByName("UpdatedBy"); field.IsValid() && field.CanSet() {
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.ValueOf(&userID))
			} else if field.Kind() == reflect.Uint {
				field.SetUint(uint64(userID))
			}
		}

	case "delete":
		// Set deleted_by field
		if field := v.FieldByName("DeletedBy"); field.IsValid() && field.CanSet() {
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.ValueOf(&userID))
			} else if field.Kind() == reflect.Uint {
				field.SetUint(uint64(userID))
			}
		}
	}
}

// AuditableDB wraps gorm.DB to automatically set audit fields
type AuditableDB struct {
	*gorm.DB
	UserID uint
}

// NewAuditableDB creates a new auditable database instance
func NewAuditableDB(db *gorm.DB, userID uint) *AuditableDB {
	return &AuditableDB{
		DB:     db,
		UserID: userID,
	}
}

// Create wraps gorm Create with audit fields
func (adb *AuditableDB) Create(value interface{}) *gorm.DB {
	SetAuditFields(value, "create", adb.UserID)
	return adb.DB.Create(value)
}

// Save wraps gorm Save with audit fields (acts as update for existing records)
func (adb *AuditableDB) Save(value interface{}) *gorm.DB {
	SetAuditFields(value, "update", adb.UserID)
	return adb.DB.Save(value)
}

// Updates wraps gorm Updates with updated_by field
func (adb *AuditableDB) Updates(value interface{}) *gorm.DB {
	// For Updates, we need to add updated_by to the update map
	if updateMap, ok := value.(map[string]interface{}); ok {
		updateMap["updated_by"] = adb.UserID
	}
	return adb.DB.Updates(value)
}

// Delete wraps gorm Delete with deleted_by field
func (adb *AuditableDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	// Set deleted_by before soft delete
	if len(conds) > 0 {
		// Update the deleted_by field before soft delete
		adb.DB.Model(value).Where(conds[0], conds[1:]...).Update("deleted_by", adb.UserID)
	} else {
		SetAuditFields(value, "delete", adb.UserID)
	}
	return adb.DB.Delete(value, conds...)
}
