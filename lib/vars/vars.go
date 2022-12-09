package vars

import (
	"fmt"
	"os"

	"github.com/joaquinrovira/upv-oos-reservations/lib/exitcodes"
	"github.com/joaquinrovira/upv-oos-reservations/lib/logging"

	"github.com/joho/godotenv"
)

type env string

const (
	User         env = "UPV_USER"
	Pass         env = "UPV_PASS"
	ActivityType env = "UPV_ACTIVITY_TYPE"
	ActivityCode env = "UPV_ACTIVITY_CODE"
	CustomCron   env = "CUSTOM_CRON"
	ConfigFile   env = "CONFIG_FILE"
	Debug        env = "DEBUG"
)

var required = []env{
	User,
	Pass,
	ActivityType,
	ActivityCode,
}

var defaults = map[env]string{
	ConfigFile: "./config.json",
}

func init() {
	// load from .env (if present)
	godotenv.Load()
	checkRequired()
	applyDefaults()
}

func checkRequired() {
	missingVars := false
	for _, v := range required {
		if os.Getenv(string(v)) == "" {
			logging.Out().Error().Msg(fmt.Sprintf("missing environment variable '%s'", v))
			missingVars = true
		}
	}
	if missingVars {
		os.Exit(exitcodes.MISSING_REQUIRED_ENV)
	}
}

func applyDefaults() {
	for k, v := range defaults {
		if os.Getenv(string(k)) == "" { // Apply defaults without overwriting existing values
			os.Setenv(string(k), v)
		}
	}
}

func Get(v env) string {
	return os.Getenv(string(v))
}

func Has(v env) bool {
	return os.Getenv(string(v)) != ""
}
