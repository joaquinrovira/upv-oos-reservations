package main

import (
	"github.com/joaquinrovira/upv-oos-reservations/internal/model/config"
	"github.com/joaquinrovira/upv-oos-reservations/internal/vars"
	"github.com/joaquinrovira/upv-oos-reservations/lib"
)

func main() {
	// Build config
	user := vars.Get(vars.User)
	pass := vars.Get(vars.Pass)

	target, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	cfg := lib.Config{User: user, Pass: pass, Target: target}

	// Initialize and run agent
	if agent, err := lib.New(cfg); err != nil {
		panic(err)
	} else if err := agent.Run(); err != nil {
		panic(err)
	}
}
