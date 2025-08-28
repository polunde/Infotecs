package app

import (
	"fmt"
	"infotecs/internal/db"
	"infotecs/internal/entity"
	"infotecs/internal/handlers"
	"infotecs/internal/httpserver"
	"infotecs/internal/repository"
	"infotecs/internal/service"
	"infotecs/internal/utils"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

func Run() error {
	if err := godotenv.Load("config/.env"); err != nil {
		return fmt.Errorf("could not setup config: %w", err)
	}

	dbCfg := &db.DBConfig{
		Host:   os.Getenv("DB_ADDRESS"),
		Port:   os.Getenv("DB_PORT"),
		User:   os.Getenv("DB_USERNAME"),
		Pass:   os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
	}

	gormDB, err := db.InitDB(dbCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	var count int64
	if err := gormDB.Model(&entity.Wallet{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed count wallets: %w", err)
	}
	if count == 0 {
		wallets := utils.GenerateWallets(10, 100)
		if err := gormDB.Create(&wallets).Error; err != nil {
			return fmt.Errorf("failed to create initial wallets: %w", err)
		}
		log.Printf("Created 10 initial wallets with UUID addresses")
	}

	walletRepo := repository.NewWalletRepository(gormDB)
	walletService := service.NewWalletService(walletRepo)

	router := handlers.NewRouter(walletService)

	host := os.Getenv("API_INTERFACE")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	srv := httpserver.NewServer(router, &httpserver.ServerConfig{
		Host: host,
		Port: port,
	})

	srv.Start()
	log.Printf("Server started on port %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case err := <-srv.Notify():
		return fmt.Errorf("server stopped unexpectedly: %w", err)
	case sig := <-quit:
		log.Printf("Interrupt signal received: %v", sig)
	}

	return nil
}
