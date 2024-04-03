package models

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	Address1    string `gorm:"size:40;not null"`
	Address2    string `gorm:"size:40"`
	Address3    string `gorm:"size:30"`
	Address4    string `gorm:"size:30"`
	Address5    string `gorm:"size:30"`
	City        string `gorm:"size:30;not null"`
	State       string `gorm:"size:3;not null"`
	PostalCode  string `gorm:"size:10;not null"`
	CountryCode string `gorm:"size:3;not null"`
}

type PersonalInfo struct {
	FirstName   string    `gorm:"size:50;not null"`
	LastName    string    `gorm:"size:50;not null"`
	Email       string    `gorm:"size:50;not null"`
	PhoneNumber string    `gorm:"size:20;not null"`
	DateOfBirth time.Time `gorm:"type:date;not null"`
}

type User struct {
	UID          string         `gorm:"primaryKey"`
	Username     string         `gorm:"unique;not null"`
	PRN          string         `gorm:"size:20"`
	PersonalInfo PersonalInfo   `gorm:"embedded;embeddedPrefix:personal_info_"`
	Address      Address        `gorm:"embedded;embeddedPrefix:addr_"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	BankAccounts []BankAccount `gorm:"foreignKey:UserID"`
	PlaidItems   []PlaidItem   `gorm:"foreignKey:UserID"`

	IsVerified    bool   `gorm:"default:false"`
	IdentityToken string `gorm:"type:text"`
}

type BankAccount struct {
	gorm.Model
	UserID          string `gorm:"index;not null"`
	InstitutionName string
	AccountName     string
	AccountType     string
	Mask            string
	User            User        `gorm:"foreignKey:UserID"`
	PlaidItems      []PlaidItem `gorm:"foreignKey:BankAccountID"`
}

type PlaidItem struct {
	gorm.Model
	UserID        string
	BankAccountID uint   `gorm:"index;not null"`
	AccessToken   string `gorm:"type:text"`
	InstitutionID string
	User          User        `gorm:"foreignKey:UserID"`
	BankAccount   BankAccount `gorm:"foreignKey:BankAccountID"`
}

type SensitiveData struct {
	IDType int32  `gorm:"size:10;not null"`
	ID     string `gorm:"size:20;not null"`
}
