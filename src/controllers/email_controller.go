package controllers

import (
	"context"
	"palm/src/config"
	"palm/src/services"
)

// EmailController handles HTTP requests related to emails
type EmailController struct {
	emailService *services.EmailService
}

// NewEmailController creates a new email controller
func NewEmailController(emailService *services.EmailService) *EmailController {
	config.Logger.Debug().Msg("Initializing email controller")
	return &EmailController{
		emailService: emailService,
	}
}

// ListEmailsResponse is the response for the ListEmails method
type ListEmailsResponse struct {
	Emails     []EmailResponse `json:"emails"`
	TotalCount int64           `json:"totalCount"`
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
	TotalPages int             `json:"totalPages"`
}

// EmailResponse represents the email data returned to the frontend
type EmailResponse struct {
	ID          uint                 `json:"id"`
	AccountID   uint                 `json:"accountId"`
	Subject     string               `json:"subject"`
	Body        string               `json:"body"`
	SenderName  string               `json:"senderName"`
	SenderEmail string               `json:"senderEmail"`
	ReceivedAt  string               `json:"receivedAt"`
	IsRead      bool                 `json:"isRead"`
	Importance  string               `json:"importance"`
	Recipients  []RecipientResponse  `json:"recipients"`
	Attachments []AttachmentResponse `json:"attachments,omitempty"`
}

// RecipientResponse represents a recipient in the response
type RecipientResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Type  string `json:"type"` // To, CC, BCC
}

// AttachmentResponse represents an attachment in the response
type AttachmentResponse struct {
	ID       uint   `json:"id"`
	Filename string `json:"filename"`
	Size     uint   `json:"size"`
	MimeType string `json:"mimeType"`
}

// ListEmails returns a paginated list of emails for an account
func (c *EmailController) ListEmails(ctx context.Context, accountID uint, page int, pageSize int) (*ListEmailsResponse, error) {
	config.Logger.Debug().
		Uint("accountID", accountID).
		Int("page", page).
		Int("pageSize", pageSize).
		Msg("List emails request received")

	result, err := c.emailService.List(ctx, accountID, pageSize, page)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("accountID", accountID).
			Msg("Failed to list emails")
		return nil, err
	}

	response := &ListEmailsResponse{
		Emails:     make([]EmailResponse, 0, len(result.Emails)),
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}

	// Convert service DTOs to response format
	for _, email := range result.Emails {
		emailResp := mapEmailToResponse(email)
		response.Emails = append(response.Emails, emailResp)
	}

	config.Logger.Debug().
		Uint("accountID", accountID).
		Int("emailCount", len(response.Emails)).
		Int64("totalCount", response.TotalCount).
		Msg("Emails listed successfully")

	return response, nil
}

// GetEmail returns a single email with its details
func (c *EmailController) GetEmail(ctx context.Context, messageID uint) (*EmailResponse, error) {
	config.Logger.Debug().
		Uint("messageID", messageID).
		Msg("Get email request received")

	email, err := c.emailService.GetByID(ctx, messageID)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Uint("messageID", messageID).
			Msg("Failed to get email")
		return nil, err
	}

	response := mapEmailToResponse(email)

	config.Logger.Debug().
		Uint("messageID", messageID).
		Msg("Email retrieved successfully")

	return &response, nil
}

// mapEmailToResponse converts an EmailDTO to an EmailResponse
func mapEmailToResponse(email *services.EmailDTO) EmailResponse {
	var subject, body string
	if email.Message.Subject != nil {
		subject = *email.Message.Subject
	}
	if email.Message.Body != nil {
		body = *email.Message.Body
	}

	// Format received date as ISO string or empty if nil
	receivedAt := ""
	if email.Message.ReceivedDatetime != nil {
		receivedAt = email.Message.ReceivedDatetime.Format("2006-01-02T15:04:05Z07:00")
	}

	// Map recipients
	recipients := make([]RecipientResponse, 0, len(email.Recipients))
	for _, r := range email.Recipients {
		name := ""
		if r.Name != nil {
			name = *r.Name
		}

		recipients = append(recipients, RecipientResponse{
			ID:    r.ID,
			Email: r.Email,
			Name:  name,
			Type:  string(r.RecipientType),
		})
	}

	// Map attachments
	attachments := make([]AttachmentResponse, 0, len(email.Attachments))
	for _, a := range email.Attachments {
		attachments = append(attachments, AttachmentResponse{
			ID:       a.ID,
			Filename: a.Filename,
			Size:     a.Size,
			MimeType: a.MimeType,
		})
	}

	// Handle sender name (which is a pointer)
	senderName := ""
	if email.Message.SenderName != nil {
		senderName = *email.Message.SenderName
	}

	return EmailResponse{
		ID:          email.Message.ID,
		AccountID:   email.Message.AccountID,
		Subject:     subject,
		Body:        body,
		SenderName:  senderName,
		SenderEmail: email.Message.SenderEmail,
		ReceivedAt:  receivedAt,
		IsRead:      email.Message.IsRead,
		Importance:  string(email.Message.Importance),
		Recipients:  recipients,
		Attachments: attachments,
	}
}
