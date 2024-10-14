package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        uint   `gorm:"primarykey"`
	FirstName string `gorm:"size:100"`
	LastName  string `gorm:"size:100"`
	Company   string `gorm:"size:100"`
	Address   string `gorm:"size:255"`
	City      string `gorm:"size:100"`
	County    string `gorm:"size:100"`
	Postal    string `gorm:"size:20"`
	Phone     string `gorm:"size:20"`
	Email     string `gorm:"size:100"`
	Web       string `gorm:"size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Customer{})
}
