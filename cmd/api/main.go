
package main

import (
	"log"
	"os"

	"github.com/MuhiddinW89/attor/internal/bootstrap"
	"github.com/MuhiddinW89/attor/pkg/config"
)

func main() {
	cfg := config.Load()

	app, err := bootstrap.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
