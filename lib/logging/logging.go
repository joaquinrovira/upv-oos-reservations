package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log zerolog.Logger

func init() {
	// Setup rotating logfile
	logfile := &lumberjack.Logger{
		Filename:   "output.log",
		MaxSize:    100, // 100 MB
		MaxBackups: 3,   // 3 backup logs
		MaxAge:     92,  // 3 months
	}

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// default level is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// initialize internal logger
	log = zerolog.New(os.Stderr).With().Timestamp().Logger()

	// log a human-friendly, colorized output
	log = log.Output(zerolog.ConsoleWriter{Out: logfile, TimeFormat: time.RFC3339, NoColor: true})
}

func Log() *zerolog.Logger {
	return &log
}
