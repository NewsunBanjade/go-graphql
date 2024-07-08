package postgres

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	twittergraphql "github.com/newsunbanjade/twitter_graphqp"
)

type UserRepo struct {
	DB *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (ur *UserRepo) Create(ctx context.Context, user twittergraphql.User) (twittergraphql.User, error) {
	tx, err := ur.DB.Pool.Begin(ctx)
	if err != nil {
		return twittergraphql.User{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	user, err = createUser(ctx, tx, user)
	if err != nil {
		return twittergraphql.User{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return twittergraphql.User{}, fmt.Errorf("error commiting: %v", err)
	}
	return user, nil

}

func createUser(ctx context.Context, tx pgx.Tx, user twittergraphql.User) (twittergraphql.User, error) {
	query := `INSERT INTO users (email,username,password) VALUES ($1,$2,$3) RETURNING *;`
	u := twittergraphql.User{}
	if err := pgxscan.Get(ctx, tx, &u, query, user.Email, user.Username, user.Password); err != nil {
		return twittergraphql.User{}, fmt.Errorf("error insert: %v", err)
	}

	return u, nil

}

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (twittergraphql.User, error) {

	query := `SELECT * FROM users WHERE username=$1`
	u := twittergraphql.User{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, username); err != nil {
		if pgxscan.NotFound(err) {

			return twittergraphql.User{}, twittergraphql.ErrNotFound
		}
		return twittergraphql.User{}, fmt.Errorf("error select : %v", err)
	}

	return u, nil
}
func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (twittergraphql.User, error) {

	query := `SELECT * FROM users WHERE email=$1`
	u := twittergraphql.User{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, email); err != nil {
		if pgxscan.NotFound(err) {

			return twittergraphql.User{}, twittergraphql.ErrNotFound
		}
		return twittergraphql.User{}, fmt.Errorf("error select : %v", err)
	}

	return u, nil
}
