package entities

// RecipientType represents the type of recipient (To, Cc, Bcc)
type RecipientType string

const (
	RecipientTypeTo  RecipientType = "To"
	RecipientTypeCc  RecipientType = "Cc"
	RecipientTypeBcc RecipientType = "Bcc"
)

// Recipient represents an email recipient
type Recipient struct {
	RecipientID   int64         `json:"recipient_id" gorm:"primaryKey"`
	MessageID     int64         `json:"message_id" gorm:"not null"`
	Email         string        `json:"email" gorm:"not null"`
	Name          *string       `json:"name,omitempty"`
	RecipientType RecipientType `json:"recipient_type" gorm:"not null"`
}
