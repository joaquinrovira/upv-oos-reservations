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
	Time float64 // TODO: find a better way to do this
}

func (t TargetValue) Format(f fmt.State, c rune) {
	hour := int64(t.Time)
	minute := int64((t.Time - float64(hour)) * 60)
	f.Write([]byte(fmt.Sprintf("%s at %2d:%2d", t.Day.String(), hour, minute)))
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
		_, err = model.NewSlotTime(v.Time)
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
		day := target.Day
		time, _ := model.NewSlotTime(target.Time) // error can be ignored, values checked on New()
		slot := value.SlotAt(day, time)

		// Validate slot
		if slot == nil {
			internal.Log().Err(fmt.Errorf("no time slot available")).Msg(fmt.Sprintf("unable to fulfill request for %v", target))
			continue
		}
		if slot.Availability == 0 {
			internal.Log().Err(fmt.Errorf("no availability on this slot")).Msg(fmt.Sprintf("unable to fulfill request for %v", target))
			continue
		}
		if slot.URL == nil {
			internal.Log().Err(fmt.Errorf("missing reservation URL on slot")).Msg(fmt.Sprintf("unable to fulfill request for %v", target))
			continue
		}

		// Send reservation request
		res, err := requests.ReserveSlot(a.Client, slot)
		if err != nil {
			internal.Log().Err(err).Msg(fmt.Sprintf("unable to fulfill request for %v", target))
			continue
		}
		if res.StatusCode < 200 || 400 <= res.StatusCode {
			internal.Log().Err(fmt.Errorf("response with status code %v", res.Status)).Msg(fmt.Sprintf("unable to fulfill request for %v", target))
			continue
		}

		internal.Log().Info().Msg(fmt.Sprintf("[%v] request %v fullfilled successfully", res.Status, target))
	}

	return err
}
