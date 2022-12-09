package agent

import (
	"fmt"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib/logging"
	"github.com/joaquinrovira/upv-oos-reservations/lib/util"
	"github.com/joaquinrovira/upv-oos-reservations/lib/vars"
	"github.com/reugn/go-quartz/quartz"
)

func (a *Agent) Run() (err error) {
	a.Login()

	targets := a.target.Clone()

	logging.LogFile().Debug().Msg("Running...")

	value, err := a.GetReservationsData()
	if err != nil {
		return err
	}

	// Remove targets that have already been fulfilled
	reservations := value.GetReservations()
	for day := range reservations {
		logging.LogFile().Info().Msgf("%v already fulfilled in time slot %v, skipping", day.String(), reservations[day].At)
		delete(targets, day)
	}

	// Remove targets for today or the days before
	today := time.Now().Weekday()
	if today != time.Saturday {
		logging.LogFile().Info().Msgf("today is %v, avoiding already gone weekdays (including today)", today.String())
		for day := range targets {
			if day <= today {
				logging.LogFile().Info().Msgf("skipping %v", day.String())
				delete(targets, day)
			}
		}
	}

	for day, target := range targets {
		target_err := a.handleTargetList(value, day, target)
		if target_err != nil {
			logging.LogFile().Err(target_err).Msgf("unable to fulfill request for %v %v", day.String(), target)
		} else {
			logging.LogFile().Info().Msgf("request %v %v fullfilled successfully", day.String(), target)
		}
	}

	return err
}
func (a *Agent) RunWithScheduler() (err error) {
	sched := quartz.NewStdScheduler()
	sched.Start()

	runJob := quartz.NewFunctionJob(func() (int, error) { return 0, a.Run() })

	triggers := []*util.CronTrigger{util.CronSaturdayAt10, util.CronEvery15Minutes}

	if vars.Get(vars.Debug) != "" {
		trigger := util.CronEvery15Seconds
		logging.Out().Debug().Msgf("debug enabled, including high-frequency schedule (%v)", trigger.Expression())
		triggers = append(triggers, trigger)
	}

	for _, trigger := range triggers {
		sched.ScheduleJob(runJob, trigger)
		logging.Out().Info().Msgf("registered job trigger %-20v i.e., %v", fmt.Sprintf("(%v)", trigger.Expression()), trigger.Description())
	}

	<-a.ctx.Done()
	sched.Stop()

	return
}
