package main

import (
	"context"
	"fmt"

	"palm/src/config"
)

// App struct
type App struct {
	ctx context.Context
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

	config.Logger.Info().Msg("Application started successfully")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	config.Logger.Debug().Str("name", name).Msg("Greet function called")
	return fmt.Sprintf("Hello %s, relax!", name)
}
