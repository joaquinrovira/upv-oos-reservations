package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/internal/model/timerange"
	"github.com/joaquinrovira/upv-oos-reservations/internal/vars"
	"github.com/joaquinrovira/upv-oos-reservations/lib"
)

func LoadConfig() (target lib.TargetValue, err error) {
	config, err := os.ReadFile(vars.Get(vars.ConfigFile))
	if err != nil {
		return
	}

	value := new(TargetValueJSON)
	err = json.Unmarshal(config, value)
	if err != nil {
		return
	}

	target, err = value.Parse()

	return
}

type TargetValueJSON map[string][]timerange.TimeRange

func (json TargetValueJSON) Parse() (target lib.TargetValue, err error) {
	target = make(lib.TargetValue)

	for k, v := range json {
		weedkay, ok := daysOfWeek[k]
		if !ok {
			err = fmt.Errorf("unkown weekday '%s'", k)
			return
		}
		target[weedkay] = v
	}
	return
}

var daysOfWeek = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}
