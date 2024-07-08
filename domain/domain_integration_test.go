//go:build integration
// +build integration

package domain

import (
	"context"
	"log"
	"os"
	"testing"

	twittergraphql "github.com/newsunbanjade/twitter_graphqp"
	"github.com/newsunbanjade/twitter_graphqp/config"
	"github.com/newsunbanjade/twitter_graphqp/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	conf        *config.Config
	db          *postgres.DB
	authService *AuthService
	userRepo    twittergraphql.UserRepo
)

func TestMain(m *testing.M) {

	ctx := context.Background()
	config.LoadEnv(".env.test")
	PasswordCost = bcrypt.MinCost
	conf = config.New()
	db = postgres.New(ctx, conf)
	defer db.Close()
	if err := db.Drop(); err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}
	userRepo = postgres.NewUserRepo(db)
	authService = NewAuthService(userRepo)
	os.Exit(m.Run())
}
