package database

import (
	"database/sql"
	"log"
	"os" // <-- 1. Added this to read environment variables

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	// 2. Check if Render gave us a cloud database URL
	connStr := os.Getenv("DATABASE_URL")

	// 3. If there is no cloud database (meaning you are running it on your own PC),
	//    fallback to your local login string.
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
}
