package config

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// InitLogger initializes the global logger with zerolog
func InitLogger() {
	// Set global time format to ISO8601
	zerolog.TimeFieldFormat = time.RFC3339

	// Configure console writer with color and caller info
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Configure multi-writer if you need to log to file as well
	writers := []io.Writer{consoleWriter}
	multi := zerolog.MultiLevelWriter(writers...)

	// Set global logger
	Logger = zerolog.New(multi).
		With().
		Timestamp().
		Caller().
		Logger()

	// Also set the global package-level logger
	log.Logger = Logger
}
