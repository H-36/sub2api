package service_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/repository"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newWelfareRedeemServiceTestClient(t *testing.T) (*sql.DB, *dbent.Client) {
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

	return db, client
}

func createWelfareRedeemUser(t *testing.T, ctx context.Context, client *dbent.Client, email string) *service.User {
	t.Helper()

	created, err := client.User.Create().
		SetEmail(email).
		SetPasswordHash("test-password-hash").
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		Save(ctx)
	require.NoError(t, err)

	return &service.User{
		ID:          created.ID,
		Email:       created.Email,
		Role:        created.Role,
		Status:      created.Status,
		Balance:     created.Balance,
		Concurrency: created.Concurrency,
	}
}

func TestRedeemService_WelfareRedeemFlow(t *testing.T) {
	ctx := context.Background()
	db, client := newWelfareRedeemServiceTestClient(t)

	redeemRepo := repository.NewRedeemCodeRepository(client)
	userRepo := repository.NewUserRepository(client, db)
	redeemSvc := service.NewRedeemService(redeemRepo, userRepo, nil, nil, nil, client, nil)

	user1 := createWelfareRedeemUser(t, ctx, client, "welfare-1@test.com")
	user2 := createWelfareRedeemUser(t, ctx, client, "welfare-2@test.com")
	user3 := createWelfareRedeemUser(t, ctx, client, "welfare-3@test.com")

	code := &service.RedeemCode{
		Code:      "233",
		Type:      service.RedeemTypeWelfare,
		Value:     10,
		Status:    service.StatusUnused,
		MaxClaims: 2,
	}
	require.NoError(t, redeemRepo.Create(ctx, code))

	first, err := redeemSvc.Redeem(ctx, user1.ID, code.Code)
	require.NoError(t, err)
	require.Equal(t, service.RedeemTypeWelfare, first.Type)
	require.Equal(t, 1, first.ClaimedCount)
	require.Equal(t, service.StatusUnused, first.Status)

	_, err = redeemSvc.Redeem(ctx, user1.ID, code.Code)
	require.True(t, errors.Is(err, service.ErrRedeemCodeClaimed))

	second, err := redeemSvc.Redeem(ctx, user2.ID, code.Code)
	require.NoError(t, err)
	require.Equal(t, 2, second.ClaimedCount)
	require.Equal(t, service.StatusUsed, second.Status)

	_, err = redeemSvc.Redeem(ctx, user3.ID, code.Code)
	require.True(t, errors.Is(err, service.ErrRedeemCodeUsed))

	updatedCode, err := redeemRepo.GetByCode(ctx, code.Code)
	require.NoError(t, err)
	require.Equal(t, 2, updatedCode.ClaimedCount)
	require.Equal(t, service.StatusUsed, updatedCode.Status)

	updatedUser1, err := userRepo.GetByID(ctx, user1.ID)
	require.NoError(t, err)
	require.Equal(t, 10.0, updatedUser1.Balance)

	updatedUser2, err := userRepo.GetByID(ctx, user2.ID)
	require.NoError(t, err)
	require.Equal(t, 10.0, updatedUser2.Balance)

	updatedUser3, err := userRepo.GetByID(ctx, user3.ID)
	require.NoError(t, err)
	require.Equal(t, 0.0, updatedUser3.Balance)
}
