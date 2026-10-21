package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/MuhiddinW89/attor/pkg/config"
	"github.com/MuhiddinW89/attor/pkg/database"
)

func main() {
	email := strings.TrimSpace(os.Getenv("INITIAL_OWNER_EMAIL"))
	password := os.Getenv("INITIAL_OWNER_PASSWORD")
	storeName := strings.TrimSpace(os.Getenv("INITIAL_STORE_NAME"))

	if email == "" || password == "" || storeName == "" {
		log.Fatal("Missing required variables: INITIAL_OWNER_EMAIL, INITIAL_OWNER_PASSWORD, INITIAL_STORE_NAME")
	}

	cfg := config.Load()

	if cfg == nil {
		log.Fatal("Database configuration is missing")
	}

	pool, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	service := identity.NewService(identity.NewPostgresRepository(pool))

	_, err = service.CreateInitialOwner(ctx, email, password, storeName)
	if err != nil {
		log.Fatalf("Initial owner was not created: %v", err)
	}

	log.Print("Initial owner created successfully")
}
