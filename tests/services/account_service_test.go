package services_test

import (
	"context"
	"palm/src/entities"
	"palm/src/repositories/sqlite"
	"palm/src/services"
	"palm/tests/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	service := services.NewAccountService(repo)
	ctx := context.Background()

	tests := []struct {
		name        string
		email       string
		accountType string
		wantErr     bool
		errType     error
	}{
		{
			name:        "Valid Microsoft Account",
			email:       "test@microsoft.com",
			accountType: entities.AccountTypeMicrosoft,
			wantErr:     false,
		},
		{
			name:        "Valid Google Account",
			email:       "test@google.com",
			accountType: entities.AccountTypeGoogle,
			wantErr:     false,
		},
		{
			name:        "Invalid Account Type",
			email:       "test@example.com",
			accountType: "InvalidType",
			wantErr:     true,
			errType:     services.ErrInvalidAccountType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := service.CreateAccount(ctx, tt.email, tt.accountType)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, account)
			assert.Equal(t, tt.email, account.Email)
			assert.Equal(t, tt.accountType, account.AccountType)
			assert.Greater(t, account.ID, int64(0), "Account ID should be greater than 0")
		})
	}
}

func TestGetAccount(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	service := services.NewAccountService(repo)
	ctx := context.Background()

	// Create test account
	testAccount := &entities.Account{
		Email:       "get-test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := db.Create(testAccount).Error
	require.NoError(t, err)
	require.Greater(t, testAccount.ID, int64(0))

	t.Run("Existing Account", func(t *testing.T) {
		account, err := service.GetAccount(ctx, testAccount.ID)
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, testAccount.ID, account.ID)
		assert.Equal(t, testAccount.Email, account.Email)
		assert.Equal(t, testAccount.AccountType, account.AccountType)
	})

	t.Run("Non-existing Account", func(t *testing.T) {
		_, err := service.GetAccount(ctx, 9999)
		assert.Error(t, err)
		// The error is wrapped, so we can't directly check for ErrAccountNotFound
		assert.Contains(t, err.Error(), "failed to get account")
	})
}

func TestGetAccountByEmail(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	service := services.NewAccountService(repo)
	ctx := context.Background()

	// Create test account
	testAccount := &entities.Account{
		Email:       "get-by-email@example.com",
		AccountType: entities.AccountTypeGoogle,
	}
	err := db.Create(testAccount).Error
	require.NoError(t, err)

	t.Run("Existing Email", func(t *testing.T) {
		account, err := service.GetAccountByEmail(ctx, testAccount.Email)
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, testAccount.ID, account.ID)
		assert.Equal(t, testAccount.Email, account.Email)
		assert.Equal(t, testAccount.AccountType, account.AccountType)
	})

	t.Run("Non-existing Email", func(t *testing.T) {
		_, err := service.GetAccountByEmail(ctx, "nonexistent@example.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get account by email")
	})
}

func TestDeleteAccount(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	service := services.NewAccountService(repo)
	ctx := context.Background()

	// Create test account
	testAccount := &entities.Account{
		Email:       "delete-test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := db.Create(testAccount).Error
	require.NoError(t, err)
	require.Greater(t, testAccount.ID, int64(0))

	t.Run("Delete Existing Account", func(t *testing.T) {
		err := service.DeleteAccount(ctx, testAccount.ID)
		assert.NoError(t, err)

		// Verify account is deleted
		var count int64
		db.Model(&entities.Account{}).Where("id = ?", testAccount.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Delete Non-existing Account", func(t *testing.T) {
		err := service.DeleteAccount(ctx, 9999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete account")
	})
}

func TestListAccounts(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	service := services.NewAccountService(repo)
	ctx := context.Background()

	// Empty list
	t.Run("Empty List", func(t *testing.T) {
		accounts, err := service.ListAccounts(ctx)
		assert.NoError(t, err)
		assert.Empty(t, accounts)
	})

	// Create test accounts
	testAccounts := []*entities.Account{
		{Email: "list-test1@example.com", AccountType: entities.AccountTypeMicrosoft},
		{Email: "list-test2@example.com", AccountType: entities.AccountTypeGoogle},
		{Email: "list-test3@example.com", AccountType: entities.AccountTypeMicrosoft},
	}

	for _, account := range testAccounts {
		err := db.Create(account).Error
		require.NoError(t, err)
	}

	t.Run("List Multiple Accounts", func(t *testing.T) {
		accounts, err := service.ListAccounts(ctx)
		assert.NoError(t, err)
		assert.Len(t, accounts, len(testAccounts))
	})
}
