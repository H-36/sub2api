package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_CustomCodePassesToService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminSvc := newStubAdminService()
	handler := &RedeemHandler{adminService: adminSvc}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := map[string]any{
		"count": 1,
		"code":  "  CUSTOM-CODE-1  ",
		"type":  "balance",
		"value": 10.0,
	}
	jsonBytes, err := json.Marshal(body)
	require.NoError(t, err)

	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/admin/redeem-codes/generate", bytes.NewReader(jsonBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Generate(c)

	require.Equal(t, http.StatusOK, w.Code)
	require.NotNil(t, adminSvc.lastGenerateRedeemCodesInput)
	assert.Equal(t, "CUSTOM-CODE-1", adminSvc.lastGenerateRedeemCodesInput.Code)
	assert.Equal(t, 1, adminSvc.lastGenerateRedeemCodesInput.Count)
}

func TestGenerate_CustomCodeRequiresSingleCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminSvc := newStubAdminService()
	handler := &RedeemHandler{adminService: adminSvc}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := map[string]any{
		"count": 2,
		"code":  "CUSTOM-CODE-2",
		"type":  "balance",
		"value": 10.0,
	}
	jsonBytes, err := json.Marshal(body)
	require.NoError(t, err)

	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/admin/redeem-codes/generate", bytes.NewReader(jsonBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Generate(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Nil(t, adminSvc.lastGenerateRedeemCodesInput)
}

// newCreateAndRedeemHandler creates a RedeemHandler with a non-nil (but minimal)
// RedeemService so that CreateAndRedeem's nil guard passes and we can test the
// parameter-validation layer that runs before any service call.
func newCreateAndRedeemHandler() *RedeemHandler {
	return &RedeemHandler{
		adminService:  newStubAdminService(),
		redeemService: &service.RedeemService{}, // non-nil to pass nil guard
	}
}

// postCreateAndRedeemValidation calls CreateAndRedeem and returns the response
// status code. For cases that pass validation and proceed into the service layer,
// a panic may occur (because RedeemService internals are nil); this is expected
// and treated as "validation passed" (returns 0 to indicate panic).
func postCreateAndRedeemValidation(t *testing.T, handler *RedeemHandler, body any) (code int) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBytes, err := json.Marshal(body)
	require.NoError(t, err)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/admin/redeem-codes/create-and-redeem", bytes.NewReader(jsonBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			// Panic means we passed validation and entered service layer (expected for minimal stub).
			code = 0
		}
	}()
	handler.CreateAndRedeem(c)
	return w.Code
}

func TestCreateAndRedeem_TypeDefaultsToBalance(t *testing.T) {
	// 不传 type 字段时应默认 balance，不触发 subscription 校验。
	// 验证通过后进入 service 层会 panic（返回 0），说明默认值生效。
	h := newCreateAndRedeemHandler()
	code := postCreateAndRedeemValidation(t, h, map[string]any{
		"code":    "test-balance-default",
		"value":   10.0,
		"user_id": 1,
	})

	assert.NotEqual(t, http.StatusBadRequest, code,
		"omitting type should default to balance and pass validation")
}

func TestCreateAndRedeem_SubscriptionRequiresGroupID(t *testing.T) {
	h := newCreateAndRedeemHandler()
	code := postCreateAndRedeemValidation(t, h, map[string]any{
		"code":          "test-sub-no-group",
		"type":          "subscription",
		"value":         29.9,
		"user_id":       1,
		"validity_days": 30,
		// group_id 缺失
	})

	assert.Equal(t, http.StatusBadRequest, code)
}

func TestCreateAndRedeem_SubscriptionRequiresNonZeroValidityDays(t *testing.T) {
	groupID := int64(5)
	h := newCreateAndRedeemHandler()

	// zero should be rejected
	t.Run("zero", func(t *testing.T) {
		code := postCreateAndRedeemValidation(t, h, map[string]any{
			"code":          "test-sub-bad-days-zero",
			"type":          "subscription",
			"value":         29.9,
			"user_id":       1,
			"group_id":      groupID,
			"validity_days": 0,
		})

		assert.Equal(t, http.StatusBadRequest, code)
	})

	// negative should pass validation (used for refund/reduction)
	t.Run("negative_passes_validation", func(t *testing.T) {
		code := postCreateAndRedeemValidation(t, h, map[string]any{
			"code":          "test-sub-negative-days",
			"type":          "subscription",
			"value":         29.9,
			"user_id":       1,
			"group_id":      groupID,
			"validity_days": -7,
		})

		assert.NotEqual(t, http.StatusBadRequest, code,
			"negative validity_days should pass validation for refund")
	})
}

func TestCreateAndRedeem_SubscriptionValidParamsPassValidation(t *testing.T) {
	groupID := int64(5)
	h := newCreateAndRedeemHandler()
	code := postCreateAndRedeemValidation(t, h, map[string]any{
		"code":          "test-sub-valid",
		"type":          "subscription",
		"value":         29.9,
		"user_id":       1,
		"group_id":      groupID,
		"validity_days": 31,
	})

	assert.NotEqual(t, http.StatusBadRequest, code,
		"valid subscription params should pass validation")
}

func TestCreateAndRedeem_BalanceIgnoresSubscriptionFields(t *testing.T) {
	h := newCreateAndRedeemHandler()
	// balance 类型不传 group_id 和 validity_days，不应报 400
	code := postCreateAndRedeemValidation(t, h, map[string]any{
		"code":    "test-balance-no-extras",
		"type":    "balance",
		"value":   50.0,
		"user_id": 1,
	})

	assert.NotEqual(t, http.StatusBadRequest, code,
		"balance type should not require group_id or validity_days")
}

func TestGetClaims_ReturnsClaimRecords(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminSvc := newStubAdminService()
	now := time.Date(2026, time.April, 17, 9, 30, 0, 0, time.UTC)
	adminSvc.redeemClaims[123] = []service.RedeemCodeClaim{
		{
			ID:           1,
			RedeemCodeID: 123,
			UserID:       42,
			Amount:       10,
			ClaimedAt:    now,
			User: &service.User{
				ID:    42,
				Email: "claimant@example.com",
			},
		},
	}

	router := gin.New()
	router.GET("/api/v1/admin/redeem-codes/:id/claims", (&RedeemHandler{adminService: adminSvc}).GetClaims)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/redeem-codes/123/claims", nil)
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Code int `json:"code"`
		Data []struct {
			ID           int64   `json:"id"`
			RedeemCodeID int64   `json:"redeem_code_id"`
			UserID       int64   `json:"user_id"`
			Amount       float64 `json:"amount"`
			ClaimedAt    string  `json:"claimed_at"`
			User         *struct {
				ID    int64  `json:"id"`
				Email string `json:"email"`
			} `json:"user"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, 0, resp.Code)
	require.Len(t, resp.Data, 1)
	assert.Equal(t, int64(1), resp.Data[0].ID)
	assert.Equal(t, int64(123), resp.Data[0].RedeemCodeID)
	assert.Equal(t, int64(42), resp.Data[0].UserID)
	assert.Equal(t, 10.0, resp.Data[0].Amount)
	assert.Equal(t, "claimant@example.com", resp.Data[0].User.Email)
	assert.Equal(t, now.Format(time.RFC3339), resp.Data[0].ClaimedAt)
}
