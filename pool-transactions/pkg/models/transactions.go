package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	TransactionID uuid.UUID `gorm:"type:uuid;primaryKey"`
	FromAccountID string    `gorm:"size:20;not null"`
	ToAccountID   string    `gorm:"size:20;not null"`
	Amount        float64   `gorm:"not null"`
	Description   string    `gorm:"size:255"`
	Status        string    `gorm:"size:50"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}
