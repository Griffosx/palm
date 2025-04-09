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

func TestCreateRecipient(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewRecipientRepository(db)
	service := services.NewRecipientService(repo)
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

	tests := []struct {
		name      string
		recipient *entities.Recipient
		wantErr   bool
		errType   error
	}{
		{
			name: "To Recipient",
			recipient: &entities.Recipient{
				MessageID:     message.MessageID,
				Email:         "to@example.com",
				RecipientType: entities.RecipientTypeTo,
			},
			wantErr: false,
		},
		{
			name: "CC Recipient",
			recipient: &entities.Recipient{
				MessageID:     message.MessageID,
				Email:         "cc@example.com",
				RecipientType: entities.RecipientTypeCc,
			},
			wantErr: false,
		},
		{
			name: "Invalid Recipient Type",
			recipient: &entities.Recipient{
				MessageID:     message.MessageID,
				Email:         "invalid@example.com",
				RecipientType: "InvalidType",
			},
			wantErr: true,
			errType: services.ErrInvalidRecipientType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateRecipient(ctx, tt.recipient)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			assert.NoError(t, err)
			assert.Greater(t, tt.recipient.RecipientID, int64(0), "Recipient ID should be greater than 0")

			// Verify in DB
			var dbRecipient entities.Recipient
			err = db.First(&dbRecipient, tt.recipient.RecipientID).Error
			assert.NoError(t, err)
			assert.Equal(t, tt.recipient.MessageID, dbRecipient.MessageID)
			assert.Equal(t, tt.recipient.Email, dbRecipient.Email)
			assert.Equal(t, tt.recipient.RecipientType, dbRecipient.RecipientType)
		})
	}
}

func TestGetRecipientsByMessageID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewRecipientRepository(db)
	service := services.NewRecipientService(repo)
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

	// Create test recipients
	recipients := []*entities.Recipient{
		{
			MessageID:     message1.MessageID,
			Email:         "to@example.com",
			RecipientType: entities.RecipientTypeTo,
		},
		{
			MessageID:     message1.MessageID,
			Email:         "cc@example.com",
			RecipientType: entities.RecipientTypeCc,
		},
		{
			MessageID:     message2.MessageID,
			Email:         "to2@example.com",
			RecipientType: entities.RecipientTypeTo,
		},
	}

	for _, r := range recipients {
		err := db.Create(r).Error
		require.NoError(t, err)
	}

	t.Run("Get Recipients for Existing Message", func(t *testing.T) {
		foundRecipients, err := service.GetRecipientsByMessageID(ctx, message1.MessageID)
		assert.NoError(t, err)
		assert.Len(t, foundRecipients, 2)

		// Check email addresses
		emails := []string{foundRecipients[0].Email, foundRecipients[1].Email}
		assert.Contains(t, emails, "to@example.com")
		assert.Contains(t, emails, "cc@example.com")

		// Ensure message2's recipient isn't included
		assert.NotContains(t, emails, "to2@example.com")
	})

	t.Run("Get Recipients for Non-existing Message", func(t *testing.T) {
		foundRecipients, err := service.GetRecipientsByMessageID(ctx, 9999)
		assert.NoError(t, err)
		assert.Empty(t, foundRecipients)
	})
}

func TestDeleteRecipientsByMessageID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewRecipientRepository(db)
	service := services.NewRecipientService(repo)
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

	// Create test recipients
	recipients := []*entities.Recipient{
		{
			MessageID:     message1.MessageID,
			Email:         "to@example.com",
			RecipientType: entities.RecipientTypeTo,
		},
		{
			MessageID:     message1.MessageID,
			Email:         "cc@example.com",
			RecipientType: entities.RecipientTypeCc,
		},
		{
			MessageID:     message2.MessageID,
			Email:         "to2@example.com",
			RecipientType: entities.RecipientTypeTo,
		},
	}

	for _, r := range recipients {
		err := db.Create(r).Error
		require.NoError(t, err)
	}

	t.Run("Delete Recipients for Existing Message", func(t *testing.T) {
		err := service.DeleteRecipientsByMessageID(ctx, message1.MessageID)
		assert.NoError(t, err)

		// Verify recipients for message1 are deleted
		var count int64
		db.Model(&entities.Recipient{}).Where("message_id = ?", message1.MessageID).Count(&count)
		assert.Equal(t, int64(0), count)

		// Verify recipients for message2 still exist
		db.Model(&entities.Recipient{}).Where("message_id = ?", message2.MessageID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Delete Recipients for Non-existing Message", func(t *testing.T) {
		err := service.DeleteRecipientsByMessageID(ctx, 9999)
		assert.NoError(t, err)
	})
}
