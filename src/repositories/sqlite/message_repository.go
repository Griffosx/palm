package sqlite

import (
	"context"
	"errors"
	"fmt"
	"palm/src/config"
	"palm/src/entities"
	"palm/src/repositories"

	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repositories.MessageRepository {
	config.Logger.Debug().Msg("Initializing message repository")
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *entities.Message) error {
	fmt.Printf("\n\nCreating message: %+v\n\n", message)
	config.Logger.Debug().
		Uint("accountID", message.AccountID).
		Str("senderEmail", message.SenderEmail).
		Msg("Creating new message")

	err := r.db.WithContext(ctx).Create(message).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("accountID", message.AccountID).
			Str("senderEmail", message.SenderEmail).
			Msg("Failed to create message")
	} else {
		config.Logger.Info().
			Uint("messageID", message.ID).
			Uint("accountID", message.AccountID).
			Msg("Message created successfully")
	}
	return err
}

func (r *messageRepository) GetByID(ctx context.Context, id uint) (*entities.Message, error) {
	config.Logger.Debug().Uint("id", id).Msg("Getting message by ID")

	var message entities.Message
	err := r.db.WithContext(ctx).First(&message, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config.Logger.Warn().Uint("id", id).Msg("Message not found")
			return nil, repositories.ErrMessageNotFound
		}
		config.Logger.Error().Err(err).Uint("id", id).Msg("Error retrieving message")
		return nil, err
	}
	config.Logger.Debug().
		Uint("id", id).
		Uint("accountID", message.AccountID).
		Msg("Message retrieved successfully")

	return &message, nil
}

func (r *messageRepository) Update(ctx context.Context, message *entities.Message) error {
	config.Logger.Debug().
		Uint("messageID", message.ID).
		Uint("accountID", message.AccountID).
		Msg("Updating message")

	// First check if the message exists
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Message{}).Where("message_id = ?", message.ID).Count(&count).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", message.ID).
			Msg("Error checking if message exists")
		return err
	}

	if count == 0 {
		config.Logger.Warn().
			Uint("messageID", message.ID).
			Msg("Message not found for update")
		return repositories.ErrMessageNotFound
	}

	// Then update the message
	result := r.db.WithContext(ctx).Save(message)
	if result.Error != nil {
		config.Logger.Error().
			Err(result.Error).
			Uint("messageID", message.ID).
			Msg("Error updating message")
		return result.Error
	}

	config.Logger.Info().
		Uint("messageID", message.ID).
		Msg("Message updated successfully")
	return nil
}

func (r *messageRepository) Delete(ctx context.Context, id uint) error {
	config.Logger.Debug().Uint("id", id).Msg("Deleting message")

	result := r.db.WithContext(ctx).Delete(&entities.Message{}, id)
	if result.Error != nil {
		config.Logger.Error().Err(result.Error).Uint("id", id).Msg("Error deleting message")
		return result.Error
	}
	if result.RowsAffected == 0 {
		config.Logger.Warn().Uint("id", id).Msg("Message not found for deletion")
		return repositories.ErrMessageNotFound
	}
	config.Logger.Info().Uint("id", id).Msg("Message deleted successfully")
	return nil
}
