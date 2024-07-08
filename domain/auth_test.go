package domain

import (
	"context"
	"errors"
	"testing"

	twittergraphql "github.com/newsunbanjade/twitter_graphqp"
	"github.com/newsunbanjade/twitter_graphqp/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	validInput := twittergraphql.RegisterInput{
		Username:        "Newsun",
		Email:           "Mail@mail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}
	t.Run("can Register", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twittergraphql.User{}, twittergraphql.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twittergraphql.User{}, twittergraphql.ErrNotFound)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(twittergraphql.User{
			ID:       "123",
			Email:    validInput.Email,
			Username: validInput.Username,
		}, nil)
		service := NewAuthService(userRepo)
		res, err := service.Register(ctx, validInput)
		require.NoError(t, err)
		require.NotEmpty(t, res.User.ID)
		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.Email)
		require.NotEmpty(t, res.User.Username)

		userRepo.AssertExpectations(t)

	})
	t.Run("username Taken", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twittergraphql.User{}, nil)

		service := NewAuthService(userRepo)
		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, twittergraphql.ErrUsernameTaken)

		userRepo.AssertNotCalled(t, "Create")

		userRepo.AssertExpectations(t)

	})

	t.Run("Email Taken", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twittergraphql.User{}, twittergraphql.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twittergraphql.User{}, nil)

		service := NewAuthService(userRepo)
		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, twittergraphql.ErrEmailTaken)

		userRepo.AssertNotCalled(t, "Create")

		userRepo.AssertExpectations(t)

	})

	t.Run("create error", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twittergraphql.User{}, twittergraphql.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twittergraphql.User{}, twittergraphql.ErrNotFound)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(twittergraphql.User{
			ID:       "123",
			Email:    validInput.Email,
			Username: validInput.Username,
		}, errors.New("something"))
		service := NewAuthService(userRepo)
		_, err := service.Register(ctx, validInput)
		require.Error(t, err)

		userRepo.AssertExpectations(t)

	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		service := NewAuthService(userRepo)
		_, err := service.Register(ctx, twittergraphql.RegisterInput{})
		require.ErrorIs(t, err, twittergraphql.ErrValidation)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertNotCalled(t, "Create")

		userRepo.AssertExpectations(t)

	})
}

func TestAuthServiceLogin(t *testing.T) {
	validInput := twittergraphql.LoginInput{
		Email:    "mymail@uk.co",
		Password: "Server@ASP",
	}
	t.Run("can login", func(t *testing.T) {
		ctx := context.Background()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(validInput.Password), PasswordCost)
		require.NoError(t, err)
		userRepo := &mocks.UserRepo{}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twittergraphql.User{
			Email:    validInput.Email,
			Password: string(hashedPassword),
		}, nil)
		service := NewAuthService(userRepo)
		_, err = service.Login(ctx, validInput)
		require.NoError(t, err)
		userRepo.AssertExpectations(t)

	})

	t.Run("wrong password", func(t *testing.T) {
		ctx := context.Background()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(validInput.Password), PasswordCost)
		require.NoError(t, err)
		userRepo := &mocks.UserRepo{}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twittergraphql.User{
			Email:    validInput.Email,
			Password: string(hashedPassword),
		}, nil)
		service := NewAuthService(userRepo)
		validInput.Password = "somethingelse"
		_, err = service.Login(ctx, validInput)
		require.ErrorIs(t, err, twittergraphql.ErrBadCredentials)
		userRepo.AssertExpectations(t)

	})

	t.Run("email not found", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twittergraphql.User{}, twittergraphql.ErrNotFound)
		service := NewAuthService(userRepo)
		validInput.Password = "somethingelse"
		_, err := service.Login(ctx, validInput)
		require.ErrorIs(t, err, twittergraphql.ErrBadCredentials)
		userRepo.AssertExpectations(t)

	})

	t.Run("get user by email error", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twittergraphql.User{}, errors.New("some errors"))
		service := NewAuthService(userRepo)
		_, err := service.Login(ctx, validInput)
		require.Error(t, err)
		userRepo.AssertExpectations(t)

	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		service := NewAuthService(userRepo)
		_, err := service.Login(ctx, twittergraphql.LoginInput{
			Email:    "bob",
			Password: "",
		})
		require.ErrorIs(t, err, twittergraphql.ErrValidation)
		userRepo.AssertExpectations(t)

	})

}
