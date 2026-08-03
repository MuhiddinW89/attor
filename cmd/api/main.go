package main

import (
	"log"

	"github.com/MuhiddinW89/attor/internal/bootstrap"
	"github.com/MuhiddinW89/attor/pkg/config"
)

func main() {

	cfg := config.Load()

	app, err := bootstrap.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on :8080")

	log.Fatal(
		app.Listen(":8080"),
	)
}
