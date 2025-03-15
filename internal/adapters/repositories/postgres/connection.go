package postgres

import (
	"fmt"
	"log"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/config"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Restaurant{},
		&domain.Volunteer{},
		&domain.Event{},
		&domain.VolunteerApplication{},
		&domain.EventVolunteer{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connected and migrations completed successfully")
	return db, nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
