package entities

import "gorm.io/gorm"

// RecipientType represents the type of recipient (To, Cc, Bcc)
type RecipientType string

const (
	RecipientTypeTo  RecipientType = "To"
	RecipientTypeCc  RecipientType = "Cc"
	RecipientTypeBcc RecipientType = "Bcc"
)

// Recipient represents an email recipient
type Recipient struct {
	gorm.Model
	Email         string        `json:"email" gorm:"not null"`
	Name          *string       `json:"name,omitempty"`
	RecipientType RecipientType `json:"recipient_type" gorm:"not null"`
	MessageID     uint          `json:"message_id"`
	Message       Message       `json:"message,omitempty"`
}
