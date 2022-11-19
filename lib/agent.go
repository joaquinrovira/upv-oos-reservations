package lib

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

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

func New(c Config) (a *Agent, err error) {
	// Check config
	for _, v := range c.Target {
		_, err = model.NewSlotTime(v.Time)
		if err != nil {
			return
		}
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
			fmt.Printf("unable to fulfill requrest for %v, no time slot available\n", target)
			continue
		}
		if slot.Availability == 0 {
			fmt.Printf("unable to fulfill requrest for %v, no availability on this slot\n", target)
			continue
		}
		if slot.URL == nil {
			fmt.Printf("unable to fulfill requrest for %v, missing reservation URL on slot\n", target)
			continue
		}

		// Send reservation request
		res, err := requests.ReserveSlot(a.Client, slot)
		if err != nil {
			fmt.Printf("unable to fulfill requrest for %v, %s\n", target, err.Error())
			continue
		}
		if res.StatusCode < 200 || 400 <= res.StatusCode {
			fmt.Printf("unable to fulfill requrest for %v, response with status code %v\n", target, res.Status)
			continue
		}
		fmt.Println(res.Status, err)
	}

	return err
}
