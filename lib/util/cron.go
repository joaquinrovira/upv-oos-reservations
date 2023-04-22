package util

import (
	"time"

	"github.com/lnquy/cron"
	"github.com/reugn/go-quartz/quartz"
)

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

	ct.quartz, err = quartz.NewCronTriggerWithLoc(ct.expression, time.Local)
	if err != nil {
		return
	}

	return
}

func MustCronTrigger(expr string) (ct *CronTrigger) {
	ct = new(CronTrigger)
	ct.expression = expr
	var err error

	exprDesc, _ := cron.NewDescriptor()
	ct.description, err = exprDesc.ToDescription(ct.expression, cron.Locale_en)
	if err != nil {
		panic(err)
	}

	ct.quartz, err = quartz.NewCronTriggerWithLoc(ct.expression, time.Local)
	if err != nil {
		panic(err)
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

// Description returns the description of the trigger.
func (ct *CronTrigger) Source() *quartz.CronTrigger {
	return ct.quartz
}
