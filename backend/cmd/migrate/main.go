package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// load .env for local/dev
	_ = godotenv.Load()

	// command: up | down | force
	cmd := flag.String("cmd", "up", "migration command: up | down | force")
	version := flag.Int("version", 0, "version to force (only used with -cmd force)")
	flag.Parse()

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// connect DB using stdlib sql
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
		cfg.Database.User, cfg.Database.Password, cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping error: %v", err)
	}

	// create migrate driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("migrate driver error: %v", err)
	}

	// init migrate
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("migrate init error: %v", err)
	}

	// execute command
	switch *cmd {
	case "up":
		log.Println("⬆ running migrations UP")
		err = m.Up()

	case "down":
		log.Println("⬇ running migrations DOWN")
		err = m.Down()

	case "force":
		log.Printf("🔧 forcing version %d", *version)
		err = m.Force(*version)

	default:
		log.Fatalf("unknown command: %s (use up, down, or force)", *cmd)
	}

	// handle result
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("✅ migration finished successfully")
}
