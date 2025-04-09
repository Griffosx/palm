package entities

import "gorm.io/gorm"

// Attachment represents an email attachment
type Attachment struct {
	gorm.Model
	Filename  string  `json:"filename" gorm:"not null"`
	MimeType  string  `json:"mime_type" gorm:"not null"`
	Size      uint    `json:"size" gorm:"not null"`
	LocalPath *string `json:"local_path,omitempty"`
	MessageID uint    `json:"message_id"`
	Message   Message `json:"message,omitempty"`
}
