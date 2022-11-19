package internal

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// default level is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// initialize internal logger
	log = zerolog.New(os.Stderr).With().Timestamp().Logger()

	// log a human-friendly, colorized output
	log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

func Log() *zerolog.Logger {
	return &log
}
