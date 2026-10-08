package handler

import (
	"errors"
	"fmt"
	"net/http"

	"community-waste-collection-api/internal/entity"
	"community-waste-collection-api/internal/service"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service  service.PaymentService
	maxBytes int64
}

func NewPaymentHandler(s service.PaymentService, maxBytes int64) *PaymentHandler {
	if maxBytes <= 0 {
		maxBytes = 5 * 1024 * 1024
	}
	return &PaymentHandler{service: s, maxBytes: maxBytes}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var input entity.CreatePaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, "household_id, waste_id, dan amount (>0) wajib diisi")
		return
	}
	result, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		h.writeError(c, err)
		return
	}
	Created(c, result)
}

func (h *PaymentHandler) List(c *gin.Context) {
	page, limit := service.NormalizePagination(
		atoiDefault(c.Query("page"), 1),
		atoiDefault(c.Query("limit"), 10),
	)
	data, total, err := h.service.List(
		c.Request.Context(),
		c.Query("status"),
		c.Query("household_id"),
		c.Query("start_date"),
		c.Query("end_date"),
		page, limit,
	)
	if err != nil {
		h.writeError(c, err)
		return
	}
	SuccessWithMeta(c, http.StatusOK, data, PagedMeta(total, page, limit))
}

func (h *PaymentHandler) Confirm(c *gin.Context) {
	fh, err := c.FormFile("proof")
	if err != nil {
		Fail(c, http.StatusBadRequest, entity.ErrProofRequired.Error())
		return
	}
	if fh.Size == 0 {
		Fail(c, http.StatusBadRequest, "file bukti kosong")
		return
	}
	if fh.Size > h.maxBytes {
		Fail(c, http.StatusBadRequest, fmt.Sprintf("ukuran file maksimal %dMB", h.maxBytes/(1024*1024)))
		return
	}
	f, err := fh.Open()
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer f.Close()
	result, err := h.service.Confirm(c.Request.Context(), c.Param("id"), fh.Filename, f)
	if err != nil {
		h.writeError(c, err)
		return
	}
	OK(c, result)
}

func (h *PaymentHandler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, entity.ErrInvalidUUID),
		errors.Is(err, entity.ErrValidation),
		errors.Is(err, entity.ErrPaymentAmountMismatch),
		errors.Is(err, entity.ErrProofRequired):
		Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, entity.ErrHouseholdNotFound),
		errors.Is(err, entity.ErrPickupNotFound),
		errors.Is(err, entity.ErrPaymentNotFound):
		Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, entity.ErrPaymentNotPending):
		Fail(c, http.StatusConflict, err.Error())
	default:
		Fail(c, http.StatusInternalServerError, err.Error())
	}
}
