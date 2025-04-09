package sqlite_test

import (
	"context"
	"palm/src/entities"
	"palm/src/repositories/sqlite"
	"palm/tests/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipientRepository_Create(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewRecipientRepository(db)
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

	name := "John Doe"
	recipient := &entities.Recipient{
		MessageID:     message.MessageID,
		Email:         "recipient@example.com",
		Name:          &name,
		RecipientType: entities.RecipientTypeTo,
	}

	err = repo.Create(ctx, recipient)
	assert.NoError(t, err)
	assert.Greater(t, recipient.RecipientID, int64(0), "Recipient ID should be set after create")

	// Verify in DB
	var dbRecipient entities.Recipient
	err = db.First(&dbRecipient, recipient.RecipientID).Error
	assert.NoError(t, err)
	assert.Equal(t, message.MessageID, dbRecipient.MessageID)
	assert.Equal(t, recipient.Email, dbRecipient.Email)
	assert.Equal(t, *recipient.Name, *dbRecipient.Name)
	assert.Equal(t, recipient.RecipientType, dbRecipient.RecipientType)
}

func TestRecipientRepository_GetByMessageID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewRecipientRepository(db)
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

	// Create test recipients
	toName := "To Recipient"
	ccName := "CC Recipient"
	recipients := []*entities.Recipient{
		{
			MessageID:     message.MessageID,
			Email:         "to@example.com",
			Name:          &toName,
			RecipientType: entities.RecipientTypeTo,
		},
		{
			MessageID:     message.MessageID,
			Email:         "cc@example.com",
			Name:          &ccName,
			RecipientType: entities.RecipientTypeCc,
		},
	}

	for _, r := range recipients {
		err := db.Create(r).Error
		require.NoError(t, err)
	}

	t.Run("Get Recipients for Existing Message", func(t *testing.T) {
		fetchedRecipients, err := repo.GetByMessageID(ctx, message.MessageID)
		assert.NoError(t, err)
		assert.Len(t, fetchedRecipients, 2)

		// Verify recipient emails
		emails := []string{fetchedRecipients[0].Email, fetchedRecipients[1].Email}
		assert.Contains(t, emails, "to@example.com")
		assert.Contains(t, emails, "cc@example.com")
	})

	t.Run("Get Recipients for Non-existing Message", func(t *testing.T) {
		fetchedRecipients, err := repo.GetByMessageID(ctx, 9999)
		assert.NoError(t, err) // Should not error for non-existing message
		assert.Empty(t, fetchedRecipients)
	})
}

func TestRecipientRepository_DeleteByMessageID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewRecipientRepository(db)
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

	// Create test recipients for both messages
	recipients := []*entities.Recipient{
		{
			MessageID:     message1.MessageID,
			Email:         "recipient1@example.com",
			RecipientType: entities.RecipientTypeTo,
		},
		{
			MessageID:     message1.MessageID,
			Email:         "recipient2@example.com",
			RecipientType: entities.RecipientTypeCc,
		},
		{
			MessageID:     message2.MessageID,
			Email:         "recipient3@example.com",
			RecipientType: entities.RecipientTypeTo,
		},
	}

	for _, r := range recipients {
		err := db.Create(r).Error
		require.NoError(t, err)
	}

	t.Run("Delete Recipients for Existing Message", func(t *testing.T) {
		err := repo.DeleteByMessageID(ctx, message1.MessageID)
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
		err := repo.DeleteByMessageID(ctx, 9999)
		assert.NoError(t, err) // Should not error for non-existing message ID
	})
}
