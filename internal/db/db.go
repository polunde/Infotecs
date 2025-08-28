package db

import (
	"fmt"
	"infotecs/internal/entity"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
}

var db *gorm.DB

func InitDB(cfg *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DBName)
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&entity.Wallet{}, &entity.Transaction{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	return db, nil
}
