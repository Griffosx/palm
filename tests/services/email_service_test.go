package services_test

import (
	"context"
	"palm/src/entities"
	"palm/src/repositories"
	"palm/src/repositories/sqlite"
	"palm/src/services"
	"palm/tests/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestAccount creates a test account for email tests
func createTestAccount(t *testing.T, ctx context.Context, accountRepo repositories.AccountRepository, email string) *entities.Account {
	account := &entities.Account{
		Email:       email,
		AccountType: entities.AccountTypeGoogle,
	}
	result := accountRepo.Create(ctx, account)
	require.NoError(t, result.Error)
	require.NotZero(t, account.ID, "Account ID should not be zero after creation")
	return account
}

// createEmailDTO creates a test email DTO
func createEmailDTO(accountID uint, subject string) *services.EmailDTO {
	subjectStr := subject
	bodyStr := "Test body"
	bodyPreviewStr := "Test body preview"
	senderName := "Test Sender"
	receivedTime := time.Now()
	sentTime := time.Now().Add(-time.Hour)
	conversationID := "test-conversation-123"
	recipientName := "Recipient One"

	message := &entities.Message{
		AccountID:        accountID,
		Subject:          &subjectStr,
		Body:             &bodyStr,
		BodyPreview:      &bodyPreviewStr,
		SenderEmail:      "sender@example.com",
		SenderName:       &senderName,
		ReceivedDatetime: &receivedTime,
		SentDatetime:     &sentTime,
		IsDraft:          false,
		IsRead:           false,
		Importance:       entities.ImportanceNormal,
		ConversationID:   &conversationID,
	}

	recipients := []*entities.Recipient{
		{
			Email:         "recipient1@example.com",
			Name:          &recipientName,
			RecipientType: entities.RecipientTypeTo,
		},
	}

	attachments := []*entities.Attachment{
		{
			Filename: "test.txt",
			MimeType: "text/plain",
			Size:     100,
		},
	}

	return &services.EmailDTO{
		Message:     message,
		Recipients:  recipients,
		Attachments: attachments,
	}
}

// TestNewEmailService tests the creation of a new email service
func TestNewEmailService(t *testing.T) {
	db := utils.SetupTestDB(t)
	defer utils.TeardownTestDB(t)

	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)

	emailService := services.NewEmailService(db, messageRepo, recipientRepo, attachmentRepo)
	assert.NotNil(t, emailService)
}

// TestEmailService_ValidateEmail tests the validation logic
func TestEmailService_ValidateEmail(t *testing.T) {
	db := utils.SetupTestDB(t)
	defer utils.TeardownTestDB(t)
	ctx := context.Background()

	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)
	accountRepo := sqlite.NewAccountRepository(db)

	emailService := services.NewEmailService(db, messageRepo, recipientRepo, attachmentRepo)

	// Create test account
	account := createTestAccount(t, ctx, accountRepo, "validate-test@example.com")

	// Test case 1: Valid email
	email := createEmailDTO(account.ID, "Test Validation")
	err := emailService.Create(ctx, email)
	assert.NoError(t, err)

	// Test case 2: Missing message
	invalidEmail := &services.EmailDTO{
		Message: nil,
		Recipients: []*entities.Recipient{
			{
				Email:         "recipient@example.com",
				RecipientType: entities.RecipientTypeTo,
			},
		},
	}
	err = emailService.Create(ctx, invalidEmail)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "message cannot be null")

	// Test case 3: No recipients
	invalidEmail = &services.EmailDTO{
		Message: &entities.Message{
			AccountID:   account.ID,
			SenderEmail: "sender@example.com",
		},
		Recipients: []*entities.Recipient{},
	}
	err = emailService.Create(ctx, invalidEmail)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one recipient is required")

	// Test case 4: Invalid importance
	subjectStr := "Test Importance"
	invalidEmail = &services.EmailDTO{
		Message: &entities.Message{
			AccountID:   account.ID,
			Subject:     &subjectStr,
			SenderEmail: "sender@example.com",
			Importance:  "INVALID_IMPORTANCE",
		},
		Recipients: []*entities.Recipient{
			{
				Email:         "recipient@example.com",
				RecipientType: entities.RecipientTypeTo,
			},
		},
	}
	err = emailService.Create(ctx, invalidEmail)
	assert.Error(t, err)
}

// TestEmailService_Create tests the Create method
func TestEmailService_Create(t *testing.T) {
	db := utils.SetupTestDB(t)
	defer utils.TeardownTestDB(t)
	ctx := context.Background()

	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)
	accountRepo := sqlite.NewAccountRepository(db)

	emailService := services.NewEmailService(db, messageRepo, recipientRepo, attachmentRepo)

	// Create test account
	account := createTestAccount(t, ctx, accountRepo, "create-test@example.com")

	// Create a valid email
	email := createEmailDTO(account.ID, "Test Create Email")

	// Create the email
	err := emailService.Create(ctx, email)
	require.NoError(t, err)
	require.NotZero(t, email.Message.ID)

	// Verify message was created in the database
	createdMessage, err := messageRepo.GetByID(ctx, email.Message.ID)
	require.NoError(t, err)
	assert.Equal(t, "Test Create Email", *createdMessage.Subject)

	// Verify recipients were created
	recipients, err := recipientRepo.GetByMessageID(ctx, email.Message.ID)
	require.NoError(t, err)
	assert.Len(t, recipients, 1)
	assert.Equal(t, "recipient1@example.com", recipients[0].Email)

	// Verify attachments were created
	attachments, err := attachmentRepo.GetByMessageID(ctx, email.Message.ID)
	require.NoError(t, err)
	assert.Len(t, attachments, 1)
	assert.Equal(t, "test.txt", attachments[0].Filename)
}

// TestEmailService_GetByID tests the GetByID method
func TestEmailService_GetByID(t *testing.T) {
	db := utils.SetupTestDB(t)
	defer utils.TeardownTestDB(t)
	ctx := context.Background()

	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)
	accountRepo := sqlite.NewAccountRepository(db)

	emailService := services.NewEmailService(db, messageRepo, recipientRepo, attachmentRepo)

	// Create test account
	account := createTestAccount(t, ctx, accountRepo, "getbyid-test@example.com")

	// Create a valid email
	email := createEmailDTO(account.ID, "Test GetByID Email")

	// Create the email
	err := emailService.Create(ctx, email)
	require.NoError(t, err)

	// Get the email by ID
	retrievedEmail, err := emailService.GetByID(ctx, email.Message.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedEmail)

	// Verify email data
	assert.Equal(t, email.Message.ID, retrievedEmail.Message.ID)
	assert.Equal(t, "Test GetByID Email", *retrievedEmail.Message.Subject)
	assert.Equal(t, account.ID, retrievedEmail.Message.AccountID)
	assert.Len(t, retrievedEmail.Recipients, 1)
	assert.Equal(t, "recipient1@example.com", retrievedEmail.Recipients[0].Email)
	assert.Len(t, retrievedEmail.Attachments, 1)
	assert.Equal(t, "test.txt", retrievedEmail.Attachments[0].Filename)

	// Test getting a non-existent email
	nonExistentEmail, err := emailService.GetByID(ctx, 9999)
	assert.Error(t, err)
	assert.Nil(t, nonExistentEmail)
	assert.ErrorIs(t, err, services.ErrEmailNotFound)
}

// TestEmailService_List tests the List method
func TestEmailService_List(t *testing.T) {
	db := utils.SetupTestDB(t)
	defer utils.TeardownTestDB(t)
	ctx := context.Background()

	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)
	accountRepo := sqlite.NewAccountRepository(db)

	emailService := services.NewEmailService(db, messageRepo, recipientRepo, attachmentRepo)

	// Create test account
	account := createTestAccount(t, ctx, accountRepo, "list-test@example.com")

	// Create multiple emails
	for i := 0; i < 5; i++ {
		email := createEmailDTO(account.ID, "Test List Email "+string(rune('A'+i)))
		err := emailService.Create(ctx, email)
		require.NoError(t, err)
	}

	// Test listing emails with pagination
	result, err := emailService.List(ctx, account.ID, 2, 1)
	require.NoError(t, err)
	assert.Len(t, result.Emails, 2)
	assert.Equal(t, int64(5), result.TotalCount)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 2, result.PageSize)
	assert.Equal(t, 3, result.TotalPages)

	// Test second page
	result, err = emailService.List(ctx, account.ID, 2, 2)
	require.NoError(t, err)
	assert.Len(t, result.Emails, 2)

	// Test last page (should have 1 item)
	result, err = emailService.List(ctx, account.ID, 2, 3)
	require.NoError(t, err)
	assert.Len(t, result.Emails, 1)

	// Test invalid page size
	_, err = emailService.List(ctx, account.ID, 0, 1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, services.ErrInvalidPageSize)

	// Test empty account
	emptyAccount := createTestAccount(t, ctx, accountRepo, "empty-account@example.com")
	result, err = emailService.List(ctx, emptyAccount.ID, 10, 1)
	require.NoError(t, err)
	assert.Len(t, result.Emails, 0)
	assert.Equal(t, int64(0), result.TotalCount)
}

// TestEmailService_Delete tests the Delete method
func TestEmailService_Delete(t *testing.T) {
	db := utils.SetupTestDB(t)
	defer utils.TeardownTestDB(t)
	ctx := context.Background()

	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)
	accountRepo := sqlite.NewAccountRepository(db)

	emailService := services.NewEmailService(db, messageRepo, recipientRepo, attachmentRepo)

	// Create test account
	account := createTestAccount(t, ctx, accountRepo, "delete-test@example.com")

	// Create an email
	email := createEmailDTO(account.ID, "Test Delete Email")
	err := emailService.Create(ctx, email)
	require.NoError(t, err)

	// Verify it was created
	messageID := email.Message.ID
	_, err = emailService.GetByID(ctx, messageID)
	require.NoError(t, err)

	// Verify recipients and attachments exist
	recipients, err := recipientRepo.GetByMessageID(ctx, messageID)
	require.NoError(t, err)
	assert.NotEmpty(t, recipients)

	attachments, err := attachmentRepo.GetByMessageID(ctx, messageID)
	require.NoError(t, err)
	assert.NotEmpty(t, attachments)

	// Delete the email
	err = emailService.Delete(ctx, int64(messageID))
	require.NoError(t, err)

	// Verify email was deleted
	_, err = emailService.GetByID(ctx, messageID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, services.ErrEmailNotFound)

	// Verify recipients were deleted
	recipients, err = recipientRepo.GetByMessageID(ctx, messageID)
	require.NoError(t, err)
	assert.Empty(t, recipients)

	// Verify attachments were deleted
	attachments, err = attachmentRepo.GetByMessageID(ctx, messageID)
	require.NoError(t, err)
	assert.Empty(t, attachments)

	// Test deleting non-existent email
	err = emailService.Delete(ctx, 9999)
	assert.Error(t, err)
	assert.ErrorIs(t, err, services.ErrEmailNotFound)
}

// TestEmailService_ListCount tests the ListCount method
func TestEmailService_ListCount(t *testing.T) {
	db := utils.SetupTestDB(t)
	defer utils.TeardownTestDB(t)
	ctx := context.Background()

	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)
	accountRepo := sqlite.NewAccountRepository(db)

	emailService := services.NewEmailService(db, messageRepo, recipientRepo, attachmentRepo)

	// Create test account
	account := createTestAccount(t, ctx, accountRepo, "count-test@example.com")

	// Initially should have no emails
	count, err := emailService.ListCount(ctx, account.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Create multiple emails
	for i := 0; i < 3; i++ {
		email := createEmailDTO(account.ID, "Test Count Email "+string(rune('A'+i)))
		err := emailService.Create(ctx, email)
		require.NoError(t, err)
	}

	// Should now have 3 emails
	count, err = emailService.ListCount(ctx, account.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// Different account should have different count
	emptyAccount := createTestAccount(t, ctx, accountRepo, "empty-count@example.com")
	count, err = emailService.ListCount(ctx, emptyAccount.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestEmailService_CreateWithNonExistentAccount tests creating an email with a non-existent account
func TestEmailService_CreateWithNonExistentAccount(t *testing.T) {
	db := utils.SetupTestDB(t)
	defer utils.TeardownTestDB(t)
	ctx := context.Background()

	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)
	accountRepo := sqlite.NewAccountRepository(db)

	emailService := services.NewEmailService(db, messageRepo, recipientRepo, attachmentRepo)

	// Create a valid test account first
	account := createTestAccount(t, ctx, accountRepo, "real-account@example.com")

	// Create a successful email with the valid account to verify our count later
	validEmail := createEmailDTO(account.ID, "Valid Email")
	err := emailService.Create(ctx, validEmail)
	require.NoError(t, err)

	// Initial count should be 1
	initialCount, err := emailService.ListCount(ctx, account.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), initialCount)

	// Create an email with a non-existent account ID
	nonExistentAccountID := uint(9999)
	email := createEmailDTO(nonExistentAccountID, "Test Non-Existent Account")

	// Attempt to create the email
	err = emailService.Create(ctx, email)

	// Verify proper error is returned
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create email")
	assert.ErrorIs(t, err, services.ErrEmailCreationFailed)

	// Count emails after the failed attempt
	// This verifies that no emails were created during the failed attempt
	finalCount, err := emailService.ListCount(ctx, account.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), finalCount, "No additional emails should have been created after the failed attempt")

	// If possible, also verify no records exist for the non-existent account
	var messageCount int64
	err = db.Model(&entities.Message{}).Where("account_id = ?", nonExistentAccountID).Count(&messageCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), messageCount, "No messages should exist for the non-existent account")
}
