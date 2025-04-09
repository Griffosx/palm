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
		Int64("messageID", attachment.MessageID).
		Str("filename", attachment.Filename).
		Str("mimeType", attachment.MimeType).
		Int64("size", attachment.Size).
		Msg("Creating new attachment")

	err := r.db.WithContext(ctx).Create(attachment).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Int64("messageID", attachment.MessageID).
			Str("filename", attachment.Filename).
			Msg("Failed to create attachment")
	} else {
		config.Logger.Info().
			Int64("attachmentID", attachment.AttachmentID).
			Int64("messageID", attachment.MessageID).
			Msg("Attachment created successfully")
	}
	return err
}

func (r *attachmentRepository) GetByID(ctx context.Context, id int64) (*entities.Attachment, error) {
	config.Logger.Debug().Int64("id", id).Msg("Getting attachment by ID")

	var attachment entities.Attachment
	err := r.db.WithContext(ctx).First(&attachment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config.Logger.Warn().Int64("id", id).Msg("Attachment not found")
			return nil, repositories.ErrAttachmentNotFound
		}
		config.Logger.Error().Err(err).Int64("id", id).Msg("Error retrieving attachment")
		return nil, err
	}
	config.Logger.Debug().
		Int64("id", id).
		Int64("messageID", attachment.MessageID).
		Msg("Attachment retrieved successfully")

	return &attachment, nil
}

func (r *attachmentRepository) GetByMessageID(ctx context.Context, messageID int64) ([]*entities.Attachment, error) {
	config.Logger.Debug().Int64("messageID", messageID).Msg("Getting attachments by message ID")

	var attachments []*entities.Attachment
	err := r.db.WithContext(ctx).Where("message_id = ?", messageID).Find(&attachments).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Int64("messageID", messageID).
			Msg("Error retrieving attachments")
		return nil, err
	}

	config.Logger.Debug().
		Int64("messageID", messageID).
		Int("count", len(attachments)).
		Msg("Attachments retrieved successfully")

	return attachments, nil
}

func (r *attachmentRepository) DeleteByMessageID(ctx context.Context, messageID int64) error {
	config.Logger.Debug().Int64("messageID", messageID).Msg("Deleting attachments by message ID")

	result := r.db.WithContext(ctx).Where("message_id = ?", messageID).Delete(&entities.Attachment{})
	if result.Error != nil {
		config.Logger.Error().
			Err(result.Error).
			Int64("messageID", messageID).
			Msg("Error deleting attachments")
		return result.Error
	}

	config.Logger.Info().
		Int64("messageID", messageID).
		Int64("count", result.RowsAffected).
		Msg("Attachments deleted successfully")

	return nil
}
