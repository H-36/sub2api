//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type redeemClaimsRepoStub struct {
	redeemRepoStub
	claimsByCodeID  map[int64][]RedeemCodeClaim
	listClaimsCalls []int64
}

func (s *redeemClaimsRepoStub) ListClaimsByRedeemCode(ctx context.Context, redeemCodeID int64) ([]RedeemCodeClaim, error) {
	s.listClaimsCalls = append(s.listClaimsCalls, redeemCodeID)
	return append([]RedeemCodeClaim(nil), s.claimsByCodeID[redeemCodeID]...), nil
}

func TestAdminService_GetRedeemCodeClaims_ForWelfareCode(t *testing.T) {
	now := time.Date(2026, time.April, 17, 9, 30, 0, 0, time.UTC)
	redeemRepo := &redeemClaimsRepoStub{
		redeemRepoStub: redeemRepoStub{
			codeByID: map[int64]*RedeemCode{
				7: {
					ID:   7,
					Type: RedeemTypeWelfare,
				},
			},
		},
		claimsByCodeID: map[int64][]RedeemCodeClaim{
			7: {
				{
					ID:           1,
					RedeemCodeID: 7,
					UserID:       88,
					Amount:       10,
					ClaimedAt:    now,
				},
			},
		},
	}
	svc := &adminServiceImpl{redeemCodeRepo: redeemRepo}

	claims, err := svc.GetRedeemCodeClaims(context.Background(), 7)
	require.NoError(t, err)
	require.Len(t, claims, 1)
	assert.Equal(t, []int64{7}, redeemRepo.listClaimsCalls)
	assert.Equal(t, int64(88), claims[0].UserID)
	assert.Equal(t, 10.0, claims[0].Amount)
	assert.Equal(t, now, claims[0].ClaimedAt)
}

func TestAdminService_GetRedeemCodeClaims_SkipsNonWelfareCode(t *testing.T) {
	redeemRepo := &redeemClaimsRepoStub{
		redeemRepoStub: redeemRepoStub{
			codeByID: map[int64]*RedeemCode{
				9: {
					ID:   9,
					Type: RedeemTypeBalance,
				},
			},
		},
	}
	svc := &adminServiceImpl{redeemCodeRepo: redeemRepo}

	claims, err := svc.GetRedeemCodeClaims(context.Background(), 9)
	require.NoError(t, err)
	assert.Empty(t, claims)
	assert.Empty(t, redeemRepo.listClaimsCalls)
}
