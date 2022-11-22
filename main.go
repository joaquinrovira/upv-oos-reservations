package main

import (
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/internal/model"
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
			time.Monday: []model.TimeRange{
				{
					StartTime: model.TimeOfDay{Hour: 20},
					EndTime:   model.TimeOfDay{Hour: 22},
				},
				{
					StartTime: model.TimeOfDay{Hour: 8},
					EndTime:   model.TimeOfDay{Hour: 21},
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
