package handler

import (
	"errors"
	"net/http"

	"community-waste-collection-api/internal/entity"
	"community-waste-collection-api/internal/service"

	"github.com/gin-gonic/gin"
)

type HouseholdHandler struct {
	service service.HouseholdService
}

func NewHouseholdHandler(s service.HouseholdService) *HouseholdHandler {
	return &HouseholdHandler{service: s}
}

func (h *HouseholdHandler) Create(c *gin.Context) {
	var input entity.CreateHouseholdInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, "owner_name dan address wajib diisi")
		return
	}
	result, err := h.service.CreateHousehold(c.Request.Context(), input)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	Created(c, result)
}

func (h *HouseholdHandler) GetAll(c *gin.Context) {
	page, limit := service.NormalizePagination(
		atoiDefault(c.Query("page"), 1),
		atoiDefault(c.Query("limit"), 10),
	)
	data, total, err := h.service.GetAllHouseholds(c.Request.Context(), page, limit)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessWithMeta(c, http.StatusOK, data, PagedMeta(total, page, limit))
}

func (h *HouseholdHandler) GetByID(c *gin.Context) {
	data, err := h.service.GetHouseholdByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, entity.ErrInvalidUUID) {
			Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, entity.ErrHouseholdNotFound) {
			Fail(c, http.StatusNotFound, err.Error())
			return
		}
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	OK(c, data)
}

func (h *HouseholdHandler) Delete(c *gin.Context) {
	if err := h.service.DeleteHousehold(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, entity.ErrInvalidUUID) {
			Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, entity.ErrHouseholdNotFound) {
			Fail(c, http.StatusNotFound, err.Error())
			return
		}
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "household berhasil dihapus"})
}
