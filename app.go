package main

import (
	"context"
	"fmt"

	"palm/src/config"
	"palm/src/controllers"
	"palm/src/repositories/sqlite"
	"palm/src/services"

	"gorm.io/gorm"
)

// App struct
type App struct {
	ctx             context.Context
	db              *gorm.DB
	emailController *controllers.EmailController
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize logger
	config.InitLogger()

	// Get database connection
	db, err := config.PalmDB(false)
	if err != nil {
		config.Logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	a.db = db

	// Initialize repositories
	messageRepo := sqlite.NewMessageRepository(db)
	recipientRepo := sqlite.NewRecipientRepository(db)
	attachmentRepo := sqlite.NewAttachmentRepository(db)

	// Initialize services
	emailService := services.NewEmailService(db, messageRepo, recipientRepo, attachmentRepo)

	// Initialize controllers
	a.emailController = controllers.NewEmailController(emailService)

	config.Logger.Info().Msg("Application started successfully")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	config.Logger.Debug().Str("name", name).Msg("Greet function called")
	return fmt.Sprintf("Hello %s, relax!", name)
}

// ListEmails returns a paginated list of emails for the given account
func (a *App) ListEmails(accountID uint, page int, pageSize int) (*controllers.ListEmailsResponse, error) {
	config.Logger.Debug().
		Uint("accountID", accountID).
		Int("page", page).
		Int("pageSize", pageSize).
		Msg("ListEmails called from frontend")

	return a.emailController.ListEmails(a.ctx, accountID, page, pageSize)
}

// GetEmail returns the details of a specific email
func (a *App) GetEmail(messageID uint) (*controllers.EmailResponse, error) {
	config.Logger.Debug().
		Uint("messageID", messageID).
		Msg("GetEmail called from frontend")

	return a.emailController.GetEmail(a.ctx, messageID)
}
