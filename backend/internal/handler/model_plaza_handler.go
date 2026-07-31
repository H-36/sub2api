package handler

import (
	"log/slog"
	"sort"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

const modelPlazaLegacyPriceScale = 1_000_000

// ModelPlazaHandler 处理「模型广场」查询。
//
// 广场路由挂 OptionalJWT 中间件：匿名可访问（除非 require_auth 开启），带 token 则
// 识别用户。可见性规则（橱窗语义，与「可用渠道」的可绑定语义不同）：
//   - 匿名：仅非专属分组（订阅型照常展示）；
//   - 登录：非专属分组 + user_allowed_groups 授权的专属分组（不检查订阅有效性）。
type ModelPlazaHandler struct {
	channelService *service.ChannelService
	apiKeyService  *service.APIKeyService
	settingService *service.SettingService
}

// NewModelPlazaHandler 创建模型广场 handler。
func NewModelPlazaHandler(
	channelService *service.ChannelService,
	apiKeyService *service.APIKeyService,
	settingService *service.SettingService,
) *ModelPlazaHandler {
	return &ModelPlazaHandler{
		channelService: channelService,
		apiKeyService:  apiKeyService,
		settingService: settingService,
	}
}

// modelPlazaOfficialPricing LiteLLM 官方参考价（USD per token）。
type modelPlazaOfficialPricing struct {
	InputPrice        *float64 `json:"input_price"`
	OutputPrice       *float64 `json:"output_price"`
	CacheWritePrice   *float64 `json:"cache_write_price"`
	CacheWrite1hPrice *float64 `json:"cache_write_1h_price,omitempty"`
	CacheReadPrice    *float64 `json:"cache_read_price"`
}

// modelPlazaModel 广场模型条目：渠道定价（白名单形态）+ 官方参考价。
type modelPlazaModel struct {
	Name            string                     `json:"name"`
	Platform        string                     `json:"platform"`
	Pricing         *userSupportedModelPricing `json:"pricing"`
	OfficialPricing *modelPlazaOfficialPricing `json:"official_pricing"`

	// 兼容本地旧版 /models 模型广场视图的扁平字段。新视图使用 Pricing/OfficialPricing。
	BillingMode       string   `json:"billing_mode"`
	InputPrice1M      *float64 `json:"input_price_1m"`
	OutputPrice1M     *float64 `json:"output_price_1m"`
	CacheWritePrice1M *float64 `json:"cache_write_price_1m"`
	CacheReadPrice1M  *float64 `json:"cache_read_price_1m"`
}

// modelPlazaGroup 广场分组条目（白名单字段）。
type modelPlazaGroup struct {
	ID                 int64             `json:"id"`
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	Platform           string            `json:"platform"`
	SubscriptionType   string            `json:"subscription_type"`
	RateMultiplier     float64           `json:"rate_multiplier"`
	UserRateMultiplier *float64          `json:"user_rate_multiplier,omitempty"`
	PeakRateEnabled    bool              `json:"peak_rate_enabled"`
	PeakStart          string            `json:"peak_start"`
	PeakEnd            string            `json:"peak_end"`
	PeakRateMultiplier float64           `json:"peak_rate_multiplier"`
	IsExclusive        bool              `json:"is_exclusive"`
	ModelCount         int               `json:"model_count"`
	Models             []modelPlazaModel `json:"models"`
}

// modelPlazaSummary / modelPlazaPlatform 兼容本地旧版 /models 视图。
type modelPlazaSummary struct {
	PlatformCount int `json:"platform_count"`
	GroupCount    int `json:"group_count"`
	ModelCount    int `json:"model_count"`
}

type modelPlazaPlatform struct {
	Platform   string            `json:"platform"`
	Label      string            `json:"label"`
	GroupCount int               `json:"group_count"`
	Groups     []modelPlazaGroup `json:"groups"`
}

// modelPlazaResponse 广场页响应。groups/description 是上游新视图口径；
// summary/platforms 用于兼容本地历史 /models 视图。
type modelPlazaResponse struct {
	Description string               `json:"description"`
	Groups      []modelPlazaGroup    `json:"groups"`
	Summary     modelPlazaSummary    `json:"summary"`
	Platforms   []modelPlazaPlatform `json:"platforms"`
}

// Get 返回模型广场数据。
// GET /api/v1/model-plaza
func (h *ModelPlazaHandler) Get(c *gin.Context) {
	if h.settingService == nil || h.channelService == nil {
		response.NotFound(c, "Model plaza is not enabled")
		return
	}
	rt := h.settingService.GetModelPlazaRuntime(c.Request.Context())
	if !rt.Enabled {
		response.NotFound(c, "Model plaza is not enabled")
		return
	}

	subject, authed := middleware.GetAuthSubjectFromContext(c)
	if rt.RequireAuth && !authed {
		response.Unauthorized(c, "Authentication required")
		return
	}

	groups, err := h.channelService.ListPlazaGroups(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// allowedExclusive == nil 表示匿名；登录用户恒为非 nil（可能为空集合）。
	var allowedExclusive map[int64]struct{}
	var userRates map[int64]float64
	if authed {
		allowedExclusive, err = h.apiKeyService.GetUserAllowedGroupIDSet(c.Request.Context(), subject.UserID)
		if err != nil {
			// 可见性数据拿不到时不能静默降级成匿名视图（会错漏专属分组），直接报错。
			response.ErrorFrom(c, err)
			return
		}
		userRates, err = h.apiKeyService.GetUserGroupRates(c.Request.Context(), subject.UserID)
		if err != nil {
			// 专属倍率仅是展示增强，失败降级为分组默认倍率。
			slog.Warn("model_plaza_user_rates_failed", "error", err, "user_id", subject.UserID)
			userRates = nil
		}
	}

	visible := filterPlazaVisibleGroups(groups, allowedExclusive)

	out := make([]modelPlazaGroup, 0, len(visible))
	for i := range visible {
		out = append(out, toModelPlazaGroupDTO(&visible[i], userRates))
	}
	platforms, summary := buildLegacyModelPlazaView(out)
	response.Success(c, modelPlazaResponse{
		Description: rt.Description,
		Groups:      out,
		Summary:     summary,
		Platforms:   platforms,
	})
}

// filterPlazaVisibleGroups 按登录态裁剪分组可见性。
// allowedExclusive == nil 表示匿名（仅非专属）；非 nil 表示登录（非专属 + 授权专属）。
func filterPlazaVisibleGroups(
	groups []service.PlazaGroup,
	allowedExclusive map[int64]struct{},
) []service.PlazaGroup {
	visible := make([]service.PlazaGroup, 0, len(groups))
	for _, g := range groups {
		if g.IsExclusive {
			if allowedExclusive == nil {
				continue
			}
			if _, ok := allowedExclusive[g.ID]; !ok {
				continue
			}
		}
		visible = append(visible, g)
	}
	return visible
}

// toModelPlazaGroupDTO 将 service 层广场分组映射为白名单 DTO，并合并用户专属倍率。
func toModelPlazaGroupDTO(g *service.PlazaGroup, userRates map[int64]float64) modelPlazaGroup {
	models := make([]modelPlazaModel, 0, len(g.Models))
	for i := range g.Models {
		m := &g.Models[i]
		pricing := toUserPricing(m.Pricing)
		official := toModelPlazaOfficialPricing(m.OfficialPricing)
		models = append(models, toModelPlazaModelDTO(m.Name, m.Platform, pricing, official))
	}
	dto := modelPlazaGroup{
		ID:                 g.ID,
		Name:               g.Name,
		Description:        g.Description,
		Platform:           g.Platform,
		SubscriptionType:   g.SubscriptionType,
		RateMultiplier:     g.RateMultiplier,
		PeakRateEnabled:    g.PeakRateEnabled,
		PeakStart:          g.PeakStart,
		PeakEnd:            g.PeakEnd,
		PeakRateMultiplier: g.PeakRateMultiplier,
		IsExclusive:        g.IsExclusive,
		ModelCount:         len(models),
		Models:             models,
	}
	if rate, ok := userRates[g.ID]; ok {
		dto.UserRateMultiplier = &rate
	}
	return dto
}

func toModelPlazaModelDTO(
	name string,
	platform string,
	pricing *userSupportedModelPricing,
	official *modelPlazaOfficialPricing,
) modelPlazaModel {
	model := modelPlazaModel{
		Name:            name,
		Platform:        platform,
		Pricing:         pricing,
		OfficialPricing: official,
		BillingMode:     string(service.BillingModeToken),
	}
	if pricing != nil && pricing.BillingMode != "" {
		model.BillingMode = pricing.BillingMode
	}

	model.InputPrice1M = scaledModelPlazaPrice(firstModelPlazaPrice(
		func() *float64 {
			if pricing == nil {
				return nil
			}
			return pricing.InputPrice
		}(),
		func() *float64 {
			if official == nil {
				return nil
			}
			return official.InputPrice
		}(),
	))
	model.OutputPrice1M = scaledModelPlazaPrice(firstModelPlazaPrice(
		func() *float64 {
			if pricing == nil {
				return nil
			}
			return pricing.OutputPrice
		}(),
		func() *float64 {
			if official == nil {
				return nil
			}
			return official.OutputPrice
		}(),
	))
	model.CacheWritePrice1M = scaledModelPlazaPrice(firstModelPlazaPrice(
		func() *float64 {
			if pricing == nil {
				return nil
			}
			return pricing.CacheWritePrice
		}(),
		func() *float64 {
			if official == nil {
				return nil
			}
			return firstModelPlazaPrice(official.CacheWritePrice, official.CacheWrite1hPrice)
		}(),
	))
	model.CacheReadPrice1M = scaledModelPlazaPrice(firstModelPlazaPrice(
		func() *float64 {
			if pricing == nil {
				return nil
			}
			return pricing.CacheReadPrice
		}(),
		func() *float64 {
			if official == nil {
				return nil
			}
			return official.CacheReadPrice
		}(),
	))
	return model
}

// toModelPlazaOfficialPricing 转换官方参考价；nil 透传（前端显示 "-"）。
func toModelPlazaOfficialPricing(p *service.PlazaOfficialPricing) *modelPlazaOfficialPricing {
	if p == nil {
		return nil
	}
	return &modelPlazaOfficialPricing{
		InputPrice:        p.InputPrice,
		OutputPrice:       p.OutputPrice,
		CacheWritePrice:   p.CacheWritePrice,
		CacheWrite1hPrice: p.CacheWrite1hPrice,
		CacheReadPrice:    p.CacheReadPrice,
	}
}

func buildLegacyModelPlazaView(groups []modelPlazaGroup) ([]modelPlazaPlatform, modelPlazaSummary) {
	buckets := make(map[string][]modelPlazaGroup)
	totalModels := 0
	for _, group := range groups {
		buckets[group.Platform] = append(buckets[group.Platform], group)
		totalModels += group.ModelCount
	}

	platforms := make([]modelPlazaPlatform, 0, len(buckets))
	visited := make(map[string]struct{}, len(buckets))
	for _, platform := range []string{
		service.PlatformOpenAI,
		service.PlatformAnthropic,
		service.PlatformGemini,
		service.PlatformAntigravity,
		service.PlatformGrok,
		service.PlatformComposite,
	} {
		if groupsForPlatform, ok := buckets[platform]; ok {
			platforms = append(platforms, toLegacyModelPlazaPlatform(platform, groupsForPlatform))
			visited[platform] = struct{}{}
		}
	}

	unknown := make([]string, 0)
	for platform := range buckets {
		if _, ok := visited[platform]; !ok {
			unknown = append(unknown, platform)
		}
	}
	sort.Strings(unknown)
	for _, platform := range unknown {
		platforms = append(platforms, toLegacyModelPlazaPlatform(platform, buckets[platform]))
	}

	return platforms, modelPlazaSummary{
		PlatformCount: len(platforms),
		GroupCount:    len(groups),
		ModelCount:    totalModels,
	}
}

func toLegacyModelPlazaPlatform(platform string, groups []modelPlazaGroup) modelPlazaPlatform {
	sortedGroups := append([]modelPlazaGroup(nil), groups...)
	sort.SliceStable(sortedGroups, func(i, j int) bool {
		return strings.ToLower(sortedGroups[i].Name) < strings.ToLower(sortedGroups[j].Name)
	})
	return modelPlazaPlatform{
		Platform:   platform,
		Label:      modelPlazaPlatformLabel(platform),
		GroupCount: len(sortedGroups),
		Groups:     sortedGroups,
	}
}

func modelPlazaPlatformLabel(platform string) string {
	switch platform {
	case service.PlatformOpenAI:
		return "OpenAI"
	case service.PlatformAnthropic:
		return "Anthropic"
	case service.PlatformGemini:
		return "Gemini"
	case service.PlatformAntigravity:
		return "Antigravity"
	case service.PlatformGrok:
		return "Grok"
	case service.PlatformComposite:
		return "Composite"
	default:
		return platform
	}
}

func firstModelPlazaPrice(values ...*float64) *float64 {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func scaledModelPlazaPrice(value *float64) *float64 {
	if value == nil {
		return nil
	}
	scaled := *value * modelPlazaLegacyPriceScale
	return &scaled
}
