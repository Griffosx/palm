package sqlite_test

import (
	"context"
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

func TestMessageRepository_Create(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewMessageRepository(db)
	ctx := context.Background()

	// Create test account
	account := &entities.Account{
		Email:       "test@example.com",
		AccountType: entities.AccountTypeMicrosoft,
	}
	err := db.Create(account).Error
	require.NoError(t, err)

	now := time.Now()
	subject := "Test Subject"
	message := &entities.Message{
		AccountID:        account.ID,
		Subject:          &subject,
		SenderEmail:      "sender@example.com",
		ReceivedDatetime: &now,
		Importance:       entities.ImportanceNormal,
	}

	err = repo.Create(ctx, message)
	assert.NoError(t, err)
	assert.Greater(t, message.MessageID, int64(0), "Message ID should be set after create")

	// Verify in DB
	var dbMessage entities.Message
	err = db.First(&dbMessage, message.MessageID).Error
	assert.NoError(t, err)
	assert.Equal(t, account.ID, dbMessage.AccountID)
	assert.Equal(t, *message.Subject, *dbMessage.Subject)
	assert.Equal(t, message.SenderEmail, dbMessage.SenderEmail)
	assert.Equal(t, message.Importance, dbMessage.Importance)
}

func TestMessageRepository_GetByID(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewMessageRepository(db)
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

	t.Run("Existing ID", func(t *testing.T) {
		message, err := repo.GetByID(ctx, testMessage.MessageID)
		assert.NoError(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, testMessage.MessageID, message.MessageID)
		assert.Equal(t, testMessage.AccountID, message.AccountID)
		assert.Equal(t, *testMessage.Subject, *message.Subject)
		assert.Equal(t, testMessage.SenderEmail, message.SenderEmail)
	})

	t.Run("Non-existing ID", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 9999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrMessageNotFound)
	})
}

func TestMessageRepository_Update(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewMessageRepository(db)
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

	t.Run("Update Existing Message", func(t *testing.T) {
		newSubject := "Updated Subject"
		testMessage.Subject = &newSubject
		testMessage.IsRead = true
		testMessage.Importance = entities.ImportanceHigh

		err := repo.Update(ctx, testMessage)
		assert.NoError(t, err)

		// Verify in DB
		var updatedMessage entities.Message
		err = db.First(&updatedMessage, testMessage.MessageID).Error
		assert.NoError(t, err)
		assert.Equal(t, *testMessage.Subject, *updatedMessage.Subject)
		assert.Equal(t, testMessage.IsRead, updatedMessage.IsRead)
		assert.Equal(t, testMessage.Importance, updatedMessage.Importance)
	})

	t.Run("Update Non-existing Message", func(t *testing.T) {
		nonExistingMessage := &entities.Message{
			MessageID:   9999,
			AccountID:   account.ID,
			SenderEmail: "notfound@example.com",
			Importance:  entities.ImportanceLow,
		}

		err := repo.Update(ctx, nonExistingMessage)
		assert.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrMessageNotFound)
	})
}

func TestMessageRepository_Delete(t *testing.T) {
	db := utils.SetupTestDB(t)
	repo := sqlite.NewMessageRepository(db)
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
		err := repo.Delete(ctx, testMessage.MessageID)
		assert.NoError(t, err)

		// Verify it's deleted
		var found entities.Message
		err = db.First(&found, testMessage.MessageID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Delete Non-existing Message", func(t *testing.T) {
		err := repo.Delete(ctx, 9999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrMessageNotFound)
	})
}
