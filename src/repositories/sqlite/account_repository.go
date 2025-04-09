package sqlite

import (
	"context"
	"errors"
	"palm/src/config"
	"palm/src/entities"
	"palm/src/repositories"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) repositories.AccountRepository {
	config.Logger.Debug().Msg("Initializing account repository")
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *entities.Account) error {
	config.Logger.Debug().
		Str("email", account.Email).
		Msg("Creating new account")

	err := r.db.WithContext(ctx).Create(account).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Str("email", account.Email).
			Msg("Failed to create account")
	} else {
		config.Logger.Info().
			Uint("id", account.ID).
			Str("email", account.Email).
			Msg("Account created successfully")
	}
	return err
}

func (r *accountRepository) GetByID(ctx context.Context, id uint) (*entities.Account, error) {
	config.Logger.Debug().Uint("id", id).Msg("Getting account by ID")

	var account entities.Account
	err := r.db.WithContext(ctx).First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config.Logger.Warn().Uint("id", id).Msg("Account not found")
			return nil, repositories.ErrAccountNotFound
		}
		config.Logger.Error().Err(err).Uint("id", id).Msg("Error retrieving account")
		return nil, err
	}
	config.Logger.Debug().Uint("id", id).Str("email", account.Email).Msg("Account retrieved successfully")
	return &account, nil
}

func (r *accountRepository) GetByEmail(ctx context.Context, email string) (*entities.Account, error) {
	config.Logger.Debug().Str("email", email).Msg("Getting account by email")

	var account entities.Account
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config.Logger.Warn().Str("email", email).Msg("Account not found")
			return nil, repositories.ErrAccountNotFound
		}
		config.Logger.Error().Err(err).Str("email", email).Msg("Error retrieving account")
		return nil, err
	}
	config.Logger.Debug().Uint("id", account.ID).Str("email", email).Msg("Account retrieved successfully")
	return &account, nil
}

func (r *accountRepository) Delete(ctx context.Context, id uint) error {
	config.Logger.Debug().Uint("id", id).Msg("Deleting account")

	result := r.db.WithContext(ctx).Delete(&entities.Account{}, id)
	if result.Error != nil {
		config.Logger.Error().Err(result.Error).Uint("id", id).Msg("Error deleting account")
		return result.Error
	}
	if result.RowsAffected == 0 {
		config.Logger.Warn().Uint("id", id).Msg("Account not found for deletion")
		return repositories.ErrAccountNotFound
	}
	config.Logger.Info().Uint("id", id).Msg("Account deleted successfully")
	return nil
}

func (r *accountRepository) List(ctx context.Context) ([]*entities.Account, error) {
	config.Logger.Debug().Msg("Listing all accounts")

	var accounts []*entities.Account
	err := r.db.WithContext(ctx).Find(&accounts).Error
	if err != nil {
		config.Logger.Error().Err(err).Msg("Error listing accounts")
		return nil, err
	}
	config.Logger.Debug().Int("count", len(accounts)).Msg("Accounts retrieved successfully")
	return accounts, nil
}
