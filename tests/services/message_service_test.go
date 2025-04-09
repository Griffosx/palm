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

func TestCreateMessage(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewMessageRepository(db)
	service := services.NewMessageService(repo)
	ctx := context.Background()

	// Create test account
	account := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := db.Create(account).Error
	require.NoError(t, err)

	tests := []struct {
		name    string
		message *entities.Message
		wantErr bool
		errType error
	}{
		{
			name: "Valid Message",
			message: &entities.Message{
				AccountID:   account.ID,
				SenderEmail: "sender@example.com",
				Importance:  entities.ImportanceNormal,
			},
			wantErr: false,
		},
		{
			name: "Valid Message with High Importance",
			message: &entities.Message{
				AccountID:   account.ID,
				SenderEmail: "sender@example.com",
				Importance:  entities.ImportanceHigh,
			},
			wantErr: false,
		},
		{
			name: "Invalid Importance",
			message: &entities.Message{
				AccountID:   account.ID,
				SenderEmail: "sender@example.com",
				Importance:  "InvalidImportance",
			},
			wantErr: true,
			errType: services.ErrInvalidImportance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateMessage(ctx, tt.message)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			assert.NoError(t, err)
			assert.Greater(t, tt.message.MessageID, int64(0), "Message ID should be greater than 0")

			// Verify in DB
			var dbMessage entities.Message
			err = db.First(&dbMessage, tt.message.MessageID).Error
			assert.NoError(t, err)
			assert.Equal(t, tt.message.AccountID, dbMessage.AccountID)
			assert.Equal(t, tt.message.SenderEmail, dbMessage.SenderEmail)
			assert.Equal(t, tt.message.Importance, dbMessage.Importance)
		})
	}
}

func TestGetMessage(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewMessageRepository(db)
	service := services.NewMessageService(repo)
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
	testMessage := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject,
		SenderEmail: "sender@example.com",
		Importance:  entities.ImportanceNormal,
	}
	err = db.Create(testMessage).Error
	require.NoError(t, err)
	require.Greater(t, testMessage.MessageID, int64(0))

	t.Run("Existing Message", func(t *testing.T) {
		message, err := service.GetMessage(ctx, testMessage.MessageID)
		assert.NoError(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, testMessage.MessageID, message.MessageID)
		assert.Equal(t, testMessage.AccountID, message.AccountID)
		assert.Equal(t, *testMessage.Subject, *message.Subject)
		assert.Equal(t, testMessage.SenderEmail, message.SenderEmail)
		assert.Equal(t, testMessage.Importance, message.Importance)
	})

	t.Run("Non-existing Message", func(t *testing.T) {
		_, err := service.GetMessage(ctx, 9999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get message")
	})
}

func TestUpdateMessage(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewMessageRepository(db)
	service := services.NewMessageService(repo)
	ctx := context.Background()

	// Create test account
	account := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := db.Create(account).Error
	require.NoError(t, err)

	// Create test message
	subject := "Original Subject"
	testMessage := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject,
		SenderEmail: "sender@example.com",
		Importance:  entities.ImportanceNormal,
		IsRead:      false,
	}
	err = db.Create(testMessage).Error
	require.NoError(t, err)
	require.Greater(t, testMessage.MessageID, int64(0))

	t.Run("Update Valid Message", func(t *testing.T) {
		newSubject := "Updated Subject"
		testMessage.Subject = &newSubject
		testMessage.IsRead = true
		testMessage.Importance = entities.ImportanceHigh

		err := service.UpdateMessage(ctx, testMessage)
		assert.NoError(t, err)

		// Verify in DB
		var updatedMessage entities.Message
		err = db.First(&updatedMessage, testMessage.MessageID).Error
		assert.NoError(t, err)
		assert.Equal(t, *testMessage.Subject, *updatedMessage.Subject)
		assert.Equal(t, testMessage.IsRead, updatedMessage.IsRead)
		assert.Equal(t, testMessage.Importance, updatedMessage.Importance)
	})

	t.Run("Update with Invalid Importance", func(t *testing.T) {
		testMessage.Importance = "InvalidImportance"

		err := service.UpdateMessage(ctx, testMessage)
		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidImportance)
	})
}

func TestDeleteMessage(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewMessageRepository(db)
	service := services.NewMessageService(repo)
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
	testMessage := &entities.Message{
		AccountID:   account.ID,
		Subject:     &subject,
		SenderEmail: "sender@example.com",
		Importance:  entities.ImportanceNormal,
	}
	err = db.Create(testMessage).Error
	require.NoError(t, err)
	require.Greater(t, testMessage.MessageID, int64(0))

	t.Run("Delete Existing Message", func(t *testing.T) {
		err := service.DeleteMessage(ctx, testMessage.MessageID)
		assert.NoError(t, err)

		// Verify message is deleted
		var count int64
		db.Model(&entities.Message{}).Where("message_id = ?", testMessage.MessageID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Delete Non-existing Message", func(t *testing.T) {
		err := service.DeleteMessage(ctx, 9999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete message")
	})
}
