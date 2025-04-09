package services

import (
	"context"
	"errors"
	"fmt"
	"palm/src/config"
	"palm/src/entities"
	"palm/src/repositories"
)

// Custom error types
var (
	ErrInvalidAccountType = errors.New("invalid account type")
	ErrAccountNotFound    = errors.New("account not found")
)

type AccountService struct {
	repo repositories.AccountRepository
}

func NewAccountService(repo repositories.AccountRepository) *AccountService {
	config.Logger.Debug().Msg("Initializing account service")
	return &AccountService{repo: repo}
}

// validateAccountType validates that the account type is one of the allowed values
func (s *AccountService) validateAccountType(accountType string) error {
	config.Logger.Debug().Str("accountType", accountType).Msg("Validating account type")

	if accountType != entities.AccountTypeMicrosoft && accountType != entities.AccountTypeGoogle {
		config.Logger.Warn().
			Str("accountType", accountType).
			Msg("Invalid account type")
		return ErrInvalidAccountType
	}
	return nil
}

func (s *AccountService) CreateAccount(ctx context.Context, email, accountType string) (*entities.Account, error) {
	config.Logger.Info().
		Str("email", email).
		Str("accountType", accountType).
		Msg("Creating new account")

	if err := s.validateAccountType(accountType); err != nil {
		return nil, err
	}

	account := &entities.Account{
		Email:       email,
		AccountType: accountType,
	}

	if err := s.repo.Create(ctx, account); err != nil {
		config.Logger.Error().
			Err(err).
			Str("email", email).
			Str("accountType", accountType).
			Msg("Failed to create account")
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	config.Logger.Info().
		Int64("id", account.ID).
		Str("email", email).
		Str("accountType", accountType).
		Msg("Account created successfully")

	return account, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id int64) (*entities.Account, error) {
	config.Logger.Debug().Int64("id", id).Msg("Getting account by ID")

	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Int64("id", id).
			Msg("Failed to get account")
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	config.Logger.Debug().
		Int64("id", id).
		Str("email", account.Email).
		Str("accountType", account.AccountType).
		Msg("Account retrieved successfully")

	return account, nil
}

func (s *AccountService) GetAccountByEmail(ctx context.Context, email string) (*entities.Account, error) {
	config.Logger.Debug().Str("email", email).Msg("Getting account by email")

	account, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Str("email", email).
			Msg("Failed to get account by email")
		return nil, fmt.Errorf("failed to get account by email: %w", err)
	}

	config.Logger.Debug().
		Int64("id", account.ID).
		Str("email", email).
		Str("accountType", account.AccountType).
		Msg("Account retrieved successfully by email")

	return account, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, id int64) error {
	config.Logger.Info().Int64("id", id).Msg("Deleting account")

	if err := s.repo.Delete(ctx, id); err != nil {
		config.Logger.Error().
			Err(err).
			Int64("id", id).
			Msg("Failed to delete account")
		return fmt.Errorf("failed to delete account: %w", err)
	}

	config.Logger.Info().Int64("id", id).Msg("Account deleted successfully")
	return nil
}

func (s *AccountService) ListAccounts(ctx context.Context) ([]*entities.Account, error) {
	config.Logger.Debug().Msg("Listing all accounts")

	accounts, err := s.repo.List(ctx)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Msg("Failed to list accounts")
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	config.Logger.Debug().
		Int("count", len(accounts)).
		Msg("Accounts listed successfully")

	return accounts, nil
}
