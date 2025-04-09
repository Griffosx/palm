package repositories

import (
	"context"
	"errors"
	"palm/src/entities"
)

// Common repository errors
var (
	ErrMessageNotFound = errors.New("message not found")
)

type MessageRepository interface {
	Create(ctx context.Context, message *entities.Message) error
	GetByID(ctx context.Context, id int64) (*entities.Message, error)
	Update(ctx context.Context, message *entities.Message) error
	Delete(ctx context.Context, id int64) error
}
