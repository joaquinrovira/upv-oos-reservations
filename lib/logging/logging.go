package logging

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logFile zerolog.Logger
var out zerolog.Logger
var logfile string = "output.log"

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix // UNIX Time is faster and smaller than most timestamps
	setupOut()
	setupLogFile()
}

func setupLogFile() {
	if os.Getenv("NO_LOGFILE") != "" {
		return
	}

	logfilePath, _ := filepath.Abs(logfile)
	Out().Info().Msgf("logging activity to %s", logfilePath)

	// Setup rotating logfile
	logfile := &lumberjack.Logger{
		Filename:   "output.log",
		MaxSize:    100, // 100 MB
		MaxBackups: 3,   // 3 backup logs
		MaxAge:     92,  // 3 months
	}

	logFile = zerolog.New(os.Stderr).With().Timestamp().Logger()
	logFile = logFile.Output(zerolog.ConsoleWriter{Out: logfile, TimeFormat: time.RFC3339, NoColor: true})
}

func setupOut() {
	// initialize internal logger
	out = zerolog.New(os.Stderr).With().Timestamp().Logger()

	// log a human-friendly, colorized output
	out = out.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}

func LogFile() *zerolog.Logger {

	if os.Getenv(string("NO_LOGFILE")) != "" {
		return &out
	}

	return &logFile
}
func Out() *zerolog.Logger {
	return &out
}
