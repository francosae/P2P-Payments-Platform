package models

import (
	"time"

	"github.com/google/uuid"
)

type Pool struct {
	PoolID       uuid.UUID     `gorm:"type:uuid;primaryKey"`
	PRN          string        `gorm:"size:20;not null"`
	PoolName     string        `gorm:"size:100;not null"`
	Description  string        `gorm:"size:255"`
	UserID       string        `gorm:"size:50;not null"`
	BalanceGoal  float64       `gorm:"not null"`
	Participants []Participant `gorm:"foreignKey:PoolID"`
	Status       string        `gorm:"size:20;not null"`
	CreatedAt    time.Time     `gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime"`
}

type PoolInvitation struct {
	InvitationID uuid.UUID `gorm:"type:uuid;primaryKey"`
	PoolID       string    `gorm:"not null"`
	InviterID    string    `gorm:"not null"`
	InviteeID    string    `gorm:"not null"`
	Status       string    `gorm:"size:20"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
