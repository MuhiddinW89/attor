package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MuhiddinW89/attor/pkg/config"
	"github.com/MuhiddinW89/attor/pkg/database"
)

func main() {
	cfg := config.Load()

	if cfg == nil {
		log.Fatal("Database configuration is missing")
	}

	pool, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dbName string
	var dbUser string
	var serverAddr *string

	err = pool.QueryRow(ctx, `
		SELECT
			current_database(),
			current_user,
			inet_server_addr()::text
	`).Scan(&dbName, &dbUser, &serverAddr)
	if err != nil {
		log.Fatalf("Database identity query failed: %v", err)
	}

	fmt.Println("Database connection: OK")
	fmt.Printf("Database: %s\n", dbName)
	fmt.Printf("User: %s\n", dbUser)

	if serverAddr != nil {
		fmt.Printf("Server address: %s\n", *serverAddr)
	} else {
		fmt.Println("Server address: unavailable")
	}
}
