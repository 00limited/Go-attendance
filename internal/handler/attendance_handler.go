package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/helper"
	"github.com/yourname/payslip-system/internal/helper/response"
	"github.com/yourname/payslip-system/internal/repository"

	"gorm.io/gorm"
)

type AttendanceHandler struct {
	Helper   helper.NewHelper
	DB       *gorm.DB
	Response response.Interface

	BaseRepo       repository.BaseRepositoryInterface
	AttendanceRepo repository.AttendanceRepository
}

func (h *AttendanceHandler) CheckinAttendancePeriod(c echo.Context) error {
	req := request.CreateAttendanceRequest{}
	if err := c.Bind(&req); err != nil {
		return h.Response.SendError(c, err.Error(), "Invalid request data")
	}

	// Get auditable database instance
	auditDB := helper.GetAuditableDB(c, h.BaseRepo.GetDB())

	_, err := h.AttendanceRepo.CheckinAttendancePeriodWithAudit(req.EmployeeID, auditDB)
	if err != nil {
		return h.Response.SendError(c, err.Error(), "Failed to create attendance period")
	}
	return h.Response.SendSuccess(c, "Attendance period created successfully", nil)
}

func (h *AttendanceHandler) CheckOutAttendancePeriod(c echo.Context) error {
	req := request.CreateAttendanceRequest{}
	if err := c.Bind(&req); err != nil {
		return h.Response.SendError(c, err.Error(), "Invalid request data")
	}

	// Get auditable database instance
	auditDB := helper.GetAuditableDB(c, h.BaseRepo.GetDB())

	_, err := h.AttendanceRepo.CheckOutAttendancePeriodWithAudit(req.EmployeeID, auditDB)
	if err != nil {
		return h.Response.SendError(c, err.Error(), "Failed to update attendance period")
	}
	return h.Response.SendSuccess(c, "Attendance checkout successful", nil)
}
