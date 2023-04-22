package logging

import (
	"os"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib/vars"
	"github.com/rs/zerolog"
)

var out zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix // UNIX Time is faster and smaller than most timestamps

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if vars.Has(vars.Debug) {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	setupOut()
}

func setupOut() {
	// initialize internal logger
	out = zerolog.New(os.Stderr).With().Timestamp().Logger()

	// log a human-friendly, colorized output
	out = out.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}

func Out() *zerolog.Logger {
	return &out
}
