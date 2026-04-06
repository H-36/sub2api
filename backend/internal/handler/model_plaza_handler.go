package handler

import (
	"context"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type modelPlazaGetter interface {
	Get(ctx context.Context) (*service.ModelPlazaResponse, error)
}

// ModelPlazaHandler handles the user-facing model plaza overview endpoint.
type ModelPlazaHandler struct {
	modelPlazaService modelPlazaGetter
}

func NewModelPlazaHandler(modelPlazaService *service.ModelPlazaService) *ModelPlazaHandler {
	return &ModelPlazaHandler{
		modelPlazaService: modelPlazaService,
	}
}

// Get returns the public model plaza overview for authenticated users.
// GET /api/v1/model-plaza
func (h *ModelPlazaHandler) Get(c *gin.Context) {
	data, err := h.modelPlazaService.Get(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, data)
}
