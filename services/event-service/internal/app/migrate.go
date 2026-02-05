package app

import (
	"log"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB, migrationsPath string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("cannot create postgres driver: %v", err)
	}

	source, err := (&file.File{}).Open(migrationsPath)
	if err != nil {
		log.Fatalf("cannot open migrations folder: %v", err)
	}

	m, err := migrate.NewWithInstance(
		"file",
		source,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("cannot create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("database migrations applied successfully")
}
