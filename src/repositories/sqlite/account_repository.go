package sqlite

import (
	"context"
	"errors"
	"palm/src/entities"
	"palm/src/repositories"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) repositories.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *entities.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) GetByID(ctx context.Context, id int64) (*entities.Account, error) {
	var account entities.Account
	err := r.db.WithContext(ctx).First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetByEmail(ctx context.Context, email string) (*entities.Account, error) {
	var account entities.Account
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&entities.Account{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repositories.ErrAccountNotFound
	}
	return nil
}

func (r *accountRepository) List(ctx context.Context) ([]*entities.Account, error) {
	var accounts []*entities.Account
	err := r.db.WithContext(ctx).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
