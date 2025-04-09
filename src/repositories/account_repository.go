package repositories

import (
	"context"
	"errors"
	"palm/src/entities"
)

// Common repository errors
var (
	ErrAccountNotFound = errors.New("account not found")
)

type AccountRepository interface {
	Create(ctx context.Context, account *entities.Account) error
	GetByID(ctx context.Context, id int64) (*entities.Account, error)
	GetByEmail(ctx context.Context, email string) (*entities.Account, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*entities.Account, error)
}
