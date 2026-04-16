package service_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newWelfareRedeemServiceTestClient(t *testing.T) *dbent.Client {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name())
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() })

	return client
}

type welfareRedeemRepoStub struct {
	nextID   int64
	codes    map[int64]*service.RedeemCode
	codeToID map[string]int64
	claims   map[int64]map[int64]struct{}
}

func newWelfareRedeemRepoStub() *welfareRedeemRepoStub {
	return &welfareRedeemRepoStub{
		nextID:   1,
		codes:    make(map[int64]*service.RedeemCode),
		codeToID: make(map[string]int64),
		claims:   make(map[int64]map[int64]struct{}),
	}
}

func (s *welfareRedeemRepoStub) Create(_ context.Context, code *service.RedeemCode) error {
	cloned := cloneRedeemCode(code)
	if cloned.ID == 0 {
		cloned.ID = s.nextID
		s.nextID++
	}
	code.ID = cloned.ID
	s.codes[cloned.ID] = cloned
	s.codeToID[cloned.Code] = cloned.ID
	return nil
}

func (s *welfareRedeemRepoStub) CreateBatch(context.Context, []service.RedeemCode) error {
	panic("unexpected CreateBatch call")
}

func (s *welfareRedeemRepoStub) GetByID(_ context.Context, id int64) (*service.RedeemCode, error) {
	code, ok := s.codes[id]
	if !ok {
		return nil, service.ErrRedeemCodeNotFound
	}
	return cloneRedeemCode(code), nil
}

func (s *welfareRedeemRepoStub) GetByCode(_ context.Context, code string) (*service.RedeemCode, error) {
	id, ok := s.codeToID[code]
	if !ok {
		return nil, service.ErrRedeemCodeNotFound
	}
	return s.GetByID(context.Background(), id)
}

func (s *welfareRedeemRepoStub) Update(_ context.Context, code *service.RedeemCode) error {
	if _, ok := s.codes[code.ID]; !ok {
		return service.ErrRedeemCodeNotFound
	}
	s.codes[code.ID] = cloneRedeemCode(code)
	return nil
}

func (s *welfareRedeemRepoStub) Delete(context.Context, int64) error {
	panic("unexpected Delete call")
}

func (s *welfareRedeemRepoStub) Use(context.Context, int64, int64) error {
	panic("unexpected Use call")
}

func (s *welfareRedeemRepoStub) HasClaimByUser(_ context.Context, redeemCodeID, userID int64) (bool, error) {
	claims := s.claims[redeemCodeID]
	_, ok := claims[userID]
	return ok, nil
}

func (s *welfareRedeemRepoStub) CreateClaim(_ context.Context, redeemCodeID, userID int64, amount float64) error {
	if _, ok := s.claims[redeemCodeID]; !ok {
		s.claims[redeemCodeID] = make(map[int64]struct{})
	}
	if _, ok := s.claims[redeemCodeID][userID]; ok {
		return service.ErrRedeemCodeClaimed
	}
	s.claims[redeemCodeID][userID] = struct{}{}
	return nil
}

func (s *welfareRedeemRepoStub) IncrementClaimedCount(_ context.Context, id, maxClaims int64) (int, error) {
	code, ok := s.codes[id]
	if !ok {
		return 0, service.ErrRedeemCodeNotFound
	}
	if code.Status == service.StatusUsed || int64(code.ClaimedCount) >= maxClaims {
		return 0, service.ErrRedeemCodeUsed
	}
	code.ClaimedCount++
	return code.ClaimedCount, nil
}

func (s *welfareRedeemRepoStub) ListClaimsByRedeemCode(context.Context, int64) ([]service.RedeemCodeClaim, error) {
	panic("unexpected ListClaimsByRedeemCode call")
}

func (s *welfareRedeemRepoStub) List(context.Context, pagination.PaginationParams) ([]service.RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *welfareRedeemRepoStub) ListWithFilters(context.Context, pagination.PaginationParams, string, string, string) ([]service.RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
}

func (s *welfareRedeemRepoStub) ListByUser(context.Context, int64, int) ([]service.RedeemCode, error) {
	panic("unexpected ListByUser call")
}

func (s *welfareRedeemRepoStub) ListByUserPaginated(context.Context, int64, pagination.PaginationParams, string) ([]service.RedeemCode, *pagination.PaginationResult, error) {
	panic("unexpected ListByUserPaginated call")
}

func (s *welfareRedeemRepoStub) SumPositiveBalanceByUser(context.Context, int64) (float64, error) {
	panic("unexpected SumPositiveBalanceByUser call")
}

type welfareUserRepoStub struct {
	users map[int64]*service.User
}

func (s *welfareUserRepoStub) Create(context.Context, *service.User) error {
	panic("unexpected Create call")
}

func (s *welfareUserRepoStub) GetByID(_ context.Context, id int64) (*service.User, error) {
	user, ok := s.users[id]
	if !ok {
		return nil, service.ErrUserNotFound
	}
	return cloneUser(user), nil
}

func (s *welfareUserRepoStub) GetByEmail(context.Context, string) (*service.User, error) {
	panic("unexpected GetByEmail call")
}

func (s *welfareUserRepoStub) GetFirstAdmin(context.Context) (*service.User, error) {
	panic("unexpected GetFirstAdmin call")
}

func (s *welfareUserRepoStub) Update(context.Context, *service.User) error {
	panic("unexpected Update call")
}

func (s *welfareUserRepoStub) Delete(context.Context, int64) error {
	panic("unexpected Delete call")
}

func (s *welfareUserRepoStub) List(context.Context, pagination.PaginationParams) ([]service.User, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *welfareUserRepoStub) ListWithFilters(context.Context, pagination.PaginationParams, service.UserListFilters) ([]service.User, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
}

func (s *welfareUserRepoStub) UpdateBalance(_ context.Context, id int64, amount float64) error {
	user, ok := s.users[id]
	if !ok {
		return service.ErrUserNotFound
	}
	user.Balance += amount
	return nil
}

func (s *welfareUserRepoStub) DeductBalance(context.Context, int64, float64) error {
	panic("unexpected DeductBalance call")
}

func (s *welfareUserRepoStub) UpdateConcurrency(context.Context, int64, int) error {
	panic("unexpected UpdateConcurrency call")
}

func (s *welfareUserRepoStub) ExistsByEmail(context.Context, string) (bool, error) {
	panic("unexpected ExistsByEmail call")
}

func (s *welfareUserRepoStub) RemoveGroupFromAllowedGroups(context.Context, int64) (int64, error) {
	panic("unexpected RemoveGroupFromAllowedGroups call")
}

func (s *welfareUserRepoStub) AddGroupToAllowedGroups(context.Context, int64, int64) error {
	panic("unexpected AddGroupToAllowedGroups call")
}

func (s *welfareUserRepoStub) RemoveGroupFromUserAllowedGroups(context.Context, int64, int64) error {
	panic("unexpected RemoveGroupFromUserAllowedGroups call")
}

func (s *welfareUserRepoStub) UpdateTotpSecret(context.Context, int64, *string) error {
	panic("unexpected UpdateTotpSecret call")
}

func (s *welfareUserRepoStub) EnableTotp(context.Context, int64) error {
	panic("unexpected EnableTotp call")
}

func (s *welfareUserRepoStub) DisableTotp(context.Context, int64) error {
	panic("unexpected DisableTotp call")
}

func cloneRedeemCode(code *service.RedeemCode) *service.RedeemCode {
	if code == nil {
		return nil
	}
	cloned := *code
	return &cloned
}

func cloneUser(user *service.User) *service.User {
	if user == nil {
		return nil
	}
	cloned := *user
	cloned.AllowedGroups = append([]int64(nil), user.AllowedGroups...)
	return &cloned
}

func TestRedeemService_WelfareRedeemFlow(t *testing.T) {
	ctx := context.Background()
	client := newWelfareRedeemServiceTestClient(t)

	redeemRepo := newWelfareRedeemRepoStub()
	userRepo := &welfareUserRepoStub{
		users: map[int64]*service.User{
			1: {ID: 1, Email: "welfare-1@test.com", Role: service.RoleUser, Status: service.StatusActive},
			2: {ID: 2, Email: "welfare-2@test.com", Role: service.RoleUser, Status: service.StatusActive},
			3: {ID: 3, Email: "welfare-3@test.com", Role: service.RoleUser, Status: service.StatusActive},
		},
	}
	redeemSvc := service.NewRedeemService(redeemRepo, userRepo, nil, nil, nil, client, nil)

	code := &service.RedeemCode{
		Code:      "233",
		Type:      service.RedeemTypeWelfare,
		Value:     10,
		Status:    service.StatusUnused,
		MaxClaims: 2,
	}
	require.NoError(t, redeemRepo.Create(ctx, code))

	first, err := redeemSvc.Redeem(ctx, 1, code.Code)
	require.NoError(t, err)
	require.Equal(t, service.RedeemTypeWelfare, first.Type)
	require.Equal(t, 1, first.ClaimedCount)
	require.Equal(t, service.StatusUnused, first.Status)

	_, err = redeemSvc.Redeem(ctx, 1, code.Code)
	require.True(t, errors.Is(err, service.ErrRedeemCodeClaimed))

	second, err := redeemSvc.Redeem(ctx, 2, code.Code)
	require.NoError(t, err)
	require.Equal(t, 2, second.ClaimedCount)
	require.Equal(t, service.StatusUsed, second.Status)

	_, err = redeemSvc.Redeem(ctx, 3, code.Code)
	require.True(t, errors.Is(err, service.ErrRedeemCodeUsed))

	updatedCode, err := redeemRepo.GetByCode(ctx, code.Code)
	require.NoError(t, err)
	require.Equal(t, 2, updatedCode.ClaimedCount)
	require.Equal(t, service.StatusUsed, updatedCode.Status)

	updatedUser1, err := userRepo.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, 10.0, updatedUser1.Balance)

	updatedUser2, err := userRepo.GetByID(ctx, 2)
	require.NoError(t, err)
	require.Equal(t, 10.0, updatedUser2.Balance)

	updatedUser3, err := userRepo.GetByID(ctx, 3)
	require.NoError(t, err)
	require.Equal(t, 0.0, updatedUser3.Balance)
}
