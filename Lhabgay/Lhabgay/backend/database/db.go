package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// DB is the shared PostgreSQL connection used by controllers and models.
var DB *sql.DB

// ConnectDB opens and verifies the PostgreSQL database connection.
func ConnectDB() {
	connStr := "host=host.docker.internal user=postgres password=postgres dbname=lhabgay_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping error: ", err)
	}

	DB = db
	log.Println("Database connected successfully")
}
