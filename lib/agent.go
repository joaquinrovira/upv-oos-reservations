package lib

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib/logging"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/timerange"
	"github.com/joaquinrovira/upv-oos-reservations/lib/requests"
	"github.com/joaquinrovira/upv-oos-reservations/lib/vars"
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

	logging.Log().Debug().Msg("Running...")

	value, err := a.GetReservationsData()
	if err != nil {
		return err
	}

	// Remove targets that have already been fulfilled
	reservations := value.GetReservations()
	for day := range reservations {
		logging.Log().Info().Msgf("%v already fulfilled in time slot %v, skipping", day.String(), reservations[day].At)
		delete(a.Cfg.Target, day)
	}

	for day, target := range a.Cfg.Target {
		target_err := a.handleTargetList(value, day, target)
		if target_err != nil {
			logging.Log().Err(target_err).Msgf("unable to fulfill request for %v %v", day.String(), target)
		} else {
			logging.Log().Info().Msgf("request %v %v fullfilled successfully", day.String(), target)
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

	return err
}

func (a *Agent) handleTarget(reservations *model.ReservationsWeek, day time.Weekday, timerange timerange.TimeRange) error {
	slot := reservations.FindSlotWithStrategy(day, timerange, model.FirstAvailable)

	// Validate slot
	if slot == nil {
		return fmt.Errorf("no empty time slot available")
	}
	if slot.Availability == 0 {
		return fmt.Errorf("no availability on this slot")
	}

	if vars.Get(vars.Debug) != "" {
		logging.Log().Debug().Msg("debug enabled, skipping sending reservation request")
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
