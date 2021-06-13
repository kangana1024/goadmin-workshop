package database

import (
	"github.com/kangana1024/goadmin-workshop/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	connection, err := gorm.Open(postgres.Open("postgresql://demo:123456@127.0.0.1/gotodo?sslmode=disable"), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		return err
	}
	DB = connection
	err = connection.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	return nil
}
