package services

import (
	"context"
	"errors"
	"fmt"
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
	return &AccountService{repo: repo}
}

// validateAccountType validates that the account type is one of the allowed values
func (s *AccountService) validateAccountType(accountType string) error {
	if accountType != entities.AccountTypeMicrosoft && accountType != entities.AccountTypeGoogle {
		return ErrInvalidAccountType
	}
	return nil
}

func (s *AccountService) CreateAccount(ctx context.Context, email, accountType string) (*entities.Account, error) {
	if err := s.validateAccountType(accountType); err != nil {
		return nil, err
	}

	account := &entities.Account{
		Email:       email,
		AccountType: accountType,
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id int64) (*entities.Account, error) {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

func (s *AccountService) GetAccountByEmail(ctx context.Context, email string) (*entities.Account, error) {
	account, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get account by email: %w", err)
	}
	return account, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	return nil
}

func (s *AccountService) ListAccounts(ctx context.Context) ([]*entities.Account, error) {
	accounts, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	return accounts, nil
}
