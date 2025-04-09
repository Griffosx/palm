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
	"gorm.io/gorm"
)

func TestAccountRepository_Create(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	account := &entities.Account{
		Email:       "create-test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}

	err := repo.Create(ctx, account)
	assert.NoError(t, err)
	assert.Greater(t, account.ID, int64(0), "Account ID should be set after create")

	// Verify in DB
	var dbAccount entities.Account
	err = db.First(&dbAccount, account.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, account.Email, dbAccount.Email)
	assert.Equal(t, account.AccountType, dbAccount.AccountType)
}

func TestAccountRepository_GetByID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	// Create a test account
	testAccount := &entities.Account{
		Email:       "get-by-id-test@example.com",
		AccountType: entities.AccountTypeGoogle,
	}
	err := db.Create(testAccount).Error
	require.NoError(t, err)
	require.Greater(t, testAccount.ID, int64(0))

	t.Run("Existing ID", func(t *testing.T) {
		account, err := repo.GetByID(ctx, testAccount.ID)
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, testAccount.ID, account.ID)
		assert.Equal(t, testAccount.Email, account.Email)
		assert.Equal(t, testAccount.AccountType, account.AccountType)
	})

	t.Run("Non-existing ID", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 9999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrAccountNotFound)
	})
}

func TestAccountRepository_GetByEmail(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	// Create a test account
	testAccount := &entities.Account{
		Email:       "get-by-email-test@example.com",
		AccountType: entities.AccountTypeGoogle,
	}
	err := db.Create(testAccount).Error
	require.NoError(t, err)

	t.Run("Existing Email", func(t *testing.T) {
		account, err := repo.GetByEmail(ctx, testAccount.Email)
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, testAccount.ID, account.ID)
		assert.Equal(t, testAccount.Email, account.Email)
		assert.Equal(t, testAccount.AccountType, account.AccountType)
	})

	t.Run("Non-existing Email", func(t *testing.T) {
		_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
		assert.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrAccountNotFound)
	})
}

func TestAccountRepository_Delete(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	// Create a test account
	testAccount := &entities.Account{
		Email:       "delete-test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := db.Create(testAccount).Error
	require.NoError(t, err)
	require.Greater(t, testAccount.ID, int64(0))

	t.Run("Delete Existing Account", func(t *testing.T) {
		err := repo.Delete(ctx, testAccount.ID)
		assert.NoError(t, err)

		// Verify it's deleted
		var found entities.Account
		err = db.First(&found, testAccount.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Delete Non-existing Account", func(t *testing.T) {
		err := repo.Delete(ctx, 9999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrAccountNotFound)
	})
}

func TestAccountRepository_List(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAccountRepository(db)
	ctx := context.Background()

	// List on empty DB
	t.Run("Empty List", func(t *testing.T) {
		accounts, err := repo.List(ctx)
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
		accounts, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, accounts, len(testAccounts))
	})
}
