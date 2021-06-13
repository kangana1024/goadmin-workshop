package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"type:VARCHAR(100)"`
	LastName  string `gorm:"type:VARCHAR(100)"`
	Email     string `gorm:"unique"`
	Password  []byte `gorm:"type:VARCHAR(100)"`
}
