package main

import (
	"log"

	"LinkUp/internal/app"

	"github.com/joho/godotenv"
)
// Auto-generated swagger comments for main
// @Summary Auto-generated summary for main
// @Description Auto-generated description for main — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func main() {
	_ = godotenv.Load()
	if err := app.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
