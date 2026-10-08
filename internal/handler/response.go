package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{"status": "success", "data": data})
}

func SuccessWithMeta(c *gin.Context, code int, data interface{}, meta interface{}) {
	c.JSON(code, gin.H{"status": "success", "data": data, "meta": meta})
}

func Created(c *gin.Context, data interface{}) { Success(c, 201, data) }

func OK(c *gin.Context, data interface{}) { Success(c, 200, data) }

func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"status": "error", "message": message})
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

// PagedMeta menyusun meta pagination yang seragam untuk semua list endpoint.
func PagedMeta(total int64, page, limit int) gin.H {
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	return gin.H{"total": total, "page": page, "limit": limit, "total_pages": totalPages}
}
