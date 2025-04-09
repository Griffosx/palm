package services

import (
	"context"
	"errors"
	"fmt"
	"palm/src/config"
	"palm/src/entities"
	"palm/src/repositories"
)

// Custom error types
var (
	ErrAttachmentNotFound = errors.New("attachment not found")
)

type AttachmentService struct {
	repo repositories.AttachmentRepository
}

func NewAttachmentService(repo repositories.AttachmentRepository) *AttachmentService {
	config.Logger.Debug().Msg("Initializing attachment service")
	return &AttachmentService{repo: repo}
}

func (s *AttachmentService) CreateAttachment(ctx context.Context, attachment *entities.Attachment) error {
	config.Logger.Info().
		Uint("messageID", attachment.MessageID).
		Str("filename", attachment.Filename).
		Str("mimeType", attachment.MimeType).
		Uint("size", attachment.Size).
		Msg("Creating new attachment")

	result := s.repo.Create(ctx, attachment)
	if result.Error != nil {
		config.Logger.Error().
			Err(result.Error).
			Uint("messageID", attachment.MessageID).
			Str("filename", attachment.Filename).
			Msg("Failed to create attachment")
		return fmt.Errorf("failed to create attachment: %w", result.Error)
	}

	config.Logger.Info().
		Uint("attachmentID", attachment.ID).
		Uint("messageID", attachment.MessageID).
		Msg("Attachment created successfully")

	return nil
}

func (s *AttachmentService) GetAttachment(ctx context.Context, id uint) (*entities.Attachment, error) {
	config.Logger.Debug().Uint("id", id).Msg("Getting attachment by ID")

	attachment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("id", id).
			Msg("Failed to get attachment")
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}

	config.Logger.Debug().
		Uint("attachmentID", id).
		Uint("messageID", attachment.MessageID).
		Msg("Attachment retrieved successfully")

	return attachment, nil
}

func (s *AttachmentService) GetAttachmentsByMessageID(ctx context.Context, messageID uint) ([]*entities.Attachment, error) {
	config.Logger.Debug().Uint("messageID", messageID).Msg("Getting attachments by message ID")

	attachments, err := s.repo.GetByMessageID(ctx, messageID)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", messageID).
			Msg("Failed to get attachments by message ID")
		return nil, fmt.Errorf("failed to get attachments: %w", err)
	}

	config.Logger.Debug().
		Uint("messageID", messageID).
		Uint("count", uint(len(attachments))).
		Msg("Attachments retrieved successfully")

	return attachments, nil
}

func (s *AttachmentService) DeleteAttachmentsByMessageID(ctx context.Context, messageID uint) error {
	config.Logger.Info().Uint("messageID", messageID).Msg("Deleting all attachments for message")

	if err := s.repo.DeleteByMessageID(ctx, messageID); err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", messageID).
			Msg("Failed to delete attachments for message")
		return fmt.Errorf("failed to delete attachments: %w", err)
	}

	config.Logger.Info().Uint("messageID", messageID).Msg("Attachments deleted successfully")
	return nil
}
