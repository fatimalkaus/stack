package db

import (
	"fmt"

	"github.com/fatimalkaus/stack/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitPostgres initializes *gorm.DB using config.
func InitPostgres(cfg config.Postgres) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
