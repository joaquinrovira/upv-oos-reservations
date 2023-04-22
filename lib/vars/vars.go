package vars

import (
	"fmt"
	"os"

	"github.com/joaquinrovira/upv-oos-reservations/lib/exitcodes"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/login"

	"github.com/joho/godotenv"
)

type env string

const (
	User         env = "UPV_USER"
	Pass         env = "UPV_PASS"
	LoginType    env = "UPV_LOGIN_TYPE"
	ActivityType env = "UPV_ACTIVITY_TYPE"
	ActivityCode env = "UPV_ACTIVITY_CODE"
	DefaultCron  env = "DEFAULT_CRON"
	CustomCron   env = "CUSTOM_CRON"
	ConfigFile   env = "CONFIG_FILE"
	Debug        env = "DEBUG"
	TZ           env = "TZ" // https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List
)

var required = []env{
	User,
	Pass,
	ActivityType,
	ActivityCode,
}

var defaults = map[env]string{
	LoginType:   string(login.StudentLogin),
	ConfigFile:  "./config.json",
	DefaultCron: fmt.Sprintf("%s,%s", cronSaturdayAt10, cronEvery15Minutes),
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
			fmt.Fprintf(os.Stderr, "missing environment variable '%s'", v)
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

var cronEveryHour string = "0 0 * * * *"
var cronEvery10Seconds string = "*/10 * * * * *"
var cronEvery15Minutes string = "0 */15 * * * *"
var cronSaturdayAt10 string = "*/10 0-15 10 * * SAT"
