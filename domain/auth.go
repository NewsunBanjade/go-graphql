package domain

import (
	"context"
	"errors"
	"fmt"

	twittergraphql "github.com/newsunbanjade/twitter_graphqp"
	"golang.org/x/crypto/bcrypt"
)

var (
	PasswordCost int = bcrypt.DefaultCost
)

type AuthService struct {
	UserRepo twittergraphql.UserRepo
}

func NewAuthService(ur twittergraphql.UserRepo) *AuthService {
	return &AuthService{
		UserRepo: ur,
	}
}

func (as *AuthService) Register(ctx context.Context, input twittergraphql.RegisterInput) (twittergraphql.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twittergraphql.AuthResponse{}, err
	}
	// check if username is already taken
	if _, err := as.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, twittergraphql.ErrNotFound) {
		return twittergraphql.AuthResponse{}, twittergraphql.ErrUsernameTaken
	}
	//check if email is already taken
	if _, err := as.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, twittergraphql.ErrNotFound) {
		return twittergraphql.AuthResponse{}, twittergraphql.ErrEmailTaken
	}

	user := twittergraphql.User{
		Email:    input.Email,
		Username: input.Username,
	}

	// hash the password
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), PasswordCost)
	if err != nil {
		return twittergraphql.AuthResponse{}, fmt.Errorf("error hashing password: %v", err)
	}
	user.Password = string(hashpassword)
	// create the user
	createdUser, err := as.UserRepo.Create(ctx, user)
	if err != nil {
		return twittergraphql.AuthResponse{}, fmt.Errorf("error Creating User: %v", err)
	}

	// return accessToken and User
	return twittergraphql.AuthResponse{AccessToken: "sample token", User: createdUser}, nil

}

func (as *AuthService) Login(ctx context.Context, input twittergraphql.LoginInput) (twittergraphql.AuthResponse, error) {
	input.Sanitize()
	if err := input.Validate(); err != nil {
		return twittergraphql.AuthResponse{}, err
	}
	user, err := as.UserRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, twittergraphql.ErrNotFound):
			return twittergraphql.AuthResponse{}, twittergraphql.ErrBadCredentials
		default:
			return twittergraphql.AuthResponse{}, err
		}

	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return twittergraphql.AuthResponse{}, twittergraphql.ErrBadCredentials
	}
	return twittergraphql.AuthResponse{AccessToken: "jwt", User: user}, nil
}
