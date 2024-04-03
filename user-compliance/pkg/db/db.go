package db

import (
	"log"

	"github.com/Sharefunds/user-compliance/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{}, &models.BankAccount{}, &models.PlaidItem{})

	return Handler{db}
}
