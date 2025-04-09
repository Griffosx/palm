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
		Uint("messageID", recipient.MessageID).
		Str("email", recipient.Email).
		Str("type", string(recipient.RecipientType)).
		Msg("Creating new recipient")

	err := r.db.WithContext(ctx).Create(recipient).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", recipient.MessageID).
			Str("email", recipient.Email).
			Msg("Failed to create recipient")
	} else {
		config.Logger.Info().
			Uint("recipientID", recipient.ID).
			Uint("messageID", recipient.MessageID).
			Msg("Recipient created successfully")
	}
	return err
}

func (r *recipientRepository) GetByMessageID(ctx context.Context, messageID uint) ([]*entities.Recipient, error) {
	config.Logger.Debug().Uint("messageID", messageID).Msg("Getting recipients by message ID")

	var recipients []*entities.Recipient
	err := r.db.WithContext(ctx).Where("message_id = ?", messageID).Find(&recipients).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", messageID).
			Msg("Error retrieving recipients")
		return nil, err
	}

	config.Logger.Debug().
		Uint("messageID", messageID).
		Int("count", len(recipients)).
		Msg("Recipients retrieved successfully")

	return recipients, nil
}

func (r *recipientRepository) DeleteByMessageID(ctx context.Context, messageID uint) error {
	config.Logger.Debug().Uint("messageID", messageID).Msg("Deleting recipients by message ID")

	result := r.db.WithContext(ctx).Where("message_id = ?", messageID).Delete(&entities.Recipient{})
	if result.Error != nil {
		config.Logger.Error().
			Err(result.Error).
			Uint("messageID", messageID).
			Msg("Error deleting recipients")
		return result.Error
	}

	config.Logger.Info().
		Uint("messageID", messageID).
		Uint("count", uint(result.RowsAffected)).
		Msg("Recipients deleted successfully")

	return nil
}
