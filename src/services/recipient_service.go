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
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrInvalidRecipientType = errors.New("invalid recipient type")
)

type RecipientService struct {
	repo repositories.RecipientRepository
}

func NewRecipientService(repo repositories.RecipientRepository) *RecipientService {
	config.Logger.Debug().Msg("Initializing recipient service")
	return &RecipientService{repo: repo}
}

// validateRecipientType validates the recipient type is one of the allowed values
func (s *RecipientService) validateRecipientType(recipientType entities.RecipientType) error {
	config.Logger.Debug().Str("recipientType", string(recipientType)).Msg("Validating recipient type")

	if recipientType != entities.RecipientTypeTo &&
		recipientType != entities.RecipientTypeCc &&
		recipientType != entities.RecipientTypeBcc {
		config.Logger.Warn().
			Str("recipientType", string(recipientType)).
			Msg("Invalid recipient type")
		return ErrInvalidRecipientType
	}
	return nil
}

func (s *RecipientService) CreateRecipient(ctx context.Context, recipient *entities.Recipient) error {
	config.Logger.Info().
		Uint("messageID", recipient.MessageID).
		Str("email", recipient.Email).
		Str("type", string(recipient.RecipientType)).
		Msg("Creating new recipient")

	if err := s.validateRecipientType(recipient.RecipientType); err != nil {
		return err
	}

	result := s.repo.Create(ctx, recipient)
	if result.Error != nil {
		config.Logger.Error().
			Err(result.Error).
			Uint("messageID", recipient.MessageID).
			Str("email", recipient.Email).
			Msg("Failed to create recipient")
		return fmt.Errorf("failed to create recipient: %w", result.Error)
	}

	config.Logger.Info().
		Uint("recipientID", recipient.ID).
		Uint("messageID", recipient.MessageID).
		Msg("Recipient created successfully")

	return nil
}

func (s *RecipientService) GetRecipientsByMessageID(ctx context.Context, messageID uint) ([]*entities.Recipient, error) {
	config.Logger.Debug().Uint("messageID", messageID).Msg("Getting recipients by message ID")

	recipients, err := s.repo.GetByMessageID(ctx, messageID)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", messageID).
			Msg("Failed to get recipients by message ID")
		return nil, fmt.Errorf("failed to get recipients: %w", err)
	}

	config.Logger.Debug().
		Uint("messageID", messageID).
		Int("count", len(recipients)).
		Msg("Recipients retrieved successfully")

	return recipients, nil
}

func (s *RecipientService) DeleteRecipientsByMessageID(ctx context.Context, messageID uint) error {
	config.Logger.Info().Uint("messageID", messageID).Msg("Deleting all recipients for message")

	if err := s.repo.DeleteByMessageID(ctx, messageID); err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", messageID).
			Msg("Failed to delete recipients for message")
		return fmt.Errorf("failed to delete recipients: %w", err)
	}

	config.Logger.Info().Uint("messageID", messageID).Msg("Recipients deleted successfully")
	return nil
}
