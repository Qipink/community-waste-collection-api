package handler

import (
	"net/http"

	"community-waste-collection-api/internal/service"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(s service.ReportService) *ReportHandler {
	return &ReportHandler{service: s}
}

func (h *ReportHandler) WasteSummary(c *gin.Context) {
	rows, err := h.service.WasteSummary(c.Request.Context())
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	OK(c, rows)
}

func (h *ReportHandler) PaymentSummary(c *gin.Context) {
	rows, revenue, total, err := h.service.PaymentSummary(c.Request.Context())
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"by_status": rows, "total_revenue": revenue, "total_payments": total,
		},
	})
}
