package db

import (
	"log"

	"github.com/Sharefunds/pool-transactions/pkg/models"

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

	db.AutoMigrate(&models.Pool{}, &models.Participant{}, &models.Transaction{}, &models.PoolInvitation{})

	return Handler{db}
}

func (h *Handler) IsMemberOfPool(userId string, poolId string) (bool, error) {
	var participant models.Participant

	result := h.DB.Where("user_id = ? AND pool_id = ?", userId, poolId).First(&participant)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
