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

func TestCreateAttachment(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAttachmentRepository(db)
	service := services.NewAttachmentService(repo)
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

	testCases := []struct {
		name       string
		attachment *entities.Attachment
		wantErr    bool
	}{
		{
			name: "PDF Attachment",
			attachment: &entities.Attachment{
				MessageID: message.MessageID,
				Filename:  "test.pdf",
				MimeType:  "application/pdf",
				Size:      1024,
			},
			wantErr: false,
		},
		{
			name: "Image Attachment",
			attachment: &entities.Attachment{
				MessageID: message.MessageID,
				Filename:  "test.jpg",
				MimeType:  "image/jpeg",
				Size:      2048,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.CreateAttachment(ctx, tc.attachment)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Greater(t, tc.attachment.AttachmentID, int64(0), "Attachment ID should be greater than 0")

			// Verify in DB
			var dbAttachment entities.Attachment
			err = db.First(&dbAttachment, tc.attachment.AttachmentID).Error
			assert.NoError(t, err)
			assert.Equal(t, tc.attachment.MessageID, dbAttachment.MessageID)
			assert.Equal(t, tc.attachment.Filename, dbAttachment.Filename)
			assert.Equal(t, tc.attachment.MimeType, dbAttachment.MimeType)
			assert.Equal(t, tc.attachment.Size, dbAttachment.Size)
		})
	}
}

func TestGetAttachment(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAttachmentRepository(db)
	service := services.NewAttachmentService(repo)
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

	t.Run("Existing Attachment", func(t *testing.T) {
		attachment, err := service.GetAttachment(ctx, testAttachment.AttachmentID)
		assert.NoError(t, err)
		assert.NotNil(t, attachment)
		assert.Equal(t, testAttachment.AttachmentID, attachment.AttachmentID)
		assert.Equal(t, testAttachment.MessageID, attachment.MessageID)
		assert.Equal(t, testAttachment.Filename, attachment.Filename)
		assert.Equal(t, testAttachment.MimeType, attachment.MimeType)
		assert.Equal(t, testAttachment.Size, attachment.Size)
	})

	t.Run("Non-existing Attachment", func(t *testing.T) {
		_, err := service.GetAttachment(ctx, 9999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get attachment")
	})
}

func TestGetAttachmentsByMessageID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAttachmentRepository(db)
	service := services.NewAttachmentService(repo)
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

	// Create test attachments
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

	t.Run("Get Attachments for Existing Message", func(t *testing.T) {
		foundAttachments, err := service.GetAttachmentsByMessageID(ctx, message1.MessageID)
		assert.NoError(t, err)
		assert.Len(t, foundAttachments, 2)

		// Check filenames
		filenames := []string{foundAttachments[0].Filename, foundAttachments[1].Filename}
		assert.Contains(t, filenames, "test1.pdf")
		assert.Contains(t, filenames, "test2.jpg")

		// Ensure message2's attachment isn't included
		assert.NotContains(t, filenames, "test3.docx")
	})

	t.Run("Get Attachments for Non-existing Message", func(t *testing.T) {
		foundAttachments, err := service.GetAttachmentsByMessageID(ctx, 9999)
		assert.NoError(t, err)
		assert.Empty(t, foundAttachments)
	})
}

func TestDeleteAttachmentsByMessageID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewAttachmentRepository(db)
	service := services.NewAttachmentService(repo)
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

	// Create test attachments
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
		err := service.DeleteAttachmentsByMessageID(ctx, message1.MessageID)
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
		err := service.DeleteAttachmentsByMessageID(ctx, 9999)
		assert.NoError(t, err)
	})
}
