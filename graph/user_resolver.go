package graph

import (
	"context"

	twittergraphql "github.com/newsunbanjade/twitter_graphqp"
)

func mapUser(u twittergraphql.User) *User {
	return &User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func (q *queryResolver) Me(ctx context.Context) (*User, error) {
	panic("implement")
}
