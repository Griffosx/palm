package sqlite_test

import (
	"context"
	"fmt"
	"palm/src/entities"
	"palm/src/repositories/sqlite"
	"palm/tests/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// createTestAccount creates a test account for message tests
func createTestAccount(t *testing.T, db *gorm.DB, email string) *entities.Account {
	account := &entities.Account{
		Email:       email,
		AccountType: entities.AccountTypeGoogle,
	}
	result := db.Create(account)
	require.NoError(t, result.Error)
	require.NotZero(t, account.ID, "Account ID should not be zero after creation")

	// Verify account was created by retrieving it
	var savedAccount entities.Account
	err := db.First(&savedAccount, account.ID).Error
	require.NoError(t, err, "Should be able to retrieve the created account")

	return account
}

// TestNewMessageRepository tests the creation of a new message repository
func TestNewMessageRepository(t *testing.T) {
	db := utils.SetupTestDB(t)
	messageRepo := sqlite.NewMessageRepository(db)
	assert.NotNil(t, messageRepo)
}

// TestMessageRepository_Create tests the Create method
func TestMessageRepository_Create(t *testing.T) {
	// Setup test database
	db := utils.SetupTestDB(t)
	ctx := context.Background()

	// Create account repository to manage account operations
	accountRepo := sqlite.NewAccountRepository(db)

	// Create message repository
	messageRepo := sqlite.NewMessageRepository(db)

	// Create test account using repository
	account := &entities.Account{
		Email:       "create-test@example.com",
		AccountType: entities.AccountTypeGoogle,
	}

	// Create account via repository
	err := accountRepo.Create(ctx, account)
	require.NoError(t, err, "Account creation should succeed")
	require.NotZero(t, account.ID, "Account ID should not be zero after creation")

	// Debug: Verify account exists via repository
	fetchedAccount, err := accountRepo.GetByEmail(ctx, account.Email)

	fmt.Printf("fetchedAccount: %+v\n", fetchedAccount)
	require.NoError(t, err, "Should be able to fetch the created account")
	t.Logf("Created account with ID: %d", fetchedAccount.ID)

	subject := "Test Subject"
	body := "Test Body"
	bodyPreview := "Test Preview"
	senderName := "Test Sender"
	receivedTime := time.Now()
	sentTime := time.Now().Add(-time.Hour)
	conversationID := "test-conversation-123"

	message := &entities.Message{
		AccountID:        fetchedAccount.ID,
		Subject:          &subject,
		Body:             &body,
		BodyPreview:      &bodyPreview,
		SenderEmail:      "sender@example.com",
		SenderName:       &senderName,
		ReceivedDatetime: &receivedTime,
		SentDatetime:     &sentTime,
		IsDraft:          false,
		IsRead:           false,
		Importance:       entities.ImportanceNormal,
		ConversationID:   &conversationID,
	}

	// Create message via repository
	err = messageRepo.Create(ctx, message)
	if err != nil {
		t.Logf("Error creating message: %v", err)

		// Debug: Check if the account still exists using repository
		checkAccount, checkErr := accountRepo.GetByID(ctx, account.ID)
		if checkErr != nil {
			t.Logf("Account lookup after error: %v", checkErr)
		} else {
			t.Logf("Account still exists with ID: %d", checkAccount.ID)
		}
	}

	assert.NoError(t, err)
	assert.NotZero(t, message.ID, "Message ID should be set after creation")
}

// // TestMessageRepository_GetByID tests the GetByID method
// func TestMessageRepository_GetByID(t *testing.T) {
// 	// Setup test database
// 	db := utils.SetupTestDB(t)
// 	ctx := context.Background()

// 	// Create test account
// 	account := createTestAccount(t, db, "getbyid-test@example.com")

// 	// Ensure the account has a valid ID
// 	require.NotZero(t, account.ID, "Account must have a valid ID before message creation")

// 	// Create message repository
// 	messageRepo := sqlite.NewMessageRepository(db)

// 	// Create a message to retrieve
// 	subject := "Retrieve Test"
// 	message := &entities.Message{
// 		AccountID:   account.ID,
// 		Subject:     &subject,
// 		SenderEmail: "retrieve@example.com",
// 		IsDraft:     false,
// 		IsRead:      false,
// 		Importance:  entities.ImportanceNormal,
// 	}

// 	err := messageRepo.Create(ctx, message)
// 	require.NoError(t, err)

// 	// Test successful retrieval
// 	t.Run("Existing message", func(t *testing.T) {
// 		retrieved, err := messageRepo.GetByID(ctx, message.MessageID)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, retrieved)
// 		assert.Equal(t, message.MessageID, retrieved.MessageID)
// 		assert.Equal(t, message.AccountID, retrieved.AccountID)
// 		assert.Equal(t, message.SenderEmail, retrieved.SenderEmail)
// 		assert.Equal(t, *message.Subject, *retrieved.Subject)
// 	})

// 	// Test non-existent message
// 	t.Run("Non-existent message", func(t *testing.T) {
// 		nonExistent, err := messageRepo.GetByID(ctx, 9999)
// 		assert.Error(t, err)
// 		assert.Equal(t, repositories.ErrMessageNotFound, err)
// 		assert.Nil(t, nonExistent)
// 	})
// }

// // TestMessageRepository_Update tests the Update method
// func TestMessageRepository_Update(t *testing.T) {
// 	// Setup test database
// 	db := utils.SetupTestDB(t)
// 	ctx := context.Background()

// 	// Create test account
// 	account := createTestAccount(t, db, "update-test@example.com")

// 	// Ensure the account has a valid ID
// 	require.NotZero(t, account.ID, "Account must have a valid ID before message creation")

// 	// Create message repository
// 	messageRepo := sqlite.NewMessageRepository(db)

// 	// Test successful update
// 	t.Run("Existing message", func(t *testing.T) {
// 		// Create a message to update
// 		subject := "Update Test"
// 		message := &entities.Message{
// 			AccountID:   account.ID,
// 			Subject:     &subject,
// 			SenderEmail: "update@example.com",
// 			IsDraft:     true,
// 			IsRead:      false,
// 			Importance:  entities.ImportanceLow,
// 		}

// 		err := messageRepo.Create(ctx, message)
// 		require.NoError(t, err)

// 		// Update the message
// 		newSubject := "Updated Subject"
// 		message.Subject = &newSubject
// 		message.IsRead = true
// 		message.Importance = entities.ImportanceHigh

// 		err = messageRepo.Update(ctx, message)
// 		assert.NoError(t, err)

// 		// Verify update
// 		updated, err := messageRepo.GetByID(ctx, message.MessageID)
// 		assert.NoError(t, err)
// 		assert.Equal(t, newSubject, *updated.Subject)
// 		assert.True(t, updated.IsRead)
// 		assert.Equal(t, entities.ImportanceHigh, updated.Importance)
// 	})

// 	// Test update with non-existent ID
// 	t.Run("Non-existent message", func(t *testing.T) {
// 		nonExistentMsg := &entities.Message{
// 			MessageID:   9999,
// 			AccountID:   account.ID,
// 			SenderEmail: "nonexistent@example.com",
// 			IsDraft:     false,
// 			IsRead:      false,
// 			Importance:  entities.ImportanceNormal,
// 		}

// 		err := messageRepo.Update(ctx, nonExistentMsg)
// 		assert.Error(t, err)
// 		assert.Equal(t, repositories.ErrMessageNotFound, err)
// 	})
// }

// // TestMessageRepository_Delete tests the Delete method
// func TestMessageRepository_Delete(t *testing.T) {
// 	// Setup test database
// 	db := utils.SetupTestDB(t)
// 	ctx := context.Background()

// 	// Create test account
// 	account := createTestAccount(t, db, "delete-test@example.com")

// 	// Ensure the account has a valid ID
// 	require.NotZero(t, account.ID, "Account must have a valid ID before message creation")

// 	// Create message repository
// 	messageRepo := sqlite.NewMessageRepository(db)

// 	// Test successful deletion
// 	t.Run("Existing message", func(t *testing.T) {
// 		// Create a message to delete
// 		subject := "Delete Test"
// 		message := &entities.Message{
// 			AccountID:   account.ID,
// 			Subject:     &subject,
// 			SenderEmail: "delete@example.com",
// 			IsDraft:     false,
// 			IsRead:      false,
// 			Importance:  entities.ImportanceNormal,
// 		}

// 		err := messageRepo.Create(ctx, message)
// 		require.NoError(t, err)

// 		// Delete the message
// 		err = messageRepo.Delete(ctx, message.MessageID)
// 		assert.NoError(t, err)

// 		// Verify deletion
// 		deleted, err := messageRepo.GetByID(ctx, message.MessageID)
// 		assert.Error(t, err)
// 		assert.Equal(t, repositories.ErrMessageNotFound, err)
// 		assert.Nil(t, deleted)
// 	})

// 	// Test delete with non-existent ID
// 	t.Run("Non-existent message", func(t *testing.T) {
// 		err := messageRepo.Delete(ctx, 9999)
// 		assert.Error(t, err)
// 		assert.Equal(t, repositories.ErrMessageNotFound, err)
// 	})
// }
