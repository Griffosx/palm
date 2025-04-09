package repositories

import (
	"context"
	"errors"
	"palm/src/entities"

	"gorm.io/gorm"
)

// Common repository errors
var (
	ErrAccountNotFound = errors.New("account not found")
)

type AccountRepository interface {
	Create(ctx context.Context, account *entities.Account) *gorm.DB
	GetByID(ctx context.Context, id uint) (*entities.Account, error)
	GetByEmail(ctx context.Context, email string) (*entities.Account, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*entities.Account, error)
}
