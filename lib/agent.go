package lib

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/internal"
	"github.com/joaquinrovira/upv-oos-reservations/internal/model"
	"github.com/joaquinrovira/upv-oos-reservations/internal/requests"
)

type Agent struct {
	Cfg    Config
	Client *http.Client
}

type Config struct {
	User   string
	Pass   string
	Target []TargetValue
}

type TargetValue struct {
	Day  time.Weekday
	Time model.TimeOfDay
}

func (t TargetValue) Format(f fmt.State, c rune) {
	f.Write([]byte(fmt.Sprintf("%s at %s", t.Day.String(), t.Time.String())))
}

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
	for _, v := range c.Target {
		err = v.Time.Validate()
		if err != nil {
			return
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

	for _, target := range a.Cfg.Target {
		err = a.handleTarget(value, target)
		if err != nil {
			internal.Log().Err(err).Msg(fmt.Sprintf("unable to fulfill request for %v", target))
		} else {
			internal.Log().Info().Msg(fmt.Sprintf("[OK] request %v fullfilled successfully", target))
		}
	}

	return err
}

func (a *Agent) handleTarget(reservations *model.ReservationsWeek, target TargetValue) error {
	slot := reservations.SlotAt(target.Day, target.Time)

	// Validate slot
	if slot == nil {
		return fmt.Errorf("no time slot available")
	}
	if slot.Availability == 0 {
		return fmt.Errorf("no availability on this slot")
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
