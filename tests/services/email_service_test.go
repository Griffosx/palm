package services_test

import (
	"context"
	"palm/src/config"
	"palm/src/entities"
	"palm/src/repositories/sqlite"
	"palm/src/services"
	"palm/tests/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupEmailServiceTest(t *testing.T) (*services.EmailService, *gorm.DB) {
	config.InitLogger() // Initialize logger for service usage
	db := utils.SetupTestDB(t)

	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)

	emailService := services.NewEmailService(db, messageRepo, recipientRepo, attachmentRepo)

	// Create a dummy account for testing relationships
	account := entities.Account{Email: "test@example.com", AccountType: entities.AccountTypeMicrosoft}
	require.NoError(t, db.Create(&account).Error)

	return emailService, db
}

// Helper to create a valid email DTO for tests
func createValidEmailDTO(accountID int64, subject string) *services.EmailDTO {
	msgSubject := subject
	msgBody := "Test Body" // Need a variable to take its address
	now := time.Now()
	return &services.EmailDTO{
		Message: &entities.Message{
			AccountID:        accountID,
			Subject:          &msgSubject,
			Body:             &msgBody, // Use pointer
			SenderEmail:      "sender@example.com",
			ReceivedDatetime: &now, // Correct field name, use pointer
			SentDatetime:     &now, // Correct field name, use pointer
			Importance:       entities.ImportanceNormal,
		},
		Recipients: []*entities.Recipient{
			{Email: "recipient1@example.com", RecipientType: entities.RecipientTypeTo},
		},
		Attachments: []*entities.Attachment{
			{Filename: "file1.txt", MimeType: "text/plain", Size: 100},
		},
	}
}

func TestEmailService_Create_Success(t *testing.T) {
	emailService, db := setupEmailServiceTest(t)
	ctx := context.Background()

	// Get the created Account ID
	var account entities.Account
	require.NoError(t, db.First(&account).Error)

	emailDTO := createValidEmailDTO(account.ID, "Test Subject Create Success")

	err := emailService.Create(ctx, emailDTO)
	require.NoError(t, err)
	require.NotZero(t, emailDTO.Message.MessageID) // Ensure message ID is populated

	// Verify data in DB
	var msg entities.Message
	err = db.First(&msg, emailDTO.Message.MessageID).Error
	require.NoError(t, err)
	assert.Equal(t, *emailDTO.Message.Subject, *msg.Subject)
	assert.Equal(t, account.ID, msg.AccountID)
	assert.Equal(t, emailDTO.Message.Importance, msg.Importance) // Check default importance

	var recipients []*entities.Recipient
	err = db.Where("message_id = ?", msg.MessageID).Find(&recipients).Error
	require.NoError(t, err)
	assert.Len(t, recipients, 1)
	assert.Equal(t, emailDTO.Recipients[0].Email, recipients[0].Email)
	assert.Equal(t, msg.MessageID, recipients[0].MessageID)

	var attachments []*entities.Attachment
	err = db.Where("message_id = ?", msg.MessageID).Find(&attachments).Error
	require.NoError(t, err)
	assert.Len(t, attachments, 1)
	assert.Equal(t, emailDTO.Attachments[0].Filename, attachments[0].Filename)
	assert.Equal(t, emailDTO.Attachments[0].MimeType, attachments[0].MimeType)
	assert.Equal(t, msg.MessageID, attachments[0].MessageID)
}

func TestEmailService_Create_ValidationErrors(t *testing.T) {
	emailService, db := setupEmailServiceTest(t)
	ctx := context.Background()

	// Get Account ID
	var account entities.Account
	require.NoError(t, db.First(&account).Error)

	// Case 1: Missing Message
	emailDTO1 := createValidEmailDTO(account.ID, "Validation Test 1")
	emailDTO1.Message = nil
	err1 := emailService.Create(ctx, emailDTO1)
	require.Error(t, err1)
	assert.Contains(t, err1.Error(), "message cannot be null")

	// Case 2: Missing Recipients
	emailDTO2 := createValidEmailDTO(account.ID, "Validation Test 2")
	emailDTO2.Recipients = []*entities.Recipient{}
	err2 := emailService.Create(ctx, emailDTO2)
	require.Error(t, err2)
	assert.Contains(t, err2.Error(), "at least one recipient is required")

	// Case 3: Invalid Importance
	emailDTO3 := createValidEmailDTO(account.ID, "Validation Test 3")
	emailDTO3.Message.Importance = "invalid_importance"
	err3 := emailService.Create(ctx, emailDTO3)
	require.Error(t, err3)
	assert.Contains(t, err3.Error(), "invalid importance value")

	// Case 4: Default Importance is Set
	emailDTO4 := createValidEmailDTO(account.ID, "Validation Test 4")
	emailDTO4.Message.Importance = "" // Explicitly empty
	err4 := emailService.Create(ctx, emailDTO4)
	require.NoError(t, err4) // Should succeed and set default
	var msg entities.Message
	require.NoError(t, db.First(&msg, emailDTO4.Message.MessageID).Error)
	assert.Equal(t, entities.ImportanceNormal, msg.Importance)
}

func TestEmailService_GetByID_Success(t *testing.T) {
	emailService, db := setupEmailServiceTest(t)
	ctx := context.Background()

	// Create an email first
	var account entities.Account
	require.NoError(t, db.First(&account).Error)
	emailToCreate := createValidEmailDTO(account.ID, "Test Subject Get Success")
	require.NoError(t, emailService.Create(ctx, emailToCreate))
	createdMessageID := emailToCreate.Message.MessageID

	// Get the email
	retrievedEmail, err := emailService.GetByID(ctx, createdMessageID)
	require.NoError(t, err)
	require.NotNil(t, retrievedEmail)

	// Assertions
	assert.Equal(t, createdMessageID, retrievedEmail.Message.MessageID)
	assert.Equal(t, *emailToCreate.Message.Subject, *retrievedEmail.Message.Subject)
	assert.Equal(t, account.ID, retrievedEmail.Message.AccountID)
	require.Len(t, retrievedEmail.Recipients, 1)
	assert.Equal(t, emailToCreate.Recipients[0].Email, retrievedEmail.Recipients[0].Email)
	require.Len(t, retrievedEmail.Attachments, 1)
	assert.Equal(t, emailToCreate.Attachments[0].Filename, retrievedEmail.Attachments[0].Filename)
	assert.Equal(t, emailToCreate.Attachments[0].MimeType, retrievedEmail.Attachments[0].MimeType)
}

func TestEmailService_GetByID_NotFound(t *testing.T) {
	emailService, _ := setupEmailServiceTest(t)
	ctx := context.Background()

	_, err := emailService.GetByID(ctx, 99999) // Non-existent ID
	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrEmailNotFound)
}

func TestEmailService_ListCount_Success(t *testing.T) {
	emailService, db := setupEmailServiceTest(t)
	ctx := context.Background()

	var account1, account2 entities.Account
	require.NoError(t, db.Where("email = ?", "test@example.com").First(&account1).Error)
	account2 = entities.Account{Email: "other@example.com", AccountType: entities.AccountTypeGoogle}
	require.NoError(t, db.Create(&account2).Error)

	// Create emails for account1
	require.NoError(t, emailService.Create(ctx, createValidEmailDTO(account1.ID, "ListCount Email 1")))
	require.NoError(t, emailService.Create(ctx, createValidEmailDTO(account1.ID, "ListCount Email 2")))
	// Create email for account2
	require.NoError(t, emailService.Create(ctx, createValidEmailDTO(account2.ID, "ListCount Other Email")))

	// Count for account1
	count1, err1 := emailService.ListCount(ctx, account1.ID)
	require.NoError(t, err1)
	assert.Equal(t, int64(2), count1)

	// Count for account2
	count2, err2 := emailService.ListCount(ctx, account2.ID)
	require.NoError(t, err2)
	assert.Equal(t, int64(1), count2)

	// Count for non-existent account
	count3, err3 := emailService.ListCount(ctx, 9999)
	require.NoError(t, err3) // Should not error, just return 0
	assert.Equal(t, int64(0), count3)
}

func TestEmailService_List_Success(t *testing.T) {
	emailService, db := setupEmailServiceTest(t)
	ctx := context.Background()

	var account entities.Account
	require.NoError(t, db.First(&account).Error)

	// Create multiple emails
	email1 := createValidEmailDTO(account.ID, "List Email 1")
	nowMinus2h := time.Now().Add(-2 * time.Hour)
	email1.Message.ReceivedDatetime = &nowMinus2h // Ensure order
	require.NoError(t, emailService.Create(ctx, email1))

	email2 := createValidEmailDTO(account.ID, "List Email 2")
	nowMinus1h := time.Now().Add(-1 * time.Hour)
	email2.Message.ReceivedDatetime = &nowMinus1h
	require.NoError(t, emailService.Create(ctx, email2))

	email3 := createValidEmailDTO(account.ID, "List Email 3")
	now := time.Now()
	email3.Message.ReceivedDatetime = &now
	require.NoError(t, emailService.Create(ctx, email3))

	// List emails (page 1, size 2)
	result, err := emailService.List(ctx, account.ID, 2, 1)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(3), result.TotalCount)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 2, result.PageSize)
	assert.Equal(t, 2, result.TotalPages)
	require.Len(t, result.Emails, 2)

	// Check order (newest first)
	assert.Equal(t, *email3.Message.Subject, *result.Emails[0].Message.Subject)
	assert.Equal(t, *email2.Message.Subject, *result.Emails[1].Message.Subject)

	// Check recipients and attachments are loaded
	assert.NotEmpty(t, result.Emails[0].Recipients)
	assert.NotEmpty(t, result.Emails[0].Attachments)
	assert.NotEmpty(t, result.Emails[1].Recipients)
	assert.NotEmpty(t, result.Emails[1].Attachments)

	// List emails (page 2, size 2)
	result2, err2 := emailService.List(ctx, account.ID, 2, 2)
	require.NoError(t, err2)
	require.NotNil(t, result2)

	assert.Equal(t, int64(3), result2.TotalCount)
	assert.Equal(t, 2, result2.Page)
	assert.Equal(t, 2, result2.PageSize)
	assert.Equal(t, 2, result2.TotalPages)
	require.Len(t, result2.Emails, 1)
	assert.Equal(t, *email1.Message.Subject, *result2.Emails[0].Message.Subject)
}

func TestEmailService_List_Empty(t *testing.T) {
	emailService, db := setupEmailServiceTest(t)
	ctx := context.Background()

	var account entities.Account
	require.NoError(t, db.First(&account).Error)

	result, err := emailService.List(ctx, account.ID, 10, 1)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(0), result.TotalCount)
	assert.Len(t, result.Emails, 0)
	assert.Equal(t, 1, result.Page) // Defaults or calculated values
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, 1, result.TotalPages) // Should be 1 even if 0 results
}

func TestEmailService_List_InvalidPageSize(t *testing.T) {
	emailService, db := setupEmailServiceTest(t)
	ctx := context.Background()

	var account entities.Account
	require.NoError(t, db.First(&account).Error)

	_, err := emailService.List(ctx, account.ID, 0, 1) // Page size 0
	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrInvalidPageSize)

	_, err = emailService.List(ctx, account.ID, 101, 1) // Page size > 100
	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrInvalidPageSize)
}

func TestEmailService_Delete_Success(t *testing.T) {
	emailService, db := setupEmailServiceTest(t)
	ctx := context.Background()

	// Create an email first
	var account entities.Account
	require.NoError(t, db.First(&account).Error)
	emailToCreate := createValidEmailDTO(account.ID, "Test Subject Delete Success")
	require.NoError(t, emailService.Create(ctx, emailToCreate))
	createdMessageID := emailToCreate.Message.MessageID

	// Verify it exists
	_, err := emailService.GetByID(ctx, createdMessageID)
	require.NoError(t, err)

	// Delete the email
	err = emailService.Delete(ctx, createdMessageID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = emailService.GetByID(ctx, createdMessageID)
	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrEmailNotFound)

	// Verify related data is deleted (cascade or explicit delete)
	var recipients []*entities.Recipient
	err = db.Where("message_id = ?", createdMessageID).Find(&recipients).Error
	require.NoError(t, err) // Should not error, just find nothing
	assert.Len(t, recipients, 0)

	var attachments []*entities.Attachment
	err = db.Where("message_id = ?", createdMessageID).Find(&attachments).Error
	require.NoError(t, err)
	assert.Len(t, attachments, 0)

	var msg entities.Message
	err = db.First(&msg, createdMessageID).Error
	require.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestEmailService_Delete_NotFound(t *testing.T) {
	emailService, _ := setupEmailServiceTest(t)
	ctx := context.Background()

	err := emailService.Delete(ctx, 99999) // Non-existent ID
	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrEmailNotFound)
}
