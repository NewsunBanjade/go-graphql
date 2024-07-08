package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type database struct {
	URL string
}

type Config struct {
	Database database
}

func LoadEnv(fileName string) {
	re := regexp.MustCompile(`^(.*` + "twitter_graphql" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	err := godotenv.Load(string(rootPath) + "/" + fileName)
	if err != nil {
		godotenv.Load()
	}
}

func New() *Config {
	godotenv.Load()
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), "5432", os.Getenv("DB_NAME"))
	return &Config{
		Database: database{
			URL: databaseURL,
		},
	}
}
