package model

import (
	"net/url"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/internal/model/daytime"
	"github.com/joaquinrovira/upv-oos-reservations/internal/model/timerange"
)

type ReservationSlot struct {
	Name         string
	Availability int64
	Reserved     bool
	URL          *url.URL
}

type ReservationsWeek struct {
	SlotTimes []timerange.TimeRange
	Slots     map[time.Weekday][]ReservationSlot
}

func (r *ReservationsWeek) SlotAt(day time.Weekday, t daytime.DayTime) *ReservationSlot {
	for idx, slot := range r.Slots[day] {
		if r.SlotTimes[idx].Contains(t) {
			return &slot
		}
	}
	return nil
}

func NewReservarionsWeek() *ReservationsWeek {
	var r ReservationsWeek
	r.Slots = make(map[time.Weekday][]ReservationSlot)
	return &r
}
