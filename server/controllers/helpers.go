package controllers

import (
	"errors"
	"net/http"
	"x-tract/data"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func parseUUIDParam(c *gin.Context, param string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(param))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return uuid.Nil, false
	}
	return id, true
}

func bindJSON[T any](c *gin.Context) (*T, bool) {
	var input T
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, false
	}
	return &input, true
}

func respondWithData(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

func respondNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func writeError(c *gin.Context, err error) {
	switch {
	case err == nil:
		return
	case errors.Is(err, data.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, data.ErrConflict), errors.Is(err, data.ErrInvalidRole), errors.Is(err, data.ErrForeignKeyViolation):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
