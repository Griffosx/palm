package sqlite_test

import (
	"context"
	"fmt"
	"palm/src/entities"
	"palm/src/repositories"
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
	result := accountRepo.Create(ctx, account)
	require.NoError(t, result.Error, "Account creation should succeed")
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
	result = messageRepo.Create(ctx, message)
	if result.Error != nil {
		t.Logf("Error creating message: %v", result.Error)

		// Debug: Check if the account still exists using repository
		checkAccount, checkErr := accountRepo.GetByID(ctx, account.ID)
		if checkErr != nil {
			t.Logf("Account lookup after error: %v", checkErr)
		} else {
			t.Logf("Account still exists with ID: %d", checkAccount.ID)
		}
	}

	assert.NoError(t, result.Error)
	assert.NotZero(t, message.ID, "Message ID should be set after creation")
	assert.Equal(t, int64(1), result.RowsAffected, "One row should be affected")
}

// TestMessageRepository_GetByID tests the GetByID method
func TestMessageRepository_GetByID(t *testing.T) {
	// Setup test database
	db := utils.SetupTestDB(t)
	ctx := context.Background()

	// Create account and message repositories
	messageRepo := sqlite.NewMessageRepository(db)

	// Create test account directly in the database
	account := createTestAccount(t, db, "getbyid-test@example.com")

	// Create a test message
	subject := "Test GetByID Subject"
	body := "Test GetByID Body"
	senderName := "Test Sender"
	message := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject,
		Body:        &body,
		SenderEmail: "sender@example.com",
		SenderName:  &senderName,
		IsDraft:     false,
		IsRead:      false,
		Importance:  entities.ImportanceNormal,
	}

	// Save the message to the database
	result := messageRepo.Create(ctx, message)
	require.NoError(t, result.Error)
	require.NotZero(t, message.ID)

	// Retrieve the message by ID
	retrievedMessage, err := messageRepo.GetByID(ctx, message.ID)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, retrievedMessage)
	assert.Equal(t, message.ID, retrievedMessage.ID)
	assert.Equal(t, account.ID, retrievedMessage.AccountID)
	assert.Equal(t, subject, *retrievedMessage.Subject)
	assert.Equal(t, body, *retrievedMessage.Body)
	assert.Equal(t, "sender@example.com", retrievedMessage.SenderEmail)
	assert.Equal(t, senderName, *retrievedMessage.SenderName)
	assert.Equal(t, entities.ImportanceNormal, retrievedMessage.Importance)

	// Test getting a non-existent message
	nonExistentMessage, err := messageRepo.GetByID(ctx, 9999)
	assert.Error(t, err)
	assert.Nil(t, nonExistentMessage)
	assert.ErrorIs(t, err, repositories.ErrMessageNotFound)
}

// TestMessageRepository_Update tests the Update method
func TestMessageRepository_Update(t *testing.T) {
	// Setup test database
	db := utils.SetupTestDB(t)
	ctx := context.Background()

	// Create message repository
	messageRepo := sqlite.NewMessageRepository(db)

	// Create test account
	account := createTestAccount(t, db, "update-test@example.com")

	// Create a test message
	subject := "Original Subject"
	body := "Original Body"
	senderName := "Original Sender"
	message := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject,
		Body:        &body,
		SenderEmail: "sender@example.com",
		SenderName:  &senderName,
		IsDraft:     false,
		IsRead:      false,
		Importance:  entities.ImportanceNormal,
	}

	// Save the message to the database
	result := messageRepo.Create(ctx, message)
	require.NoError(t, result.Error)
	require.NotZero(t, message.ID)

	// Update the message
	updatedSubject := "Updated Subject"
	updatedBody := "Updated Body"
	message.Subject = &updatedSubject
	message.Body = &updatedBody
	message.IsRead = true
	message.Importance = entities.ImportanceHigh

	err := messageRepo.Update(ctx, message)
	require.NoError(t, err)

	// Retrieve the updated message
	updatedMessage, err := messageRepo.GetByID(ctx, message.ID)
	require.NoError(t, err)

	// Assertions
	assert.Equal(t, updatedSubject, *updatedMessage.Subject)
	assert.Equal(t, updatedBody, *updatedMessage.Body)
	assert.True(t, updatedMessage.IsRead)
	assert.Equal(t, entities.ImportanceHigh, updatedMessage.Importance)

	// Test updating a non-existent message
	nonExistentMessage := &entities.Message{
		Model:     gorm.Model{ID: 9999},
		AccountID: account.ID,
		Subject:   &updatedSubject,
		Body:      &updatedBody,
	}

	err = messageRepo.Update(ctx, nonExistentMessage)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repositories.ErrMessageNotFound)
}

// TestMessageRepository_Delete tests the Delete method
func TestMessageRepository_Delete(t *testing.T) {
	// Setup test database
	db := utils.SetupTestDB(t)
	ctx := context.Background()

	// Create message repository
	messageRepo := sqlite.NewMessageRepository(db)

	// Create test account
	account := createTestAccount(t, db, "delete-test@example.com")

	// Create a test message
	subject := "Delete Test Subject"
	message := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject,
		SenderEmail: "sender@example.com",
		IsDraft:     false,
		IsRead:      false,
		Importance:  entities.ImportanceNormal,
	}

	// Save the message to the database
	result := messageRepo.Create(ctx, message)
	require.NoError(t, result.Error)
	require.NotZero(t, message.ID)

	// Delete the message
	err := messageRepo.Delete(ctx, message.ID)
	require.NoError(t, err)

	// Try to retrieve the deleted message
	deletedMessage, err := messageRepo.GetByID(ctx, message.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedMessage)
	assert.ErrorIs(t, err, repositories.ErrMessageNotFound)

	// Test deleting a non-existent message
	err = messageRepo.Delete(ctx, 9999)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repositories.ErrMessageNotFound)
}
