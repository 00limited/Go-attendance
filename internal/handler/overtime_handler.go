package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/helper"
	"github.com/yourname/payslip-system/internal/helper/response"
	"github.com/yourname/payslip-system/internal/repository"

	"gorm.io/gorm"
)

type OvertimeHandler struct {
	Helper   helper.NewHelper
	DB       *gorm.DB
	Response response.Interface

	BaseRepo     repository.BaseRepositoryInterface
	OvertimeRepo repository.OvertimeRepository
}

func (h *OvertimeHandler) CreateOvertime(c echo.Context) error {
	req := request.CreateOvertimeRequest{}
	if err := c.Bind(&req); err != nil {
		return h.Response.SendError(c, err.Error(), "Invalid request data")
	}

	// Get auditable DB instance
	auditDB := helper.GetAuditableDB(c, h.OvertimeRepo.GetDB())

	_, err := h.OvertimeRepo.CreateOvertimePeriodWithAudit(req.EmployeeID, req.Hours, req.Reason, auditDB)
	if err != nil {
		return h.Response.SendError(c, err.Error(), "Failed to create overtime period")
	}

	return h.Response.SendSuccess(c, "Overtime period created successfully", nil)
}
