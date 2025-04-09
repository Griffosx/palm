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
	ErrMessageNotFound   = errors.New("message not found")
	ErrInvalidImportance = errors.New("invalid importance value")
)

type MessageService struct {
	repo repositories.MessageRepository
}

func NewMessageService(repo repositories.MessageRepository) *MessageService {
	config.Logger.Debug().Msg("Initializing message service")
	return &MessageService{repo: repo}
}

// validateImportance validates the importance value is one of the allowed values
func (s *MessageService) validateImportance(importance entities.Importance) error {
	config.Logger.Debug().Str("importance", string(importance)).Msg("Validating importance")

	if importance != entities.ImportanceLow &&
		importance != entities.ImportanceNormal &&
		importance != entities.ImportanceHigh {
		config.Logger.Warn().
			Str("importance", string(importance)).
			Msg("Invalid importance value")
		return ErrInvalidImportance
	}
	return nil
}

func (s *MessageService) CreateMessage(ctx context.Context, message *entities.Message) error {
	config.Logger.Info().
		Uint("accountID", message.AccountID).
		Str("senderEmail", message.SenderEmail).
		Msg("Creating new message")

	if err := s.validateImportance(message.Importance); err != nil {
		return err
	}

	if err := s.repo.Create(ctx, message); err != nil {
		config.Logger.Error().
			Err(err).
			Uint("accountID", message.AccountID).
			Msg("Failed to create message")
		return fmt.Errorf("failed to create message: %w", err)
	}

	config.Logger.Info().
		Uint("messageID", message.ID).
		Uint("accountID", message.AccountID).
		Msg("Message created successfully")

	return nil
}

func (s *MessageService) GetMessage(ctx context.Context, id uint) (*entities.Message, error) {
	config.Logger.Debug().Uint("id", id).Msg("Getting message by ID")

	message, err := s.repo.GetByID(ctx, id)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("id", id).
			Msg("Failed to get message")
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	config.Logger.Debug().
		Uint("messageID", id).
		Uint("accountID", message.AccountID).
		Msg("Message retrieved successfully")

	return message, nil
}

func (s *MessageService) UpdateMessage(ctx context.Context, message *entities.Message) error {
	config.Logger.Info().
		Uint("messageID", message.ID).
		Uint("accountID", message.AccountID).
		Msg("Updating message")

	if err := s.validateImportance(message.Importance); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, message); err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", message.ID).
			Msg("Failed to update message")
		return fmt.Errorf("failed to update message: %w", err)
	}

	config.Logger.Info().
		Uint("messageID", message.ID).
		Msg("Message updated successfully")

	return nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, id uint) error {
	config.Logger.Info().Uint("id", id).Msg("Deleting message")

	if err := s.repo.Delete(ctx, id); err != nil {
		config.Logger.Error().
			Err(err).
			Uint("id", id).
			Msg("Failed to delete message")
		return fmt.Errorf("failed to delete message: %w", err)
	}

	config.Logger.Info().Uint("id", id).Msg("Message deleted successfully")
	return nil
}
