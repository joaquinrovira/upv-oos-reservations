package agent

import (
	"fmt"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib/logging"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/daytime"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/timerange"
	"github.com/joaquinrovira/upv-oos-reservations/lib/util"
	"github.com/joaquinrovira/upv-oos-reservations/lib/vars"
	"github.com/reugn/go-quartz/quartz"
)

func (a *Agent) Run() (err error) {
	a.Login()

	targets := a.target.Clone()

	logging.Out().Info().Msg("+---------------------+")
	logging.Out().Info().Msg("| Reservation trigger |")
	logging.Out().Info().Msg("+---------------------+")

	value, err := a.GetReservationsData()
	if err != nil {
		return err
	}

	// Remove targets that have already been fulfilled
	reservations := value.GetReservations()
	for day := range reservations {
		logging.Out().Info().Msgf("%-10v (skipping) already fulfilled in time slot %v", day.String(), reservations[day].At)
		delete(targets, day)
	}

	// Remove targets for today or the days before
	today := time.Now().Weekday()
	if today != time.Saturday {
		for day := range targets {
			if day <= today {
				logging.Out().Info().Msgf("%-10v (skipping) avoiding past days", day.String())
				delete(targets, day)
			}
		}
	}

	// Minimum 24hr buffer
	logging.Out().Debug().Msg("ensuring 24hr buffer")
	tomorrow := time.Now().Add(time.Hour * 24).Weekday()
	now, _ := daytime.FromTime(time.Now())
	var newTargetValue []timerange.TimeRange
	for _, rnge := range targets[tomorrow] {
		if rnge.End.Value() < now.Value() {
			logging.Out().Debug().Msgf("%-10v (avoiding) range %v with less than 24h reservation margin", tomorrow, rnge)
			continue
		}

		if rnge.Start.Value() < now.Value() {
			logging.Out().Debug().Msgf("%-10v (modifying) range %v to  range %v to ensure 24h reservation margin", tomorrow, rnge, now)
			rnge.Start = now
		}
		newTargetValue = append(newTargetValue, rnge)
	}
	targets[tomorrow] = newTargetValue

	// Handle valid targets
	for day, target := range targets {
		if len(target) < 1 {
			continue
		}
		target_err := a.handleTargetList(value, day, target)
		if target_err != nil {
			logging.Out().Err(target_err).Msgf("%-10v (failure) unable to fulfill request for %v", day.String(), target)
		} else {
			logging.Out().Info().Msgf("%-10v (success) %v fullfilled", day.String(), target)
		}
	}

	return err
}

func (a *Agent) RunWithScheduler() (err error) {
	sched := quartz.NewStdScheduler()
	sched.Start()

	runJob := quartz.NewFunctionJob(func() (int, error) { return 0, a.Run() })

	triggers := []*util.CronTrigger{util.CronSaturdayAt10, util.CronEvery15Minutes}

	if vars.Has(vars.Debug) {
		trigger := util.CronEvery15Seconds
		logging.Out().Debug().Msgf("debug enabled, including high-frequency schedule (%v)", trigger.Expression())
		triggers = append(triggers, trigger)
	}

	if vars.Has(vars.CustomCron) {
		trigger := util.CronCustom
		logging.Out().Debug().Msgf("custom cron detected, including schedule (%v)", trigger.Expression())
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
