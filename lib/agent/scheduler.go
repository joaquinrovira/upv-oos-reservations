package agent

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib/logging"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/daytime"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/timerange"
	"github.com/joaquinrovira/upv-oos-reservations/lib/util"
	"github.com/joaquinrovira/upv-oos-reservations/lib/vars"
	"github.com/reugn/go-quartz/quartz"
)

var lock sync.Mutex
var triggers []*util.CronTrigger

func registerCron(cron string) {
	if triggers == nil {
		triggers = make([]*util.CronTrigger, 0)
	}

	if trigger, err := util.NewCronTrigger(strings.TrimSpace(cron)); err != nil {
		logging.Out().Fatal().Msgf("'%s' %v", cron, err)
	} else {
		triggers = append(triggers, trigger)
		logging.Out().Info().Msgf("registered job trigger %-20v", fmt.Sprintf("(%v)", trigger.Expression()))
		logging.Out().Info().Msgf("description: \"%v\"", trigger.Description())
	}
}

func init() {
	logging.Out().Info().Msg("+-----------------------+")
	logging.Out().Info().Msg("| Registering CRON jobs |")
	logging.Out().Info().Msg("+-----------------------+")
	tz, _ := time.Now().Zone()
	logging.Out().Debug().Msgf("TZ=%v", tz)

	// Parse and validate Crons
	// Split by ',' and trim
	var split = strings.FieldsFunc(vars.Get(vars.DefaultCron), func(r rune) bool { return r == ',' })
	for _, cron := range split {
		logging.Out().Info().Msgf("env.%-15s adding default cron '%s'", vars.DefaultCron, cron)
		registerCron(cron)
	}

	// Parse and validate CustomCron
	if vars.Has(vars.CustomCron) {
		cron := vars.Get(vars.CustomCron)
		logging.Out().Info().Msgf("env.%-15s adding custom cron '%s'", vars.CustomCron, cron)
		registerCron(cron)
	}

	// Add high-frequency job when DEBUG env var is present
	if vars.Has(vars.HighFqTrigger) {
		var cron string = "*/10 * * * * *" // Every 10 seconds
		logging.Out().Info().Msgf("env.%-15s enabled adding debug cron '%s'", vars.HighFqTrigger, cron)
		registerCron(cron)
	}
}

func runJobScheduleWrapper(a *Agent, trigger *util.CronTrigger) quartz.Job {
	return quartz.NewFunctionJob(func() (int, error) {
		logging.Out().Debug().Msgf("👀 \t %s", trigger.Expression())
		lock.Lock()
		logging.Out().Debug().Msgf("🔒 \t %s", trigger.Expression())
		defer lock.Unlock()
		defer logging.Out().Debug().Msgf("🔓 \t %s", trigger.Expression())

		logging.Out().Info().Msg("+------------------------------------------+")
		logging.Out().Info().Msgf("| Reservation trigger %20s |", trigger.Expression())
		logging.Out().Info().Msg("+------------------------------------------+")

		next, _ := trigger.NextFireTime(time.Now().UnixNano())
		defer logging.Out().Debug().Msgf("next execution in %v", time.Until(time.Unix(0, next)).Round(time.Second))

		return 0, a.Run()
	})

}

func (a *Agent) RunWithScheduler() (err error) {
	sched := quartz.NewStdScheduler()
	sched.Start()

	for _, trigger := range triggers {
		job := runJobScheduleWrapper(a, trigger)
		sched.ScheduleJob(job, trigger)
	}
	<-a.ctx.Done()
	sched.Stop()

	return
}

func (a *Agent) Run() (err error) {
	a.Login()

	targets := a.target.Clone()

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
			logging.Out().Info().Msgf("%-10v (avoiding) range %v with less than 24h reservation margin", tomorrow, rnge)
			continue
		}

		if rnge.Start.Value() < now.Value() {
			logging.Out().Info().Msgf("%-10v (modifying) range %v to range %v to ensure 24h reservation margin", tomorrow, rnge, now)
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

	logging.Out().Info().Msg("finished reservation attempt!")

	return err
}
