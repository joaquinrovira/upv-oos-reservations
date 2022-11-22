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

type TimeRange struct {
	StartTime TimeOfDay
	EndTime   TimeOfDay
}

func (tr TimeRange) Validate() error {
	if tr.StartTime.Value() > tr.EndTime.Value() {
		return fmt.Errorf("start time (%v) greater than end time (%v)", tr.StartTime, tr.EndTime)
	}
	return nil
}

func (tr TimeRange) timeInRange(t TimeOfDay) bool {
	return tr.StartTime.Value() <= t.Value() && t.Value() < tr.EndTime.Value()
}

type ReservationSlot struct {
	Name         string
	Availability int64
	Reserved     bool
	URL          *url.URL
}

type ReservationsWeek struct {
	SlotTimes []TimeRange
	Slots     map[time.Weekday][]ReservationSlot
}

type SlotAtSearchStrategy int

const (
	Containts SlotAtSearchStrategy = iota
	NotBefore
)

type SlotWithinSearchStrategy int

const (
	First SlotWithinSearchStrategy = iota
)

func (r *ReservationsWeek) SlotAtContaints(day time.Weekday, t TimeOfDay) *ReservationSlot {
	for idx, slot := range r.Slots[day] {
		if r.SlotTimes[idx].timeInRange(t) {
			return &slot
		}
	}
	return nil
}

func (r *ReservationsWeek) SlotAtNotBefore(day time.Weekday, t TimeOfDay) *ReservationSlot {
	for idx, slot := range r.Slots[day] {
		if r.SlotTimes[idx].StartTime.Value() >= t.Value() {
			return &slot
		}
	}
	return nil
}

func (r *ReservationsWeek) SlotAtWithStrategy(day time.Weekday, t TimeOfDay, strat SlotAtSearchStrategy) *ReservationSlot {

	switch strat {
	case Containts:
		return r.SlotAtContaints(day, t)
	case NotBefore:
		return r.SlotAtNotBefore(day, t)
	}
	return nil
}

func (r *ReservationsWeek) SlotAt(day time.Weekday, t TimeOfDay) *ReservationSlot {
	return r.SlotAtWithStrategy(day, t, Containts)
}

func (r *ReservationsWeek) SlotWithinFirst(day time.Weekday, t TimeRange) *ReservationSlot {
	for idx, slot := range r.Slots[day] {
		slotTime := r.SlotTimes[idx]
		if t.StartTime.Value() <= slotTime.StartTime.Value() && slotTime.StartTime.Value() <= t.EndTime.Value() {
			return &slot
		}
	}
	return nil
}

func (r *ReservationsWeek) SlotWithinWithStrategy(day time.Weekday, t TimeRange, strat SlotWithinSearchStrategy) *ReservationSlot {

	switch strat {
	case First:
		return r.SlotWithinFirst(day, t)
	}
	return nil
}

func (r *ReservationsWeek) SlotWithin(day time.Weekday, t TimeRange) *ReservationSlot {
	return r.SlotWithinWithStrategy(day, t, First)
}

func NewReservarionsWeek() *ReservationsWeek {
	var r ReservationsWeek
	r.Slots = make(map[time.Weekday][]ReservationSlot)
	return &r
}
