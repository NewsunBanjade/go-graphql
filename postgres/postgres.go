package postgres

import (
	"context"
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/newsunbanjade/twitter_graphqp/config"
)

type DB struct {
	Pool *pgxpool.Pool
	conf *config.Config
}

func New(ctx context.Context, config *config.Config) *DB {
	dbConf, err := pgxpool.ParseConfig(config.Database.URL)
	if err != nil {
		log.Fatalf("can't parse postgres config: %v", err)
	}
	pool, err := pgxpool.ConnectConfig(ctx, dbConf)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	db := &DB{Pool: pool, conf: config}

	db.Ping(ctx)
	return db

}

func (db *DB) Ping(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("can't ping postgres: %v", err)
	}
	log.Println("postgres connected")
}

func (db *DB) Open(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("can't ping postgres: %v", err)
	}
	log.Println("database connected")
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) Migrate() error {

	_, b, _, _ := runtime.Caller(0)
	migrationPath := fmt.Sprintf("file:///%s/migrations", path.Dir(b))

	m, err := migrate.New(migrationPath, db.conf.Database.URL)
	if err != nil {
		return fmt.Errorf("error create the migration instance : %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error migrate up : %v", err)
	}
	log.Println(("migration done"))
	return nil
}
func (db *DB) Drop() error {

	_, b, _, _ := runtime.Caller(0)
	migrationPath := fmt.Sprintf("file:///%s/migrations", path.Dir(b))

	m, err := migrate.New(migrationPath, db.conf.Database.URL)
	if err != nil {
		return fmt.Errorf("error drop the database instance : %v", err)
	}

	if err := m.Drop(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error drop  : %v", err)
	}
	log.Println(("dropping done"))
	return nil
}

func (db *DB) Truncate(ctx context.Context) error {
	if _, err := db.Pool.Exec(ctx, `DELETE FROM users;`); err != nil {
		return fmt.Errorf("error truncate: %v", err)
	}
	return nil
}
