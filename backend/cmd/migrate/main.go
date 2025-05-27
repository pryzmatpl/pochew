package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/readitlater?sslmode=disable"
	}

	// Open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create postgres driver:", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("Migrations completed successfully!")
} 