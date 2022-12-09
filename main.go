package main

import (
	"github.com/joaquinrovira/upv-oos-reservations/lib/agent"
	"github.com/joaquinrovira/upv-oos-reservations/lib/vars"
)

func main() {
	// Build config
	user := vars.Get(vars.User)
	pass := vars.Get(vars.Pass)
	configFile := vars.Get(vars.ConfigFile)

	cfg := agent.Config{User: user, Pass: pass, ConfigFile: configFile}

	// Initialize and run agent
	if agent, err := agent.New(cfg); err != nil {
		panic(err)
	} else if err = agent.Run(); err != nil { // Run once on startup
		panic(err)
	} else if err := agent.RunWithScheduler(); err != nil { // Run on a schedule
		panic(err)
	}

}
