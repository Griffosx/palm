package entities

// Attachment represents an email attachment
type Attachment struct {
	AttachmentID int64   `json:"attachment_id" gorm:"primaryKey"`
	MessageID    int64   `json:"message_id" gorm:"not null"`
	Filename     string  `json:"filename" gorm:"not null"`
	MimeType     string  `json:"mime_type" gorm:"not null"`
	Size         int64   `json:"size" gorm:"not null"`
	LocalPath    *string `json:"local_path,omitempty"`
	Content      []byte  `json:"content,omitempty" gorm:"-"` // Exclude from database
}
