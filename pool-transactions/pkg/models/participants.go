package models

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ParticipantID uuid.UUID `gorm:"type:uuid;primaryKey"`
	PoolID        uuid.UUID `gorm:"type:uuid;not null"`
	UserID        string    `gorm:"size:50;not null"`
	Role          string    `gorm:"size:50"`
	JoinedAt      time.Time `gorm:"autoCreateTime"`
}
