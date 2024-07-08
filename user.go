package twittergraphql

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUsernameTaken = errors.New("username Taken")
	ErrEmailTaken    = errors.New("email Taken")
)

type UserRepo interface {
	Create(ctx context.Context, username User) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByEmail(ctx context.Context, username string) (User, error)
}

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
