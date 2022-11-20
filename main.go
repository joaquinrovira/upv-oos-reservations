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
		Target: []lib.TargetValue{
			{Day: time.Monday, Time: model.TimeOfDay{Hours: 13, Minutes: 30}},
		},
	}

	// Initialize and run agent
	if agent, err := lib.New(cfg); err != nil {
		panic(err)
	} else if err := agent.Run(); err != nil {
		panic(err)
	}
}
