package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func DefaultConfig() Config {
	host := os.Getenv("DB_HOST")
	if host == "" {
		log.Fatal("DB_HOST environment variable is required")
	}

	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		log.Fatal("DB_PORT environment variable is required")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Invalid DB_PORT:", err)
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		log.Fatal("DB_USER environment variable is required")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatal("DB_PASSWORD environment variable is required")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is required")
	}

	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		log.Fatal("DB_SSLMODE environment variable is required")
	}

	return Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}
}
