package util

import (
	"github.com/joaquinrovira/upv-oos-reservations/lib/logging"
	"github.com/joaquinrovira/upv-oos-reservations/lib/vars"
	"github.com/lnquy/cron"
	"github.com/reugn/go-quartz/quartz"
)

var cronEveryHour string = "0 0 * * * *"
var cronEvery10Seconds string = "*/10 * * * * *"
var cronEvery15Minutes string = "0 */15 * * * *"
var cronSaturdayAt10 string = "*/10 0-15 10 * * SAT"

var CronEveryHour *CronTrigger
var CronEvery15Seconds *CronTrigger
var CronEvery15Minutes *CronTrigger
var CronSaturdayAt10 *CronTrigger
var CronCustom *CronTrigger

func init() {
	var err error
	CronEveryHour, err = NewCronTrigger(cronEveryHour)
	if err != nil {
		panic(err)
	}
	CronEvery15Seconds, err = NewCronTrigger(cronEvery10Seconds)
	if err != nil {
		panic(err)
	}
	CronEvery15Minutes, err = NewCronTrigger(cronEvery15Minutes)
	if err != nil {
		panic(err)
	}
	CronSaturdayAt10, err = NewCronTrigger(cronSaturdayAt10)
	if err != nil {
		panic(err)
	}

	if vars.Has(vars.CustomCron) {
		customCron := vars.Get(vars.CustomCron)
		CronCustom, err = NewCronTrigger(customCron)
		if err != nil {
			logging.Out().Fatal().Err(err).Msgf("unable to parse custom cron value (%s)", customCron)
		}
	}
}

type CronTrigger struct {
	quartz      *quartz.CronTrigger
	expression  string
	description string
}

func NewCronTrigger(expr string) (ct *CronTrigger, err error) {
	ct = new(CronTrigger)
	ct.expression = expr

	exprDesc, _ := cron.NewDescriptor()
	ct.description, err = exprDesc.ToDescription(ct.expression, cron.Locale_en)
	if err != nil {
		return
	}

	ct.quartz, err = quartz.NewCronTrigger(ct.expression)
	if err != nil {
		return
	}

	return
}

// NextFireTime returns the next time at which the CronTrigger is scheduled to fire.
func (ct *CronTrigger) NextFireTime(prev int64) (int64, error) {
	return ct.quartz.NextFireTime(prev)
}

// Description returns the description of the trigger.
func (ct *CronTrigger) Description() string {
	return ct.description
}

// Description returns the description of the trigger.
func (ct *CronTrigger) Expression() string {
	return ct.expression
}
