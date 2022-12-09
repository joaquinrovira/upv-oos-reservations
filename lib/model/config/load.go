package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib/model/timerange"
	"github.com/joaquinrovira/upv-oos-reservations/lib/util"
)

/*
For each weekday, we can define a list of time ranges we want to try
to make a reservation in. The list is sorted by descreasing priority.
If a slot is open within the first time range, a reservation will be
executed for that slot.
*/
type TargetValue map[time.Weekday][]timerange.TimeRange

// Returns error for invalid configurations
func CheckConfig(target TargetValue) (err error) {
	for _, trList := range target {
		for _, tr := range trList {
			err = tr.Validate()
			if err != nil {
				return
			}
		}
	}

	return
}

func LoadConfig(filename string) (target TargetValue, err error) {
	config, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	value := new(TargetValueJSON)
	err = json.Unmarshal(config, value)
	if err != nil {
		return
	}

	target, err = value.Parse()
	if err != nil {
		return
	}

	err = CheckConfig(target)
	if err != nil {
		return
	}

	return
}

type TargetValueJSON map[string][]timerange.TimeRange

func (json TargetValueJSON) Parse() (target TargetValue, err error) {
	target = make(TargetValue)

	for k, v := range json {
		weedkay, ok := util.DaysOfWeek[k]
		if !ok {
			err = fmt.Errorf("unkown weekday '%s'", k)
			return
		}
		if len(v) != 0 {
			target[weedkay] = v
		}
	}
	return
}
