package vars

import (
	"fmt"
	"os"

	"github.com/joaquinrovira/upv-oos-reservations/internal"
	"github.com/joaquinrovira/upv-oos-reservations/internal/exitcodes"

	"github.com/joho/godotenv"
)

type variable string

const (
	User         variable = "UPV_USER"
	Pass         variable = "UPV_PASS"
	ActivityType variable = "UPV_ACTIVITY_TYPE"
	ActivityCode variable = "UPV_ACTIVITY_CODE"
	ConfigFile   variable = "CONFIG_FILE"
)

var requiredVariables = []variable{
	User,
	Pass,
}

var variableDefaults = map[variable]string{
	ActivityType: "6607",
	ActivityCode: "20705",
	ConfigFile:   "./config.json",
}

func init() {
	// load from .env (if present)
	godotenv.Load()
	checkRequired()
	applyDefaults()
}

func checkRequired() {
	missingVars := false
	for _, v := range requiredVariables {
		if os.Getenv(string(v)) == "" {
			internal.Log().Error().Msg(fmt.Sprintf("missing environment variable '%s'", v))
			missingVars = true
		}
	}
	if missingVars {
		os.Exit(exitcodes.MISSING_REQUIRED_ENV)
	}
}

func applyDefaults() {
	for k, v := range variableDefaults {
		if os.Getenv(string(k)) == "" {
			os.Setenv(string(k), v)
		}
	}
}

func Get(v variable) string {
	return os.Getenv(string(v))
}
