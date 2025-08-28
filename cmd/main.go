package main

import (
	"infotecs/internal/app"
	"log"
	"os"
)

func main() {
	if err := app.Run(); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}
}
