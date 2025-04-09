package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"palm/src/config"
	"palm/src/entities"
	"palm/src/repositories/sqlite"
	"palm/src/services"
	"time"
)

// EmailFixture represents a single email fixture from the JSON file
type EmailFixture struct {
	AccountID   uint                `json:"account_id"`
	Message     MessageFixture      `json:"message"`
	Recipients  []RecipientFixture  `json:"recipients"`
	Attachments []AttachmentFixture `json:"attachments"`
}

type MessageFixture struct {
	Subject          string `json:"subject"`
	Body             string `json:"body"`
	BodyPreview      string `json:"body_preview"`
	SenderEmail      string `json:"sender_email"`
	SenderName       string `json:"sender_name"`
	ReceivedDatetime string `json:"received_datetime"`
	SentDatetime     string `json:"sent_datetime"`
	IsDraft          bool   `json:"is_draft"`
	IsRead           bool   `json:"is_read"`
	Importance       string `json:"importance"`
	ConversationID   string `json:"conversation_id"`
}

type RecipientFixture struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	RecipientType string `json:"recipient_type"`
}

type AttachmentFixture struct {
	Filename  string `json:"filename"`
	MimeType  string `json:"mime_type"`
	Size      uint   `json:"size"`
	LocalPath string `json:"local_path"`
}

func checkDB() error {
	config.InitLogger()
	config.Logger.Info().Msg("Checking database and ensuring tables exist")

	// Initialize database
	db, err := config.PalmDB(false)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Verify that tables exist by checking if we can query them
	tables := []string{"accounts", "messages", "recipients", "attachments"}
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			return fmt.Errorf("table '%s' does not exist after migration", table)
		}
		config.Logger.Info().Str("table", table).Msg("Table exists in database")
	}

	config.Logger.Info().Msg("Database check completed successfully")
	return nil
}

// populateEmails loads fixture data from JSON and creates emails in the database
func populateEmails() error {
	config.InitLogger()
	config.Logger.Info().Msg("Starting email population from fixtures")

	// Initialize database
	db, err := config.PalmDB(false)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repositories
	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)

	// Initialize email service
	emailService := services.NewEmailService(
		db,
		messageRepo,
		recipientRepo,
		attachmentRepo,
	)

	// Load fixture data
	fixtures, err := loadFixtures()
	if err != nil {
		return fmt.Errorf("failed to load fixtures: %w", err)
	}

	config.Logger.Info().Int("count", len(fixtures)).Msg("Loaded email fixtures")

	// Create each email fixture
	ctx := context.Background()
	for i, fixture := range fixtures {
		email, err := convertFixtureToEmailDTO(fixture)
		if err != nil {
			config.Logger.Error().
				Err(err).
				Int("index", i).
				Msg("Failed to convert fixture to email DTO")
			continue
		}

		err = emailService.Create(ctx, email)
		if err != nil {
			config.Logger.Error().
				Err(err).
				Int("index", i).
				Uint("accountID", email.Message.AccountID).
				Msg("Failed to create email from fixture")
			continue
		}

		config.Logger.Info().
			Int("index", i).
			Uint("messageID", email.Message.ID).
			Uint("accountID", email.Message.AccountID).
			Str("subject", *email.Message.Subject).
			Msg("Created email from fixture")
	}

	config.Logger.Info().Msg("Finished populating emails from fixtures")
	return nil
}

// loadFixtures loads the email fixtures from the JSON file
func loadFixtures() ([]EmailFixture, error) {
	// Open fixtures file
	file, err := os.Open("scripts/assets/fixtures.json")
	if err != nil {
		return nil, fmt.Errorf("could not open fixtures file: %w", err)
	}
	defer file.Close()

	// Decode JSON data
	var fixtures []EmailFixture
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&fixtures); err != nil {
		return nil, fmt.Errorf("could not decode fixtures JSON: %w", err)
	}

	return fixtures, nil
}

// convertFixtureToEmailDTO converts a fixture to an EmailDTO
func convertFixtureToEmailDTO(fixture EmailFixture) (*services.EmailDTO, error) {
	// Create message entity
	subject := fixture.Message.Subject
	body := fixture.Message.Body
	bodyPreview := fixture.Message.BodyPreview
	senderName := fixture.Message.SenderName
	conversationID := fixture.Message.ConversationID

	// Parse time fields
	var receivedTime, sentTime *time.Time
	if fixture.Message.ReceivedDatetime != "" {
		t, err := time.Parse(time.RFC3339, fixture.Message.ReceivedDatetime)
		if err != nil {
			return nil, fmt.Errorf("invalid received_datetime format: %w", err)
		}
		receivedTime = &t
	}

	if fixture.Message.SentDatetime != "" {
		t, err := time.Parse(time.RFC3339, fixture.Message.SentDatetime)
		if err != nil {
			return nil, fmt.Errorf("invalid sent_datetime format: %w", err)
		}
		sentTime = &t
	}

	message := &entities.Message{
		AccountID:        fixture.AccountID,
		Subject:          &subject,
		Body:             &body,
		BodyPreview:      &bodyPreview,
		SenderEmail:      fixture.Message.SenderEmail,
		SenderName:       &senderName,
		ReceivedDatetime: receivedTime,
		SentDatetime:     sentTime,
		IsDraft:          fixture.Message.IsDraft,
		IsRead:           fixture.Message.IsRead,
		Importance:       entities.Importance(fixture.Message.Importance),
		ConversationID:   &conversationID,
	}

	// Create recipients
	recipients := make([]*entities.Recipient, 0, len(fixture.Recipients))
	for _, r := range fixture.Recipients {
		name := r.Name
		recipient := &entities.Recipient{
			Email:         r.Email,
			Name:          &name,
			RecipientType: entities.RecipientType(r.RecipientType),
		}
		recipients = append(recipients, recipient)
	}

	// Create attachments
	attachments := make([]*entities.Attachment, 0, len(fixture.Attachments))
	for _, a := range fixture.Attachments {
		localPath := a.LocalPath
		attachment := &entities.Attachment{
			Filename:  a.Filename,
			MimeType:  a.MimeType,
			Size:      uint(a.Size),
			LocalPath: &localPath,
		}
		attachments = append(attachments, attachment)
	}

	// Create the EmailDTO
	email := &services.EmailDTO{
		Message:     message,
		Recipients:  recipients,
		Attachments: attachments,
	}

	return email, nil
}

func main() {
	// Check if database tables exist, this will also create them if they don't
	if err := checkDB(); err != nil {
		fmt.Printf("Error checking database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Database check completed successfully, tables exist")

	// Ask if we should continue with populating data
	fmt.Print("Do you want to populate the database with sample emails? (y/n): ")
	var answer string
	fmt.Scanln(&answer)

	if answer == "y" || answer == "Y" {
		if err := populateEmails(); err != nil {
			fmt.Printf("Error populating emails: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully populated emails!")
	} else {
		fmt.Println("Skipping database population")
	}
}
