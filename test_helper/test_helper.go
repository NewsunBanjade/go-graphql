package testhelper

import (
	"context"
	"testing"

	"github.com/newsunbanjade/twitter_graphqp/postgres"
	"github.com/stretchr/testify/require"
)

func TeardownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()

	err := db.Truncate(ctx)
	require.NoError(t, err)
}
