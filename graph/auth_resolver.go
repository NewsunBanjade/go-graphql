package graph

import (
	"context"
	"errors"

	twittergraphql "github.com/newsunbanjade/twitter_graphqp"
)

func mapAuthResponse(a twittergraphql.AuthResponse) *AuthResponse {
	return &AuthResponse{
		AccessToken: a.AccessToken,
		User:        mapUser(a.User),
	}
}

func (m *mutationResolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	res, err := m.AuthService.Login(ctx, twittergraphql.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, twittergraphql.ErrValidation) || errors.Is(err, twittergraphql.ErrEmailTaken) || errors.Is(err, twittergraphql.ErrUsernameTaken):
			return nil, buildBadRequestError(ctx, err)

		default:
			return nil, err
		}
	}
	return mapAuthResponse(res), nil

}
func (m *mutationResolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	res, err := m.AuthService.Register(ctx, twittergraphql.RegisterInput{
		Email:           input.Email,
		Password:        input.Password,
		Username:        input.Username,
		ConfirmPassword: input.ConfirmPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, twittergraphql.ErrValidation) || errors.Is(err, twittergraphql.ErrBadCredentials):
			return nil, buildBadRequestError(ctx, err)

		default:
			return nil, err
		}
	}
	return mapAuthResponse(res), nil
}
