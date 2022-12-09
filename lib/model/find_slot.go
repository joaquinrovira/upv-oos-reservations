package model

import (
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib/model/daytime"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/timerange"
)

type SlotSearchStrategy int

const (
	First SlotSearchStrategy = iota
	FirstAvailable
	FirstAvailableAtLeast1HourAway
)

func FirstAvailableCondition(slot ReservationSlot, slotTime timerange.TimeRange, day time.Weekday, t timerange.TimeRange) bool {
	if t.Start.Value() > slotTime.Start.Value() || slotTime.Start.Value() > t.End.Value() {
		return false
	}
	if slot.Availability == 0 {
		return false
	}
	return true
}

func (r *ReservationsWeek) FindSlotFirstAvailable(day time.Weekday, t timerange.TimeRange) *ReservationSlot {
	for idx, slot := range r.Slots[day] {
		slotTime := r.SlotTimes[idx]
		if FirstAvailableCondition(slot, slotTime, day, t) {
			return &slot
		}
	}
	return nil
}

func FirstAvailableAtLeast1HourAwayCondition(slot ReservationSlot, slotTime timerange.TimeRange, day time.Weekday, t timerange.TimeRange) bool {
	anHourAway, _ := daytime.FromTime(time.Now().Add(time.Hour))

	if slotTime.Start.Value() < anHourAway.Value() {
		return false
	}

	if t.Start.Value() > slotTime.Start.Value() || slotTime.Start.Value() > t.End.Value() {
		return false
	}
	if slot.Availability == 0 {
		return false
	}
	return true
}

func (r *ReservationsWeek) FindSlotFirstAvailableAtLeast1HourAway(day time.Weekday, t timerange.TimeRange) *ReservationSlot {
	for idx, slot := range r.Slots[day] {
		slotTime := r.SlotTimes[idx]
		if FirstAvailableAtLeast1HourAwayCondition(slot, slotTime, day, t) {
			return &slot
		}
	}
	return nil
}

func (r *ReservationsWeek) FindSlotFirst(day time.Weekday, t timerange.TimeRange) *ReservationSlot {
	for idx, slot := range r.Slots[day] {
		slotTime := r.SlotTimes[idx]
		if t.Start.Value() > slotTime.Start.Value() || slotTime.Start.Value() > t.End.Value() {
			return &slot
		}
	}
	return nil
}

func (r *ReservationsWeek) FindSlotWithStrategy(day time.Weekday, t timerange.TimeRange, strat SlotSearchStrategy) *ReservationSlot {

	switch strat {
	case First:
		return r.FindSlotFirst(day, t)
	case FirstAvailable:
		return r.FindSlotFirstAvailable(day, t)
	case FirstAvailableAtLeast1HourAway:
		return r.FindSlotFirstAvailableAtLeast1HourAway(day, t)
	}
	return nil
}

func (r *ReservationsWeek) FindSlot(day time.Weekday, t timerange.TimeRange) *ReservationSlot {
	return r.FindSlotWithStrategy(day, t, FirstAvailable)
}
