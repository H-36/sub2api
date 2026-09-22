//go:build integration

package repository

import (
	"context"
	"os"
	"testing"

	dbmigrations "github.com/Wei-Shaw/sub2api/migrations"
	"github.com/stretchr/testify/require"
)

func TestForkDeepSeekMigration(t *testing.T) {
	tx := testTx(t)
	ctx := context.Background()
	fixture, err := os.ReadFile("../../migrations/testdata/fork_deepseek_setup.sql")
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, string(fixture))
	require.NoError(t, err)
	migration, err := dbmigrations.FS.ReadFile("240_migrate_fork_deepseek_accounts.sql")
	require.NoError(t, err)
	checks, err := os.ReadFile("../../migrations/testdata/fork_deepseek_checks.sql")
	require.NoError(t, err)
	for range 2 {
		_, err = tx.ExecContext(ctx, string(migration))
		require.NoError(t, err)
		_, err = tx.ExecContext(ctx, string(checks))
		require.NoError(t, err)
	}
}
