package db

import (
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	DB *sqlx.DB
}

func ConnectPostgres(url string) (*Storage, error) {
	data, err := sqlx.Open("postgres", url)
	if err != nil {
		return &Storage{}, err
	}

	err = migratePostgres(data)
	if err != nil {
		return &Storage{}, err
	}
	return &Storage{data}, nil
}

func migratePostgres(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/storage/db/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
