package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourname/payslip-system/internal/dto/request"
	"github.com/yourname/payslip-system/internal/helper"
	"github.com/yourname/payslip-system/internal/helper/response"
	"github.com/yourname/payslip-system/internal/repository"
	"gorm.io/gorm"
)

type ReimbusementHandler struct {
	Helper   helper.NewHelper
	DB       *gorm.DB
	Response response.Interface

	BaseRepo         repository.BaseRepositoryInterface
	ReimbusementRepo repository.ReimbusementRepository
}

func (h *ReimbusementHandler) CreateReimbusement(c echo.Context) error {
	req := request.CreateReimbusementRequest{}
	if err := c.Bind(&req); err != nil {
		return h.Response.SendError(c, err.Error(), "Invalid request data")
	}

	// Get auditable DB instance
	auditDB := helper.GetAuditableDB(c, h.ReimbusementRepo.GetDB())

	_, err := h.ReimbusementRepo.CreateReimbusementWithAudit(req.EmployeeID, req.Amount, req.Description, auditDB)
	if err != nil {
		return h.Response.SendError(c, err.Error(), "Failed to create reimbusement")
	}
	return h.Response.SendSuccess(c, "Reimbusement created successfully", nil)
}
