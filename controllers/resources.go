package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CRUDStore describes the persistence behavior required by a REST resource.
type CRUDStore[Resource any, Create any, Update any] interface {
	List(ctx context.Context) ([]Resource, error)
	Create(ctx context.Context, input *Create) (*Resource, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Resource, error)
	Update(ctx context.Context, id uuid.UUID, input *Update) (*Resource, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// CRUDController is a reusable Gin controller for standard CRUD resources.
type CRUDController[Resource any, Create any, Update any] struct {
	store CRUDStore[Resource, Create, Update]
}

// NewCRUDController builds a CRUDController from a resource store.
func NewCRUDController[Resource any, Create any, Update any](store CRUDStore[Resource, Create, Update]) *CRUDController[Resource, Create, Update] {
	return &CRUDController[Resource, Create, Update]{store: store}
}

// List handles GET collection requests.
func (ctl *CRUDController[Resource, Create, Update]) List(c *gin.Context) {
	items, err := ctl.store.List(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}

	respondWithData(c, http.StatusOK, items)
}

// Create handles POST collection requests.
func (ctl *CRUDController[Resource, Create, Update]) Create(c *gin.Context) {
	input, ok := bindJSON[Create](c)
	if !ok {
		return
	}

	item, err := ctl.store.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err)
		return
	}

	respondWithData(c, http.StatusCreated, item)
}

// GetByID handles GET item requests.
func (ctl *CRUDController[Resource, Create, Update]) GetByID(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	item, err := ctl.store.GetByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}

	respondWithData(c, http.StatusOK, item)
}

// Update handles PATCH item requests.
func (ctl *CRUDController[Resource, Create, Update]) Update(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	input, ok := bindJSON[Update](c)
	if !ok {
		return
	}

	item, err := ctl.store.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err)
		return
	}

	respondWithData(c, http.StatusOK, item)
}

// Delete handles DELETE item requests.
func (ctl *CRUDController[Resource, Create, Update]) Delete(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := ctl.store.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}

	respondNoContent(c)
}
