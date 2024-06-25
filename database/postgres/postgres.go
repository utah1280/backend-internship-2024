package postgres

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresConnection() *sqlx.DB {
	connStr := "user=root dbname=korzinka sslmode=disable password=root host=localhost"
	DB, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		DB.Close()
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	return DB
}
