package repository

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// NewPostgresDB creates a sql.DB using DATABASE_URL
func NewPostgresDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
