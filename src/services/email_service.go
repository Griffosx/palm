package services

import (
	"context"
	"errors"
	"fmt"
	"palm/src/config"
	"palm/src/entities"
	"palm/src/repositories"

	"gorm.io/gorm"
)

// Custom error types
var (
	ErrEmailCreationFailed = errors.New("failed to create email")
	ErrEmailNotFound       = errors.New("email not found")
	ErrEmailDeleteFailed   = errors.New("failed to delete email")
	ErrInvalidPageSize     = errors.New("page size must be between 1 and 100")
)

// EmailDTO represents an email to be created with all its components
type EmailDTO struct {
	Message     *entities.Message      // The message entity
	Recipients  []*entities.Recipient  // List of recipients
	Attachments []*entities.Attachment // List of attachments (optional)
}

// PaginatedEmailsResult represents the result of a paginated email list operation
type PaginatedEmailsResult struct {
	Emails     []*EmailDTO // List of emails
	TotalCount int64       // Total number of emails matching the criteria
	Page       int         // Current page
	PageSize   int         // Page size
	TotalPages int         // Total number of pages
}

// EmailService handles operations for creating emails with related entities
type EmailService struct {
	db             *gorm.DB
	messageRepo    repositories.MessageRepository
	recipientRepo  repositories.RecipientRepository
	attachmentRepo repositories.AttachmentRepository
}

// NewEmailService creates a new EmailService
func NewEmailService(
	db *gorm.DB,
	messageRepo repositories.MessageRepository,
	recipientRepo repositories.RecipientRepository,
	attachmentRepo repositories.AttachmentRepository,
) *EmailService {
	config.Logger.Debug().Msg("Initializing email service")
	return &EmailService{
		db:             db,
		messageRepo:    messageRepo,
		recipientRepo:  recipientRepo,
		attachmentRepo: attachmentRepo,
	}
}

// validateEmail validates the email data before creation
func (s *EmailService) validateEmail(email *EmailDTO) error {
	// Message is required
	if email.Message == nil {
		return errors.New("message cannot be null")
	}

	// At least one recipient is required
	if len(email.Recipients) == 0 {
		return errors.New("at least one recipient is required")
	}

	// Validate message importance if set
	if email.Message.Importance != "" {
		messageService := NewMessageService(s.messageRepo)
		if err := messageService.validateImportance(email.Message.Importance); err != nil {
			return err
		}
	} else {
		// Set default importance if not provided
		email.Message.Importance = entities.ImportanceNormal
	}

	return nil
}

// Create creates a new email with its related entities in a single transaction
func (s *EmailService) Create(ctx context.Context, email *EmailDTO) error {
	if err := s.validateEmail(email); err != nil {
		accountID := uint(0)
		if email.Message != nil {
			accountID = email.Message.AccountID
		}
		config.Logger.Error().
			Err(err).
			Uint("accountID", accountID).
			Msg("Email validation failed")
		return err
	}

	config.Logger.Info().
		Uint("accountID", email.Message.AccountID).
		Str("subject", stringOrEmpty(email.Message.Subject)).
		Str("senderEmail", email.Message.SenderEmail).
		Int("recipientCount", len(email.Recipients)).
		Int("attachmentCount", len(email.Attachments)).
		Msg("Creating new email")

	// Start a transaction
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the message first to get its ID
		if err := tx.Create(email.Message).Error; err != nil {
			config.Logger.Error().
				Err(err).
				Uint("accountID", email.Message.AccountID).
				Msg("Failed to create message")
			return err
		}

		// Set the message ID on all recipients and create them
		for _, recipient := range email.Recipients {
			recipient.MessageID = email.Message.ID
			if err := tx.Create(recipient).Error; err != nil {
				config.Logger.Error().
					Err(err).
					Uint("messageID", email.Message.ID).
					Str("email", recipient.Email).
					Msg("Failed to create recipient")
				return err
			}
		}

		// Create attachments if any
		for _, attachment := range email.Attachments {
			attachment.MessageID = email.Message.ID
			if err := tx.Create(attachment).Error; err != nil {
				config.Logger.Error().
					Err(err).
					Uint("messageID", email.Message.ID).
					Str("filename", attachment.Filename).
					Msg("Failed to create attachment")
				return err
			}
		}

		return nil
	})

	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("accountID", email.Message.AccountID).
			Msg("Email creation failed")
		return fmt.Errorf("%w: %s", ErrEmailCreationFailed, err.Error())
	}

	config.Logger.Info().
		Uint("messageID", email.Message.ID).
		Uint("accountID", email.Message.AccountID).
		Int("recipientCount", len(email.Recipients)).
		Int("attachmentCount", len(email.Attachments)).
		Msg("Email created successfully")

	return nil
}

// GetByID retrieves an email with all its components by message ID
func (s *EmailService) GetByID(ctx context.Context, messageID uint) (*EmailDTO, error) {
	config.Logger.Debug().Uint("messageID", messageID).Msg("Getting email by ID")

	// Get the message
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", messageID).
			Msg("Failed to get message")
		if errors.Is(err, repositories.ErrMessageNotFound) {
			return nil, ErrEmailNotFound
		}
		return nil, err
	}

	// Get recipients for this message
	recipients, err := s.recipientRepo.GetByMessageID(ctx, messageID)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", messageID).
			Msg("Failed to get recipients")
		return nil, err
	}

	// Get attachments for this message
	attachments, err := s.attachmentRepo.GetByMessageID(ctx, messageID)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", messageID).
			Msg("Failed to get attachments")
		return nil, err
	}

	email := &EmailDTO{
		Message:     message,
		Recipients:  recipients,
		Attachments: attachments,
	}

	config.Logger.Debug().
		Uint("messageID", messageID).
		Int("recipientCount", len(recipients)).
		Int("attachmentCount", len(attachments)).
		Msg("Email retrieved successfully")

	return email, nil
}

// ListCount returns the total number of emails for a specific account
func (s *EmailService) ListCount(ctx context.Context, accountID uint) (int64, error) {
	config.Logger.Debug().
		Uint("accountID", accountID).
		Msg("Counting emails for account")

	var totalCount int64

	// Get total count
	err := s.db.WithContext(ctx).
		Model(&entities.Message{}).
		Where("account_id = ?", accountID).
		Count(&totalCount).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("accountID", accountID).
			Msg("Failed to count messages")
		return 0, err
	}

	config.Logger.Debug().
		Uint("accountID", accountID).
		Uint("count", uint(totalCount)).
		Msg("Email count retrieved successfully")

	return totalCount, nil
}

// List retrieves a paginated list of emails for a specific account
func (s *EmailService) List(ctx context.Context, accountID uint, pageSize int, page int) (*PaginatedEmailsResult, error) {
	config.Logger.Debug().
		Uint("accountID", accountID).
		Int("pageSize", pageSize).
		Int("page", page).
		Msg("Listing emails for account")

	// Validate page size
	if pageSize < 1 || pageSize > 100 {
		config.Logger.Error().
			Int("pageSize", pageSize).
			Msg("Invalid page size")
		return nil, ErrInvalidPageSize
	}

	// Default page to 1 if not positive
	if page < 1 {
		page = 1
	}

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	var messages []*entities.Message
	var totalCount int64

	// Get total count first
	err := s.db.WithContext(ctx).
		Model(&entities.Message{}).
		Where("account_id = ?", accountID).
		Count(&totalCount).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("accountID", accountID).
			Msg("Failed to count messages")
		return nil, err
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	// Get messages with pagination
	err = s.db.WithContext(ctx).
		Where("account_id = ?", accountID).
		Order("received_datetime DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&messages).Error
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("accountID", accountID).
			Msg("Failed to list messages")
		return nil, err
	}

	// Create the result with pagination info
	result := &PaginatedEmailsResult{
		Emails:     make([]*EmailDTO, 0, len(messages)),
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	// For each message, fetch related entities and create an EmailDTO
	for _, message := range messages {
		// Get recipients for this message
		recipients, err := s.recipientRepo.GetByMessageID(ctx, message.ID)
		if err != nil {
			config.Logger.Error().
				Err(err).
				Uint("messageID", message.ID).
				Msg("Failed to get recipients")
			continue
		}

		// Get attachments for this message
		attachments, err := s.attachmentRepo.GetByMessageID(ctx, message.ID)
		if err != nil {
			config.Logger.Error().
				Err(err).
				Uint("messageID", message.ID).
				Msg("Failed to get attachments")
			continue
		}

		// Create an EmailDTO and add it to the result
		email := &EmailDTO{
			Message:     message,
			Recipients:  recipients,
			Attachments: attachments,
		}
		result.Emails = append(result.Emails, email)
	}

	config.Logger.Debug().
		Uint("accountID", uint(accountID)).
		Int("found", len(result.Emails)).
		Uint("total", uint(totalCount)).
		Int("page", page).
		Int("totalPages", totalPages).
		Msg("Emails listed successfully")

	return result, nil
}

// Delete deletes an email with all its components in a single transaction
func (s *EmailService) Delete(ctx context.Context, messageID int64) error {
	config.Logger.Info().Int64("messageID", messageID).Msg("Deleting email")

	// Start a transaction
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete attachments first (foreign key references)
		if err := tx.Where("message_id = ?", messageID).Delete(&entities.Attachment{}).Error; err != nil {
			config.Logger.Error().
				Err(err).
				Int64("messageID", messageID).
				Msg("Failed to delete attachments")
			return err
		}

		// Delete recipients
		if err := tx.Where("message_id = ?", messageID).Delete(&entities.Recipient{}).Error; err != nil {
			config.Logger.Error().
				Err(err).
				Int64("messageID", messageID).
				Msg("Failed to delete recipients")
			return err
		}

		// Delete the message last
		result := tx.Delete(&entities.Message{}, messageID)
		if result.Error != nil {
			config.Logger.Error().
				Err(result.Error).
				Int64("messageID", messageID).
				Msg("Failed to delete message")
			return result.Error
		}

		if result.RowsAffected == 0 {
			config.Logger.Warn().
				Int64("messageID", messageID).
				Msg("Message not found for deletion")
			return repositories.ErrMessageNotFound
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, repositories.ErrMessageNotFound) {
			return ErrEmailNotFound
		}
		config.Logger.Error().
			Err(err).
			Int64("messageID", messageID).
			Msg("Email deletion failed")
		return fmt.Errorf("%w: %s", ErrEmailDeleteFailed, err.Error())
	}

	config.Logger.Info().
		Int64("messageID", messageID).
		Msg("Email deleted successfully")

	return nil
}

// Helper function to safely return string value from pointer
func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
