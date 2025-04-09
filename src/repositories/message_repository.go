package repositories

import (
	"context"
	"errors"
	"palm/src/entities"

	"gorm.io/gorm"
)

// Common repository errors
var (
	ErrMessageNotFound = errors.New("message not found")
)

type MessageRepository interface {
	Create(ctx context.Context, message *entities.Message) *gorm.DB
	GetByID(ctx context.Context, id uint) (*entities.Message, error)
	Update(ctx context.Context, message *entities.Message) error
	Delete(ctx context.Context, id uint) error
}
