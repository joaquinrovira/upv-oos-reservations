package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joaquinrovira/upv-oos-reservations/lib"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/timerange"
	"github.com/joaquinrovira/upv-oos-reservations/lib/util"
	"github.com/joaquinrovira/upv-oos-reservations/lib/vars"
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
