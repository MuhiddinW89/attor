package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/MuhiddinW89/attor/pkg/config"
	"github.com/MuhiddinW89/attor/pkg/database"
)

func main() {
	cfg := config.Load()
	pool, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal("database unavailable")
	}
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	service := identity.NewService(identity.NewPostgresRepository(pool))
	_, err = service.CreateInitialOwner(
		ctx,
		os.Getenv("INITIAL_OWNER_EMAIL"),
		os.Getenv("INITIAL_OWNER_PASSWORD"),
		os.Getenv("INITIAL_STORE_NAME"),
	)
	if err != nil {
		log.Fatal("initial owner was not created")
	}
	log.Print("Initial owner created")
}
