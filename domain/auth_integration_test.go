//go:build integration
// +build integration

package domain

import (
	"context"
	"testing"

	twittergraphql "github.com/newsunbanjade/twitter_graphqp"
	"github.com/newsunbanjade/twitter_graphqp/faker"
	testhelper "github.com/newsunbanjade/twitter_graphqp/test_helper"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAuthService_Register(t *testing.T) {
	validInput := twittergraphql.RegisterInput{
		Username:        faker.Username(),
		Email:           faker.Email(),
		Password:        "password",
		ConfirmPassword: "password",
	}

	t.Run("can register a user", func(t *testing.T) {
		ctx := context.Background()

		defer testhelper.TeardownDB(ctx, t, db)

		res, err := authService.Register(ctx, validInput)
		require.NoError(t, err)
		require.NotEmpty(t, res.User.ID)
		require.Equal(t, validInput.Email, res.User.Email)
		require.Equal(t, validInput.Username, res.User.Username)
		require.NotEqual(t, validInput.Password, res.User.Password)

	})

	t.Run("existing username", func(t *testing.T) {
		ctx := context.Background()

		defer testhelper.TeardownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)
		require.NoError(t, err)

		_, err = authService.Register(ctx, twittergraphql.RegisterInput{
			Username:        validInput.Username,
			Email:           faker.Email(),
			Password:        "server123",
			ConfirmPassword: "server123",
		})

		require.ErrorIs(t, err, twittergraphql.ErrUsernameTaken)

	})

	t.Run("existing email", func(t *testing.T) {
		ctx := context.Background()

		defer testhelper.TeardownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)
		require.NoError(t, err)

		_, err = authService.Register(ctx, twittergraphql.RegisterInput{
			Username:        faker.Username(),
			Email:           validInput.Email,
			Password:        "server123",
			ConfirmPassword: "server123",
		})

		require.ErrorIs(t, err, twittergraphql.ErrEmailTaken)

	})
}
