package main

import (
	"context"
	"log"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/config"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/database"
	"github.com/joho/godotenv"
)

func main() {
	// load .env (for local dev)
	_ = godotenv.Load()

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// connect to postgres
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer db.Close()

	log.Println("✅ PostgreSQL connected successfully")

	// optional: simple test query
	var now string
	if err := db.Pool.QueryRow(context.Background(), "SELECT now()").Scan(&now); err != nil {
		log.Fatalf("test query failed: %v", err)
	}

	log.Println("✅ Test query OK, time:", now)
}
