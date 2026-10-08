package handler

import (
	"errors"
	"net/http"

	"community-waste-collection-api/internal/entity"
	"community-waste-collection-api/internal/service"

	"github.com/gin-gonic/gin"
)

type PickupHandler struct {
	service service.PickupService
}

func NewPickupHandler(s service.PickupService) *PickupHandler {
	return &PickupHandler{service: s}
}

func (h *PickupHandler) Create(c *gin.Context) {
	var input entity.CreatePickupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, "household_id dan type wajib diisi (organic/plastic/paper/electronic)")
		return
	}
	result, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	Created(c, result)
}

func (h *PickupHandler) List(c *gin.Context) {
	page, limit := service.NormalizePagination(
		atoiDefault(c.Query("page"), 1),
		atoiDefault(c.Query("limit"), 10),
	)
	data, total, err := h.service.List(c.Request.Context(), c.Query("status"), c.Query("household_id"), page, limit)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	SuccessWithMeta(c, http.StatusOK, data, PagedMeta(total, page, limit))
}

func (h *PickupHandler) Schedule(c *gin.Context) {
	var input entity.SchedulePickupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, "Tanggal pickup wajib diisi, contoh: 2026-10-12")
		return
	}
	result, err := h.service.Schedule(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	OK(c, result)
}

func (h *PickupHandler) Complete(c *gin.Context) {
	pickup, payment, err := h.service.Complete(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"pickup": pickup, "payment": payment}})
}

func (h *PickupHandler) Cancel(c *gin.Context) {
	result, err := h.service.Cancel(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	OK(c, result)
}

func (h *PickupHandler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, entity.ErrInvalidUUID),
		errors.Is(err, entity.ErrValidation):
		Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, entity.ErrHouseholdNotFound),
		errors.Is(err, entity.ErrPickupNotFound):
		Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, entity.ErrHasPendingPayment),
		errors.Is(err, entity.ErrPickupNotPending),
		errors.Is(err, entity.ErrPickupNotScheduled),
		errors.Is(err, entity.ErrPickupAlreadyFinal):
		Fail(c, http.StatusConflict, err.Error())
	case errors.Is(err, entity.ErrSafetyCheckRequired):
		Fail(c, http.StatusUnprocessableEntity, err.Error())
	default:
		Fail(c, http.StatusInternalServerError, err.Error())
	}
}
