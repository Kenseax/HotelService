package main

import (
	"HotelService/api/rest/controller"
	"HotelService/infrastructure/db"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	dbConfig := db.DefaultConfig()

	database, err := db.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	if err := runMigrations(database); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	mux := controller.SetupRoutes(database)

	port := os.Getenv("PORT")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func runMigrations(db *sql.DB) error {
	migrationSQL, err := os.ReadFile("infrastructure/db/migrations/001_create_tables.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}
