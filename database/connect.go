package database

import (
	"os"

	"github.com/kangana1024/goadmin-workshop/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	connection, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		return err
	}
	DB = connection
	err = connection.AutoMigrate(&models.User{}, &models.PasswordReset{})
	if err != nil {
		return err
	}

	return nil
}
