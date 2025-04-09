package repositories

import (
	"context"
	"errors"
	"palm/src/entities"

	"gorm.io/gorm"
)

// Common repository errors
var (
	ErrAttachmentNotFound = errors.New("attachment not found")
)

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *entities.Attachment) *gorm.DB
	GetByID(ctx context.Context, id uint) (*entities.Attachment, error)
	GetByMessageID(ctx context.Context, messageID uint) ([]*entities.Attachment, error)
	DeleteByMessageID(ctx context.Context, messageID uint) error
}
