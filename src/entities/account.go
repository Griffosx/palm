package entities

import "gorm.io/gorm"

const (
	AccountTypeMicrosoft = "Microsoft"
	AccountTypeGoogle    = "Google"
)

type Account struct {
	gorm.Model
	Email       string    `json:"email" gorm:"unique;not null"`
	AccountType string    `json:"account_type" gorm:"not null"`
	Messages    []Message `json:"messages,omitempty"`
}
