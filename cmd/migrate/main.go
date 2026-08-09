package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./cmd/migrate/... [up|down]")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		log.Fatalf("create migrate instance: %v", err)
	}

	var runErr error
	switch os.Args[1] {
	case "up":
		runErr = m.Up()
	case "down":
		runErr = m.Down()
	default:
		log.Fatal("usage: go run ./cmd/migrate/... [up|down]")
	}

	if runErr != nil && !errors.Is(runErr, migrate.ErrNoChange) {
		log.Fatalf("migration failed: %v", runErr)
	}

	fmt.Println("migration command completed")
}
