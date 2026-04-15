package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		url.QueryEscape(os.Getenv("DB_PASSWORD")),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	dir, err := filepath.Abs(os.Getenv("MIGRATIONS_DIR"))
	if err != nil {
		log.Fatalf("get absolute path to migrations: %v", err)
	}

	if err := goose.Up(db, dir); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Printf("Successfully applied the migrations")
}
