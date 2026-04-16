//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type redeemRepoStubForGenerate struct {
	created     *RedeemCode
	createErr   error
	createCalls int
}

func (s *redeemRepoStubForGenerate) Create(_ context.Context, code *RedeemCode) error {
	s.createCalls++
	copied := *code
	s.created = &copied
	return s.createErr
}

func (s *redeemRepoStubForGenerate) CreateBatch(context.Context, []RedeemCode) error {
	panic("unexpected CreateBatch call")
}

func (s *redeemRepoStubForGenerate) GetByID(context.Context, int64) (*RedeemCode, error) {
	panic("unexpected GetByID call")
}

func (s *redeemRepoStubForGenerate) GetByCode(context.Context, string) (*RedeemCode, error) {
	panic("unexpected GetByCode call")
}

func (s *redeemRepoStubForGenerate) Update(context.Context, *RedeemCode) error {
	panic("unexpected Update call")
}

func (s *redeemRepoStubForGenerate) Delete(context.Context, int64) error {
	panic("unexpected Delete call")
}

func (s *redeemRepoStubForGenerate) Use(context.Context, int64, int64) error {
	panic("unexpected Use call")
}

func (s *redeemRepoStubForGenerate) HasClaimByUser(context.Context, int64, int64) (bool, error) {
	panic("unexpected HasClaimByUser call")
}

func (s *redeemRepoStubForGenerate) CreateClaim(context.Context, int64, int64, float64) error {
	panic("unexpected CreateClaim call")
}

func (s *redeemRepoStubForGenerate) IncrementClaimedCount(context.Context, int64, int64) (int, error) {
	panic("unexpected IncrementClaimedCount call")
}

func (s *redeemRepoStubForGenerate) ListClaimsByRedeemCode(context.Context, int64) ([]RedeemCodeClaim, error) {
	panic("unexpected ListClaimsByRedeemCode call")
}

func (s *redeemRepoStubForGenerate) List(context.Context, pagination.PaginationParams) ([]RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *redeemRepoStubForGenerate) ListWithFilters(context.Context, pagination.PaginationParams, string, string, string) ([]RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
}

func (s *redeemRepoStubForGenerate) ListByUser(context.Context, int64, int) ([]RedeemCode, error) {
	panic("unexpected ListByUser call")
}

func (s *redeemRepoStubForGenerate) ListByUserPaginated(context.Context, int64, pagination.PaginationParams, string) ([]RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected ListByUserPaginated call")
}

func (s *redeemRepoStubForGenerate) SumPositiveBalanceByUser(context.Context, int64) (float64, error) {
	panic("unexpected SumPositiveBalanceByUser call")
}

func TestAdminService_GenerateRedeemCodes_UsesCustomCode(t *testing.T) {
	redeemRepo := &redeemRepoStubForGenerate{}
	svc := &adminServiceImpl{
		redeemCodeRepo: redeemRepo,
		groupRepo:      &groupRepoStubForAdmin{},
	}

	codes, err := svc.GenerateRedeemCodes(context.Background(), &GenerateRedeemCodesInput{
		Count: 1,
		Code:  "  CUSTOM-REDEEM-1  ",
		Type:  RedeemTypeBalance,
		Value: 25,
	})
	require.NoError(t, err)
	require.Len(t, codes, 1)
	require.NotNil(t, redeemRepo.created)
	require.Equal(t, "CUSTOM-REDEEM-1", redeemRepo.created.Code)
	require.Equal(t, 1, redeemRepo.createCalls)
	require.Equal(t, "CUSTOM-REDEEM-1", codes[0].Code)
}

func TestAdminService_GenerateRedeemCodes_RejectsCustomCodeWithMultiCount(t *testing.T) {
	redeemRepo := &redeemRepoStubForGenerate{}
	svc := &adminServiceImpl{
		redeemCodeRepo: redeemRepo,
		groupRepo:      &groupRepoStubForAdmin{},
	}

	_, err := svc.GenerateRedeemCodes(context.Background(), &GenerateRedeemCodesInput{
		Count: 2,
		Code:  "CUSTOM-REDEEM-2",
		Type:  RedeemTypeBalance,
		Value: 25,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "count must be 1")
	require.Nil(t, redeemRepo.created)
}

func TestAdminService_GenerateRedeemCodes_CreatesWelfareCode(t *testing.T) {
	redeemRepo := &redeemRepoStubForGenerate{}
	svc := &adminServiceImpl{
		redeemCodeRepo: redeemRepo,
		groupRepo:      &groupRepoStubForAdmin{},
	}

	codes, err := svc.GenerateRedeemCodes(context.Background(), &GenerateRedeemCodesInput{
		Count: 20,
		Code:  "233",
		Type:  RedeemTypeWelfare,
		Value: 10,
	})
	require.NoError(t, err)
	require.Len(t, codes, 1)
	require.NotNil(t, redeemRepo.created)
	require.Equal(t, RedeemTypeWelfare, redeemRepo.created.Type)
	require.Equal(t, 20, redeemRepo.created.MaxClaims)
	require.Equal(t, 0, redeemRepo.created.ClaimedCount)
	require.Equal(t, "233", redeemRepo.created.Code)
}
