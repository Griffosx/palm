package sqlite_test

import (
	"context"
	"palm/src/entities"
	"palm/src/repositories"
	"palm/src/repositories/sqlite"
	"palm/tests/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository_Create(t *testing.T) {
	// Setup
	db := utils.SetupTestDB(t)

	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	// Test data
	account := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeGoogle,
	}

	// Test
	err := repo.Create(ctx, account)

	// Assertions
	require.NoError(t, err)
	assert.NotZero(t, account.ID, "ID should be assigned")
}

func TestAccountRepository_GetByID(t *testing.T) {
	// Setup
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	// Create test account
	testAccount := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeGoogle,
	}
	err := repo.Create(ctx, testAccount)
	require.NoError(t, err)
	require.NotZero(t, testAccount.ID)

	// Test: Get existing account
	account, err := repo.GetByID(ctx, testAccount.ID)
	require.NoError(t, err)
	assert.Equal(t, testAccount.ID, account.ID)
	assert.Equal(t, testAccount.Email, account.Email)
	assert.Equal(t, testAccount.AccountType, account.AccountType)

	// Test: Get non-existent account
	account, err = repo.GetByID(ctx, 9999)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repositories.ErrAccountNotFound)
	assert.Nil(t, account)
}

func TestAccountRepository_GetByEmail(t *testing.T) {
	// Setup
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	// Create test account
	testAccount := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := repo.Create(ctx, testAccount)
	require.NoError(t, err)

	// Test: Get existing account
	account, err := repo.GetByEmail(ctx, testAccount.Email)
	require.NoError(t, err)
	assert.Equal(t, testAccount.ID, account.ID)
	assert.Equal(t, testAccount.Email, account.Email)
	assert.Equal(t, testAccount.AccountType, account.AccountType)

	// Test: Get non-existent account
	account, err = repo.GetByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
	assert.ErrorIs(t, err, repositories.ErrAccountNotFound)
	assert.Nil(t, account)
}

func TestAccountRepository_Delete(t *testing.T) {
	// Setup
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	// Create test account
	testAccount := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeGoogle,
	}
	err := repo.Create(ctx, testAccount)
	require.NoError(t, err)

	// Test: Delete existing account
	err = repo.Delete(ctx, testAccount.ID)
	require.NoError(t, err)

	// Verify deletion
	account, err := repo.GetByID(ctx, testAccount.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repositories.ErrAccountNotFound)
	assert.Nil(t, account)

	// Test: Delete non-existent account
	err = repo.Delete(ctx, 9999)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repositories.ErrAccountNotFound)
}

func TestAccountRepository_List(t *testing.T) {
	// Setup
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	// Create test accounts
	accounts := []*entities.Account{
		{Email: "test1@example.com", AccountType: entities.AccountTypeGoogle},
		{Email: "test2@example.com", AccountType: entities.AccountTypeMicrosoft},
		{Email: "test3@example.com", AccountType: entities.AccountTypeGoogle},
	}

	for _, acc := range accounts {
		err := repo.Create(ctx, acc)
		require.NoError(t, err)
	}

	// Test
	result, err := repo.List(ctx)

	// Assertions
	require.NoError(t, err)
	assert.Len(t, result, len(accounts), "Should return all created accounts")

	// Check if all accounts are in the result
	emails := make(map[string]bool)
	for _, acc := range result {
		emails[acc.Email] = true
	}

	for _, acc := range accounts {
		assert.True(t, emails[acc.Email], "Account with email %s should be in results", acc.Email)
	}
}

func TestAccountRepository_UniqueEmail(t *testing.T) {
	// Setup
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	// Create first account
	account1 := &entities.Account{
		Email:       "duplicate@example.com",
		AccountType: entities.AccountTypeGoogle,
	}
	err := repo.Create(ctx, account1)
	require.NoError(t, err)

	// Try to create account with same email
	account2 := &entities.Account{
		Email:       "duplicate@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err = repo.Create(ctx, account2)

	// Should fail due to unique constraint
	assert.Error(t, err, "Creating account with duplicate email should fail")
}
