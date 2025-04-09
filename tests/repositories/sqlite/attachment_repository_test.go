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

func TestAttachmentRepository_Create(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAttachmentRepository(db)
	ctx := context.Background()

	// Create test account
	account := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := db.Create(account).Error
	require.NoError(t, err)

	// Create test message
	subject := "Test Subject"
	message := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject,
		SenderEmail: "sender@example.com",
		Importance:  entities.ImportanceNormal,
	}
	err = db.Create(message).Error
	require.NoError(t, err)

	attachment := &entities.Attachment{
		MessageID: message.MessageID,
		Filename:  "test.pdf",
		MimeType:  "application/pdf",
		Size:      1024,
	}

	err = repo.Create(ctx, attachment)
	assert.NoError(t, err)
	assert.Greater(t, attachment.AttachmentID, int64(0), "Attachment ID should be set after create")

	// Verify in DB
	var dbAttachment entities.Attachment
	err = db.First(&dbAttachment, attachment.AttachmentID).Error
	assert.NoError(t, err)
	assert.Equal(t, message.MessageID, dbAttachment.MessageID)
	assert.Equal(t, attachment.Filename, dbAttachment.Filename)
	assert.Equal(t, attachment.MimeType, dbAttachment.MimeType)
	assert.Equal(t, attachment.Size, dbAttachment.Size)
}

func TestAttachmentRepository_GetByID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAttachmentRepository(db)
	ctx := context.Background()

	// Create test account
	account := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := db.Create(account).Error
	require.NoError(t, err)

	// Create test message
	subject := "Test Subject"
	message := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject,
		SenderEmail: "sender@example.com",
		Importance:  entities.ImportanceNormal,
	}
	err = db.Create(message).Error
	require.NoError(t, err)

	// Create test attachment
	testAttachment := &entities.Attachment{
		MessageID: message.MessageID,
		Filename:  "test.pdf",
		MimeType:  "application/pdf",
		Size:      1024,
	}
	err = db.Create(testAttachment).Error
	require.NoError(t, err)
	require.Greater(t, testAttachment.AttachmentID, int64(0))

	t.Run("Existing ID", func(t *testing.T) {
		attachment, err := repo.GetByID(ctx, testAttachment.AttachmentID)
		assert.NoError(t, err)
		assert.NotNil(t, attachment)
		assert.Equal(t, testAttachment.AttachmentID, attachment.AttachmentID)
		assert.Equal(t, testAttachment.MessageID, attachment.MessageID)
		assert.Equal(t, testAttachment.Filename, attachment.Filename)
		assert.Equal(t, testAttachment.MimeType, attachment.MimeType)
		assert.Equal(t, testAttachment.Size, attachment.Size)
	})

	t.Run("Non-existing ID", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 9999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrAttachmentNotFound)
	})
}

func TestAttachmentRepository_GetByMessageID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAttachmentRepository(db)
	ctx := context.Background()

	// Create test account
	account := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := db.Create(account).Error
	require.NoError(t, err)

	// Create test message
	subject := "Test Subject"
	message := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject,
		SenderEmail: "sender@example.com",
		Importance:  entities.ImportanceNormal,
	}
	err = db.Create(message).Error
	require.NoError(t, err)

	// Create test attachments
	attachments := []*entities.Attachment{
		{
			MessageID: message.MessageID,
			Filename:  "test1.pdf",
			MimeType:  "application/pdf",
			Size:      1024,
		},
		{
			MessageID: message.MessageID,
			Filename:  "test2.jpg",
			MimeType:  "image/jpeg",
			Size:      2048,
		},
	}

	for _, a := range attachments {
		err := db.Create(a).Error
		require.NoError(t, err)
	}

	t.Run("Get Attachments for Existing Message", func(t *testing.T) {
		fetchedAttachments, err := repo.GetByMessageID(ctx, message.MessageID)
		assert.NoError(t, err)
		assert.Len(t, fetchedAttachments, 2)

		// Verify attachment filenames
		filenames := []string{fetchedAttachments[0].Filename, fetchedAttachments[1].Filename}
		assert.Contains(t, filenames, "test1.pdf")
		assert.Contains(t, filenames, "test2.jpg")
	})

	t.Run("Get Attachments for Non-existing Message", func(t *testing.T) {
		fetchedAttachments, err := repo.GetByMessageID(ctx, 9999)
		assert.NoError(t, err) // Should not error for non-existing message
		assert.Empty(t, fetchedAttachments)
	})
}

func TestAttachmentRepository_DeleteByMessageID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAttachmentRepository(db)
	ctx := context.Background()

	// Create test account
	account := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := db.Create(account).Error
	require.NoError(t, err)

	// Create test messages
	subject1 := "Test Subject 1"
	subject2 := "Test Subject 2"
	message1 := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject1,
		SenderEmail: "sender1@example.com",
		Importance:  entities.ImportanceNormal,
	}
	message2 := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject2,
		SenderEmail: "sender2@example.com",
		Importance:  entities.ImportanceNormal,
	}
	err = db.Create(message1).Error
	require.NoError(t, err)
	err = db.Create(message2).Error
	require.NoError(t, err)

	// Create test attachments for both messages
	attachments := []*entities.Attachment{
		{
			MessageID: message1.MessageID,
			Filename:  "test1.pdf",
			MimeType:  "application/pdf",
			Size:      1024,
		},
		{
			MessageID: message1.MessageID,
			Filename:  "test2.jpg",
			MimeType:  "image/jpeg",
			Size:      2048,
		},
		{
			MessageID: message2.MessageID,
			Filename:  "test3.docx",
			MimeType:  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			Size:      3072,
		},
	}

	for _, a := range attachments {
		err := db.Create(a).Error
		require.NoError(t, err)
	}

	t.Run("Delete Attachments for Existing Message", func(t *testing.T) {
		err := repo.DeleteByMessageID(ctx, message1.MessageID)
		assert.NoError(t, err)

		// Verify attachments for message1 are deleted
		var count int64
		db.Model(&entities.Attachment{}).Where("message_id = ?", message1.MessageID).Count(&count)
		assert.Equal(t, int64(0), count)

		// Verify attachments for message2 still exist
		db.Model(&entities.Attachment{}).Where("message_id = ?", message2.MessageID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Delete Attachments for Non-existing Message", func(t *testing.T) {
		err := repo.DeleteByMessageID(ctx, 9999)
		assert.NoError(t, err) // Should not error for non-existing message ID
	})
}
