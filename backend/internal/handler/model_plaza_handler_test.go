package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type modelPlazaGetterStub struct {
	response *service.ModelPlazaResponse
	err      error
}

func (s modelPlazaGetterStub) Get(_ context.Context) (*service.ModelPlazaResponse, error) {
	return s.response, s.err
}

func TestModelPlazaHandler_Get_Success(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/model-plaza", nil)

	handler := &ModelPlazaHandler{
		modelPlazaService: modelPlazaGetterStub{
			response: &service.ModelPlazaResponse{
				Summary: service.ModelPlazaSummary{
					PlatformCount: 1,
					GroupCount:    1,
					ModelCount:    2,
				},
				Platforms: []service.ModelPlazaPlatform{
					{
						Platform:   service.PlatformOpenAI,
						Label:      "OpenAI",
						GroupCount: 1,
						Groups: []service.ModelPlazaGroup{
							{
								ID:             1,
								Name:           "Pro",
								Platform:       service.PlatformOpenAI,
								RateMultiplier: 1.5,
								ModelCount:     2,
							},
						},
					},
				},
			},
		},
	}

	handler.Get(c)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var payload response.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if payload.Code != 0 {
		t.Fatalf("expected response code 0, got %d", payload.Code)
	}

	data, ok := payload.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected response data object, got %T", payload.Data)
	}

	summary, ok := data["summary"].(map[string]any)
	if !ok {
		t.Fatalf("expected summary object, got %T", data["summary"])
	}

	if summary["platform_count"] != float64(1) {
		t.Fatalf("expected platform_count 1, got %#v", summary["platform_count"])
	}
}

func TestModelPlazaHandler_Get_Error(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/model-plaza", nil)

	handler := &ModelPlazaHandler{
		modelPlazaService: modelPlazaGetterStub{
			err: errors.New("boom"),
		},
	}

	handler.Get(c)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}
