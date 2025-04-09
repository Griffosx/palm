package repositories

import (
	"context"
	"errors"
	"palm/src/entities"
)

// Common repository errors
var (
	ErrRecipientNotFound = errors.New("recipient not found")
)

type RecipientRepository interface {
	Create(ctx context.Context, recipient *entities.Recipient) error
	GetByMessageID(ctx context.Context, messageID uint) ([]*entities.Recipient, error)
	DeleteByMessageID(ctx context.Context, messageID uint) error
}
