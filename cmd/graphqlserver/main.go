package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/newsunbanjade/twitter_graphqp/config"
	"github.com/newsunbanjade/twitter_graphqp/domain"
	"github.com/newsunbanjade/twitter_graphqp/graph"
	"github.com/newsunbanjade/twitter_graphqp/postgres"
)

func main() {
	ctx := context.Background()
	config.LoadEnv(".env")
	conf := config.New()

	db := postgres.New(ctx, conf)
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}
	//REPOS

	userRepo := postgres.NewUserRepo(db)

	//Services

	authService := domain.NewAuthService(userRepo)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Timeout(time.Second * 60))

	router.Handle("/", playground.Handler("Twitter clone", "/query"))
	router.Handle("/query", handler.NewDefaultServer(graph.NewExecutableSchema(
		graph.Config{
			Resolvers: &graph.Resolver{
				AuthService: *authService,
			},
		},
	)))

	log.Fatal(http.ListenAndServe(":3000", router))
}
