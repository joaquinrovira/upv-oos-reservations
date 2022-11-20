package model

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// var intToWeekday = []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}

// Expects string "HH:MM"
func ParseSlotTime(s string) (t TimeOfDay, err error) {
	units := strings.Split(s, ":")

	if len(units) != 2 {
		err = fmt.Errorf("invalid input '%s' expected format like HH:MM", s)
		return
	}

	t, err = NewTimeOfDayString(units[0], units[1])

	return
}

type SlotTimeRange struct {
	StartTime TimeOfDay
	EndTime   TimeOfDay
}

func (r SlotTimeRange) timeInRange(t TimeOfDay) bool {
	return r.StartTime.Value() <= t.Value() && t.Value() < r.EndTime.Value()
}

type ReservationSlot struct {
	Name         string
	URL          *url.URL
	Availability int64
	Reserved     bool
}

type ReservationsWeek struct {
	SlotTimes []SlotTimeRange
	Slots     map[time.Weekday][]ReservationSlot
}

func (r *ReservationsWeek) SlotAt(day time.Weekday, t TimeOfDay) *ReservationSlot {
	for idx, slot := range r.Slots[day] {
		if r.SlotTimes[idx].timeInRange(t) {
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
