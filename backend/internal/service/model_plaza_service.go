package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

type modelPlazaGroupReader interface {
	ListActive(ctx context.Context) ([]Group, error)
}

type modelPlazaChannelReader interface {
	GetChannelForGroup(ctx context.Context, groupID int64) (*Channel, error)
}

type modelPlazaAccountReader interface {
	ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]Account, error)
}

type modelPlazaBillingReader interface {
	GetModelPricing(model string) (*ModelPricing, error)
}

type ModelPlazaSummary struct {
	PlatformCount int `json:"platform_count"`
	GroupCount    int `json:"group_count"`
	ModelCount    int `json:"model_count"`
}

type ModelPlazaModel struct {
	Name              string   `json:"name"`
	BillingMode       string   `json:"billing_mode"`
	InputPrice1M      *float64 `json:"input_price_1m"`
	OutputPrice1M     *float64 `json:"output_price_1m"`
	CacheWritePrice1M *float64 `json:"cache_write_price_1m"`
	CacheReadPrice1M  *float64 `json:"cache_read_price_1m"`
}

type ModelPlazaGroup struct {
	ID             int64             `json:"id"`
	Name           string            `json:"name"`
	Platform       string            `json:"platform"`
	RateMultiplier float64           `json:"rate_multiplier"`
	ModelCount     int               `json:"model_count"`
	Models         []ModelPlazaModel `json:"models"`
}

type ModelPlazaPlatform struct {
	Platform   string            `json:"platform"`
	Label      string            `json:"label"`
	GroupCount int               `json:"group_count"`
	Groups     []ModelPlazaGroup `json:"groups"`
}

type ModelPlazaResponse struct {
	Summary   ModelPlazaSummary    `json:"summary"`
	Platforms []ModelPlazaPlatform `json:"platforms"`
}

type modelPlazaModelSeed struct {
	Name        string
	BillingMode BillingMode
}

var modelPlazaPlatformOrder = []string{
	PlatformOpenAI,
	PlatformAnthropic,
	PlatformGemini,
	PlatformAntigravity,
}

var modelPlazaPlatformLabels = map[string]string{
	PlatformOpenAI:      "OpenAI",
	PlatformAnthropic:   "Anthropic",
	PlatformGemini:      "Gemini",
	PlatformAntigravity: "Antigravity",
}

type ModelPlazaService struct {
	groupReader   modelPlazaGroupReader
	channelReader modelPlazaChannelReader
	accountReader modelPlazaAccountReader
	billingReader modelPlazaBillingReader
}

func NewModelPlazaService(
	groupReader modelPlazaGroupReader,
	channelReader modelPlazaChannelReader,
	accountReader modelPlazaAccountReader,
	billingReader modelPlazaBillingReader,
) *ModelPlazaService {
	return &ModelPlazaService{
		groupReader:   groupReader,
		channelReader: channelReader,
		accountReader: accountReader,
		billingReader: billingReader,
	}
}

func (s *ModelPlazaService) Get(ctx context.Context) (*ModelPlazaResponse, error) {
	groups, err := s.groupReader.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active groups: %w", err)
	}

	groupBuckets := make(map[string][]ModelPlazaGroup)
	totalGroups := 0
	totalModels := 0

	for _, group := range groups {
		if group.Status != StatusActive || group.IsExclusive {
			continue
		}

		groupView, err := s.buildGroup(ctx, group)
		if err != nil {
			return nil, err
		}

		groupBuckets[group.Platform] = append(groupBuckets[group.Platform], groupView)
		totalGroups++
		totalModels += groupView.ModelCount
	}

	platforms := buildModelPlazaPlatforms(groupBuckets)
	return &ModelPlazaResponse{
		Summary: ModelPlazaSummary{
			PlatformCount: len(platforms),
			GroupCount:    totalGroups,
			ModelCount:    totalModels,
		},
		Platforms: platforms,
	}, nil
}

func (s *ModelPlazaService) buildGroup(ctx context.Context, group Group) (ModelPlazaGroup, error) {
	seeds, err := s.collectGroupModelSeeds(ctx, group)
	if err != nil {
		return ModelPlazaGroup{}, fmt.Errorf("build group %d models: %w", group.ID, err)
	}

	models := make([]ModelPlazaModel, 0, len(seeds))
	for _, seed := range seeds {
		models = append(models, s.buildModel(seed))
	}

	return ModelPlazaGroup{
		ID:             group.ID,
		Name:           group.Name,
		Platform:       group.Platform,
		RateMultiplier: group.RateMultiplier,
		ModelCount:     len(models),
		Models:         models,
	}, nil
}

func (s *ModelPlazaService) collectGroupModelSeeds(ctx context.Context, group Group) ([]modelPlazaModelSeed, error) {
	if s.channelReader != nil {
		channel, err := s.channelReader.GetChannelForGroup(ctx, group.ID)
		if err == nil && channel != nil {
			seeds := collectChannelModelSeeds(channel, group.Platform)
			if len(seeds) > 0 {
				return seeds, nil
			}
		}
	}

	if s.accountReader == nil {
		return nil, nil
	}

	accounts, err := s.accountReader.ListSchedulableByGroupID(ctx, group.ID)
	if err != nil {
		return nil, fmt.Errorf("list schedulable accounts for group %d: %w", group.ID, err)
	}
	return collectAccountModelSeeds(accounts, group.Platform), nil
}

func (s *ModelPlazaService) buildModel(seed modelPlazaModelSeed) ModelPlazaModel {
	mode := normalizeModelPlazaBillingMode(seed.BillingMode)
	model := ModelPlazaModel{
		Name:        seed.Name,
		BillingMode: string(mode),
	}

	if mode != BillingModeToken || s.billingReader == nil {
		return model
	}

	pricing, err := s.billingReader.GetModelPricing(seed.Name)
	if err != nil || pricing == nil {
		return model
	}

	model.InputPrice1M = nonZeroFloat64Ptr(pricing.InputPricePerToken * 1_000_000)
	model.OutputPrice1M = nonZeroFloat64Ptr(pricing.OutputPricePerToken * 1_000_000)
	model.CacheWritePrice1M = nonZeroFloat64Ptr(modelPlazaCacheWritePrice(pricing) * 1_000_000)
	model.CacheReadPrice1M = nonZeroFloat64Ptr(pricing.CacheReadPricePerToken * 1_000_000)

	return model
}

func buildModelPlazaPlatforms(groupBuckets map[string][]ModelPlazaGroup) []ModelPlazaPlatform {
	platforms := make([]ModelPlazaPlatform, 0, len(groupBuckets))
	for _, platform := range modelPlazaPlatformOrder {
		groups := groupBuckets[platform]
		if len(groups) == 0 {
			continue
		}

		sort.Slice(groups, func(i, j int) bool {
			return strings.ToLower(groups[i].Name) < strings.ToLower(groups[j].Name)
		})

		platforms = append(platforms, ModelPlazaPlatform{
			Platform:   platform,
			Label:      modelPlazaPlatformLabel(platform),
			GroupCount: len(groups),
			Groups:     groups,
		})
	}
	return platforms
}

func collectChannelModelSeeds(channel *Channel, platform string) []modelPlazaModelSeed {
	if channel == nil {
		return nil
	}

	seen := make(map[string]struct{})
	seeds := make([]modelPlazaModelSeed, 0)
	for _, pricing := range channel.ModelPricing {
		if pricing.Platform != platform {
			continue
		}
		mode := normalizeModelPlazaBillingMode(pricing.BillingMode)
		for _, model := range pricing.Models {
			name := strings.TrimSpace(model)
			if !shouldExposeModelName(name) {
				continue
			}
			key := strings.ToLower(name)
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			seeds = append(seeds, modelPlazaModelSeed{Name: name, BillingMode: mode})
		}
	}

	sort.Slice(seeds, func(i, j int) bool {
		return strings.ToLower(seeds[i].Name) < strings.ToLower(seeds[j].Name)
	})
	return seeds
}

func collectAccountModelSeeds(accounts []Account, platform string) []modelPlazaModelSeed {
	seen := make(map[string]struct{})
	seeds := make([]modelPlazaModelSeed, 0)
	for _, account := range accounts {
		if account.Platform != platform {
			continue
		}

		for model := range account.GetModelMapping() {
			name := strings.TrimSpace(model)
			if !shouldExposeModelName(name) {
				continue
			}
			key := strings.ToLower(name)
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			seeds = append(seeds, modelPlazaModelSeed{Name: name, BillingMode: BillingModeToken})
		}
	}

	sort.Slice(seeds, func(i, j int) bool {
		return strings.ToLower(seeds[i].Name) < strings.ToLower(seeds[j].Name)
	})
	return seeds
}

func modelPlazaPlatformLabel(platform string) string {
	if label, ok := modelPlazaPlatformLabels[platform]; ok {
		return label
	}
	return platform
}

func normalizeModelPlazaBillingMode(mode BillingMode) BillingMode {
	switch mode {
	case BillingModePerRequest, BillingModeImage:
		return mode
	default:
		return BillingModeToken
	}
}

func nonZeroFloat64Ptr(value float64) *float64 {
	if value == 0 {
		return nil
	}
	return float64Ptr(value)
}

func modelPlazaCacheWritePrice(pricing *ModelPricing) float64 {
	if pricing == nil {
		return 0
	}
	if pricing.CacheCreationPricePerToken > 0 {
		return pricing.CacheCreationPricePerToken
	}
	if pricing.CacheCreation5mPrice > 0 {
		return pricing.CacheCreation5mPrice
	}
	return pricing.CacheCreation1hPrice
}

func shouldExposeModelName(name string) bool {
	name = strings.TrimSpace(name)
	return name != "" && !strings.Contains(name, "*")
}
