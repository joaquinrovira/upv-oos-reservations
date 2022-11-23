package main

import (
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/internal/model/daytime"
	"github.com/joaquinrovira/upv-oos-reservations/internal/model/timerange"
	"github.com/joaquinrovira/upv-oos-reservations/internal/vars"
	"github.com/joaquinrovira/upv-oos-reservations/lib"
)

func main() {
	// Build config
	user := vars.Get(vars.User)
	pass := vars.Get(vars.Pass)
	cfg := lib.Config{
		User: user,
		Pass: pass,
		Target: lib.TargetValue{
			time.Monday: []timerange.TimeRange{
				{
					Start: daytime.DayTime{Hour: 16},
					End:   daytime.DayTime{Hour: 17},
				},
				{
					Start: daytime.DayTime{Hour: 8},
					End:   daytime.DayTime{Hour: 21},
				},
			},
		},
	}

	// Initialize and run agent
	if agent, err := lib.New(cfg); err != nil {
		panic(err)
	} else if err := agent.Run(); err != nil {
		panic(err)
	}
}
