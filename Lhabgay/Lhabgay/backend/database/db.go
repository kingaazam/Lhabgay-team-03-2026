package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		connStr = "postgres://postgres:12345@host.docker.internal:5432/bookdb?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database ping error: %v", err)
	}

	log.Println("Successfully connected to the database!")

	// FORCE-FIX THE DATABASE COLUMNS LIVE ON LAUNCH
	_, err = DB.Exec(`
		ALTER TABLE books ADD COLUMN IF NOT EXISTS cover_image TEXT;
		ALTER TABLE books ADD COLUMN IF NOT EXISTS book_file_path TEXT;
	`)
	if err != nil {
		log.Println("Database setup auto-check warning:", err)
	}
}
