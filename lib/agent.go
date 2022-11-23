package lib

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/internal"
	"github.com/joaquinrovira/upv-oos-reservations/internal/model"
	"github.com/joaquinrovira/upv-oos-reservations/internal/model/timerange"
	"github.com/joaquinrovira/upv-oos-reservations/internal/requests"
	"github.com/joaquinrovira/upv-oos-reservations/internal/vars"
)

type Agent struct {
	Cfg    Config
	Client *http.Client
}

type Config struct {
	User   string
	Pass   string
	Target TargetValue
}

/*
For each weekday, we can define a list of time ranges we want to try
to make a reservation in. The list is sorted by descreasing priority.
If a slot is open within the first time range, a reservation will be
executed for that slot.
*/
type TargetValue map[time.Weekday][]timerange.TimeRange

func New(c Config) (a *Agent, err error) {
	if err = checkConfig(c); err != nil {
		return
	}

	jar, err := cookiejar.New(nil)
	client := http.Client{
		// Enable cookie jar for user authentication
		Jar: jar,
		// Do not follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	a = &Agent{Cfg: c, Client: &client}
	return
}

// Returns error for invalid configurations
func checkConfig(c Config) (err error) {
	for _, trList := range c.Target {
		for _, tr := range trList {
			err = tr.Validate()
			if err != nil {
				return
			}
		}
	}

	return
}

func (a *Agent) Run() (err error) {
	a.Login()

	value, err := a.GetReservationsData()
	if err != nil {
		return err
	}

	// TODO: if you run the program twice, if can reserve two times in the same day
	// FIX IT!
	for day, target := range a.Cfg.Target {
		err = a.handleTargetList(value, day, target)
		if err != nil {
			internal.Log().Err(err).Msg(fmt.Sprintf("unable to fulfill request for %v", target))
		} else {
			internal.Log().Info().Msg(fmt.Sprintf("[OK] request %v %v fullfilled successfully", day.String(), target))
		}
	}

	return err
}

func (a *Agent) handleTargetList(reservations *model.ReservationsWeek, day time.Weekday, targetList []timerange.TimeRange) (err error) {
	for _, target := range targetList {
		err = a.handleTarget(reservations, day, target)

		if err == nil {
			return
		}
	}

	return fmt.Errorf("unable to fulfill target list")
}

func (a *Agent) handleTarget(reservations *model.ReservationsWeek, day time.Weekday, timerange timerange.TimeRange) error {
	slot := reservations.FindSlot(day, timerange)

	// Validate slot
	if slot == nil {
		return fmt.Errorf("no time slot available")
	}
	if slot.Availability == 0 {
		return fmt.Errorf("no availability on this slot")
	}

	if vars.Get(vars.Debug) != "" {
		internal.Log().Debug().Msg("debug enabled, skipping sending reservation request")
		return nil
	}

	// Send reservation request
	res, err := requests.ReserveSlot(a.Client, slot)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || 400 <= res.StatusCode {
		return fmt.Errorf("response with status code %v", res.Status)
	}

	return nil
}
