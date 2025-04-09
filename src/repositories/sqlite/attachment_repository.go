package sqlite

import (
	"context"
	"errors"
	"palm/src/config"
	"palm/src/entities"
	"palm/src/repositories"

	"gorm.io/gorm"
)

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) repositories.AttachmentRepository {
	config.Logger.Debug().Msg("Initializing attachment repository")
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(ctx context.Context, attachment *entities.Attachment) error {
	config.Logger.Debug().
		Uint("messageID", attachment.MessageID).
		Str("filename", attachment.Filename).
		Str("mimeType", attachment.MimeType).
		Uint("size", attachment.Size).
		Msg("Creating new attachment")

	err := r.db.WithContext(ctx).Create(attachment).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", attachment.MessageID).
			Str("filename", attachment.Filename).
			Msg("Failed to create attachment")
	} else {
		config.Logger.Info().
			Uint("messageID", attachment.MessageID).
			Msg("Attachment created successfully")
	}
	return err
}

func (r *attachmentRepository) GetByID(ctx context.Context, id uint) (*entities.Attachment, error) {
	config.Logger.Debug().Uint("id", id).Msg("Getting attachment by ID")

	var attachment entities.Attachment
	err := r.db.WithContext(ctx).First(&attachment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config.Logger.Warn().Uint("id", id).Msg("Attachment not found")
			return nil, repositories.ErrAttachmentNotFound
		}
		config.Logger.Error().Err(err).Uint("id", id).Msg("Error retrieving attachment")
		return nil, err
	}
	config.Logger.Debug().
		Uint("id", id).
		Uint("messageID", attachment.MessageID).
		Msg("Attachment retrieved successfully")

	return &attachment, nil
}

func (r *attachmentRepository) GetByMessageID(ctx context.Context, messageID uint) ([]*entities.Attachment, error) {
	config.Logger.Debug().Uint("messageID", messageID).Msg("Getting attachments by message ID")

	var attachments []*entities.Attachment
	err := r.db.WithContext(ctx).Where("message_id = ?", messageID).Find(&attachments).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", messageID).
			Msg("Error retrieving attachments")
		return nil, err
	}

	config.Logger.Debug().
		Uint("messageID", messageID).
		Int("count", len(attachments)).
		Msg("Attachments retrieved successfully")

	return attachments, nil
}

func (r *attachmentRepository) DeleteByMessageID(ctx context.Context, messageID uint) error {
	config.Logger.Debug().Uint("messageID", messageID).Msg("Deleting attachments by message ID")

	result := r.db.WithContext(ctx).Where("message_id = ?", messageID).Delete(&entities.Attachment{})
	if result.Error != nil {
		config.Logger.Error().
			Err(result.Error).
			Uint("messageID", messageID).
			Msg("Error deleting attachments")
		return result.Error
	}

	config.Logger.Info().
		Uint("messageID", messageID).
		Uint("count", uint(result.RowsAffected)).
		Msg("Attachments deleted successfully")

	return nil
}
