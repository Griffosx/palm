package repositories

import (
	"context"
	"errors"
	"palm/src/entities"
)

// Common repository errors
var (
	ErrAttachmentNotFound = errors.New("attachment not found")
)

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *entities.Attachment) error
	GetByID(ctx context.Context, id int64) (*entities.Attachment, error)
	GetByMessageID(ctx context.Context, messageID int64) ([]*entities.Attachment, error)
	DeleteByMessageID(ctx context.Context, messageID int64) error
}
