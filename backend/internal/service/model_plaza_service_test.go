package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type modelPlazaGroupRepoStub struct {
	listActiveFn func(ctx context.Context) ([]Group, error)
}

func (s modelPlazaGroupRepoStub) ListActive(ctx context.Context) ([]Group, error) {
	if s.listActiveFn != nil {
		return s.listActiveFn(ctx)
	}
	return nil, nil
}

type modelPlazaChannelReaderStub struct {
	getChannelForGroupFn func(ctx context.Context, groupID int64) (*Channel, error)
}

func (s modelPlazaChannelReaderStub) GetChannelForGroup(ctx context.Context, groupID int64) (*Channel, error) {
	if s.getChannelForGroupFn != nil {
		return s.getChannelForGroupFn(ctx, groupID)
	}
	return nil, nil
}

type modelPlazaAccountReaderStub struct {
	listSchedulableByGroupIDFn func(ctx context.Context, groupID int64) ([]Account, error)
}

func (s modelPlazaAccountReaderStub) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]Account, error) {
	if s.listSchedulableByGroupIDFn != nil {
		return s.listSchedulableByGroupIDFn(ctx, groupID)
	}
	return nil, nil
}

type modelPlazaBillingReaderStub struct {
	getModelPricingFn func(model string) (*ModelPricing, error)
}

func (s modelPlazaBillingReaderStub) GetModelPricing(model string) (*ModelPricing, error) {
	if s.getModelPricingFn != nil {
		return s.getModelPricingFn(model)
	}
	return nil, nil
}

func TestModelPlazaService_Get_FiltersPublicGroupsAndConvertsToPerMillion(t *testing.T) {
	t.Parallel()

	svc := NewModelPlazaService(
		modelPlazaGroupRepoStub{
			listActiveFn: func(context.Context) ([]Group, error) {
				return []Group{
					{ID: 1, Name: "OpenAI Public", Platform: PlatformOpenAI, Status: StatusActive, IsExclusive: false, RateMultiplier: 1.2},
					{ID: 2, Name: "Private VIP", Platform: PlatformOpenAI, Status: StatusActive, IsExclusive: true, RateMultiplier: 0.9},
				}, nil
			},
		},
		modelPlazaChannelReaderStub{
			getChannelForGroupFn: func(context.Context, int64) (*Channel, error) {
				return &Channel{
					GroupIDs: []int64{1},
					ModelPricing: []ChannelModelPricing{
						{Platform: PlatformOpenAI, Models: []string{"gpt-5.4"}},
					},
				}, nil
			},
		},
		modelPlazaAccountReaderStub{},
		modelPlazaBillingReaderStub{
			getModelPricingFn: func(model string) (*ModelPricing, error) {
				require.Equal(t, "gpt-5.4", model)
				return &ModelPricing{
					InputPricePerToken:  2.5e-6,
					OutputPricePerToken: 15e-6,
				}, nil
			},
		},
	)

	got, err := svc.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, got.Summary.PlatformCount)
	require.Equal(t, 1, got.Summary.GroupCount)
	require.Equal(t, 1, got.Summary.ModelCount)
	require.Len(t, got.Platforms, 1)
	require.Equal(t, "OpenAI", got.Platforms[0].Label)
	require.Len(t, got.Platforms[0].Groups, 1)
	require.Equal(t, "gpt-5.4", got.Platforms[0].Groups[0].Models[0].Name)
	require.InDelta(t, 2.5, derefFloat64(got.Platforms[0].Groups[0].Models[0].InputPrice1M), 1e-9)
	require.InDelta(t, 15.0, derefFloat64(got.Platforms[0].Groups[0].Models[0].OutputPrice1M), 1e-9)
}

func TestModelPlazaService_Get_KeepsGroupWhenNoModels(t *testing.T) {
	t.Parallel()

	svc := NewModelPlazaService(
		modelPlazaGroupRepoStub{
			listActiveFn: func(context.Context) ([]Group, error) {
				return []Group{
					{ID: 10, Name: "Empty Public", Platform: PlatformGemini, Status: StatusActive, IsExclusive: false, RateMultiplier: 1},
				}, nil
			},
		},
		modelPlazaChannelReaderStub{},
		modelPlazaAccountReaderStub{
			listSchedulableByGroupIDFn: func(context.Context, int64) ([]Account, error) {
				return nil, nil
			},
		},
		modelPlazaBillingReaderStub{},
	)

	got, err := svc.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, got.Summary.PlatformCount)
	require.Equal(t, 1, got.Summary.GroupCount)
	require.Equal(t, 0, got.Summary.ModelCount)
	require.Len(t, got.Platforms, 1)
	require.Len(t, got.Platforms[0].Groups, 1)
	require.Equal(t, 0, got.Platforms[0].Groups[0].ModelCount)
	require.Empty(t, got.Platforms[0].Groups[0].Models)
}

func TestModelPlazaService_Get_LeavesPriceNilWhenPricingMissing(t *testing.T) {
	t.Parallel()

	svc := NewModelPlazaService(
		modelPlazaGroupRepoStub{
			listActiveFn: func(context.Context) ([]Group, error) {
				return []Group{
					{ID: 11, Name: "OpenAI Public", Platform: PlatformOpenAI, Status: StatusActive, IsExclusive: false, RateMultiplier: 1},
				}, nil
			},
		},
		modelPlazaChannelReaderStub{
			getChannelForGroupFn: func(context.Context, int64) (*Channel, error) {
				return &Channel{
					GroupIDs: []int64{11},
					ModelPricing: []ChannelModelPricing{
						{Platform: PlatformOpenAI, Models: []string{"unknown-model"}},
					},
				}, nil
			},
		},
		modelPlazaAccountReaderStub{},
		modelPlazaBillingReaderStub{
			getModelPricingFn: func(string) (*ModelPricing, error) {
				return nil, errors.New("missing pricing")
			},
		},
	)

	got, err := svc.Get(context.Background())
	require.NoError(t, err)
	require.Len(t, got.Platforms, 1)
	require.Len(t, got.Platforms[0].Groups, 1)
	require.Len(t, got.Platforms[0].Groups[0].Models, 1)
	require.Equal(t, "unknown-model", got.Platforms[0].Groups[0].Models[0].Name)
	require.Nil(t, got.Platforms[0].Groups[0].Models[0].InputPrice1M)
	require.Nil(t, got.Platforms[0].Groups[0].Models[0].OutputPrice1M)
}

func TestModelPlazaService_Get_FallsBackToAccountsWhenChannelUsesWildcard(t *testing.T) {
	t.Parallel()

	svc := NewModelPlazaService(
		modelPlazaGroupRepoStub{
			listActiveFn: func(context.Context) ([]Group, error) {
				return []Group{
					{ID: 12, Name: "OpenAI Public", Platform: PlatformOpenAI, Status: StatusActive, IsExclusive: false, RateMultiplier: 1},
				}, nil
			},
		},
		modelPlazaChannelReaderStub{
			getChannelForGroupFn: func(context.Context, int64) (*Channel, error) {
				return &Channel{
					GroupIDs: []int64{12},
					ModelPricing: []ChannelModelPricing{
						{Platform: PlatformOpenAI, Models: []string{"gpt-*"}},
					},
				}, nil
			},
		},
		modelPlazaAccountReaderStub{
			listSchedulableByGroupIDFn: func(context.Context, int64) ([]Account, error) {
				return []Account{
					{
						Platform: PlatformOpenAI,
						Credentials: map[string]any{
							"model_mapping": map[string]any{
								"gpt-4.1":      "gpt-4.1",
								"gpt-4.1-mini": "gpt-4.1-mini",
							},
						},
					},
				}, nil
			},
		},
		modelPlazaBillingReaderStub{},
	)

	got, err := svc.Get(context.Background())
	require.NoError(t, err)
	require.Len(t, got.Platforms, 1)
	require.Len(t, got.Platforms[0].Groups, 1)
	require.Len(t, got.Platforms[0].Groups[0].Models, 2)
	require.Equal(t, "gpt-4.1", got.Platforms[0].Groups[0].Models[0].Name)
	require.Equal(t, "gpt-4.1-mini", got.Platforms[0].Groups[0].Models[1].Name)
}

func derefFloat64(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}
