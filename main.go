package main

import (
	"sync"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/config"
	"github.com/joaquinrovira/upv-oos-reservations/lib/vars"
	"github.com/reugn/go-quartz/quartz"
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
	} else {
		sched := quartz.NewStdScheduler()
		wg := sync.WaitGroup{}
		wg.Add(2)

		sched.Start()
		functionJob := quartz.NewFunctionJob(func() (int, error) {
			defer wg.Done()
			return 0, agent.Run()
		})
		sched.ScheduleJob(functionJob, quartz.NewSimpleTrigger(time.Second*5))
		wg.Wait()
		sched.Stop()
	}
	//  else if err := agent.Run(); err != nil {
	// 	panic(err)
	// }
}
