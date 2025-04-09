package entities

import (
	"time"

	"gorm.io/gorm"
)

type Importance string

const (
	ImportanceLow    Importance = "Low"
	ImportanceNormal Importance = "Normal"
	ImportanceHigh   Importance = "High"
)

// Message represents an email message
type Message struct {
	gorm.Model
	Subject          *string      `json:"subject,omitempty"`
	Body             *string      `json:"body,omitempty"`
	BodyPreview      *string      `json:"body_preview,omitempty"`
	SenderEmail      string       `json:"sender_email" gorm:"not null"`
	SenderName       *string      `json:"sender_name,omitempty"`
	ReceivedDatetime *time.Time   `json:"received_datetime,omitempty"`
	SentDatetime     *time.Time   `json:"sent_datetime,omitempty"`
	IsDraft          bool         `json:"is_draft" gorm:"not null"`
	IsRead           bool         `json:"is_read" gorm:"not null"`
	Importance       Importance   `json:"importance" gorm:"not null"`
	ConversationID   *string      `json:"conversation_id,omitempty"`
	AccountID        uint         `json:"account_id"`
	Account          Account      `json:"account,omitempty"`
	Attachments      []Attachment `json:"attachments,omitempty"`
	Recipients       []Recipient  `json:"recipients,omitempty"`
}
