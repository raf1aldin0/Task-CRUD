package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

func ConnectDB(cfg *Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database is not reachable:", err)
	}

	log.Println("Database connected")
	return db
}
