package sqlite

import (
	"context"
	"palm/src/config"
	"palm/src/entities"
	"palm/src/repositories"

	"gorm.io/gorm"
)

type recipientRepository struct {
	db *gorm.DB
}

func NewRecipientRepository(db *gorm.DB) repositories.RecipientRepository {
	config.Logger.Debug().Msg("Initializing recipient repository")
	return &recipientRepository{db: db}
}

func (r *recipientRepository) Create(ctx context.Context, recipient *entities.Recipient) error {
	config.Logger.Debug().
		Int64("messageID", recipient.MessageID).
		Str("email", recipient.Email).
		Str("type", string(recipient.RecipientType)).
		Msg("Creating new recipient")

	err := r.db.WithContext(ctx).Create(recipient).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Int64("messageID", recipient.MessageID).
			Str("email", recipient.Email).
			Msg("Failed to create recipient")
	} else {
		config.Logger.Info().
			Int64("recipientID", recipient.RecipientID).
			Int64("messageID", recipient.MessageID).
			Msg("Recipient created successfully")
	}
	return err
}

func (r *recipientRepository) GetByMessageID(ctx context.Context, messageID int64) ([]*entities.Recipient, error) {
	config.Logger.Debug().Int64("messageID", messageID).Msg("Getting recipients by message ID")

	var recipients []*entities.Recipient
	err := r.db.WithContext(ctx).Where("message_id = ?", messageID).Find(&recipients).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Int64("messageID", messageID).
			Msg("Error retrieving recipients")
		return nil, err
	}

	config.Logger.Debug().
		Int64("messageID", messageID).
		Int("count", len(recipients)).
		Msg("Recipients retrieved successfully")

	return recipients, nil
}

func (r *recipientRepository) DeleteByMessageID(ctx context.Context, messageID int64) error {
	config.Logger.Debug().Int64("messageID", messageID).Msg("Deleting recipients by message ID")

	result := r.db.WithContext(ctx).Where("message_id = ?", messageID).Delete(&entities.Recipient{})
	if result.Error != nil {
		config.Logger.Error().
			Err(result.Error).
			Int64("messageID", messageID).
			Msg("Error deleting recipients")
		return result.Error
	}

	config.Logger.Info().
		Int64("messageID", messageID).
		Int64("count", result.RowsAffected).
		Msg("Recipients deleted successfully")

	return nil
}
