package database

import (
	"fmt"
	"library-system/conf"
	"library-system/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDBConnection(config *conf.Configuration) (*gorm.DB, error) {
	dbURL := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName,
	)
	dial := postgres.Open(dbURL)

	database, err := gorm.Open(dial, &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return database, err
	}
	database.AutoMigrate(&models.Book{}, &models.Borrower{})
	return database, nil
}
