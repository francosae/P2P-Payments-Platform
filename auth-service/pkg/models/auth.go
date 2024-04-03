package models

import (
	"time"
)

type User struct {
	UID       string    `gorm:"primaryKey;column:uid"`
	Email     string    `gorm:"unique;column:email;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
