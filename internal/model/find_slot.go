package model

import (
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/internal/model/timerange"
)

type SlotSearchStrategy int

const (
	First SlotSearchStrategy = iota
	FirstAvailable
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
	case FirstAvailable:
		return r.FindSlotFirstAvailable(day, t)
	case First:
		return r.FindSlotFirst(day, t)
	}
	return nil
}

func (r *ReservationsWeek) FindSlot(day time.Weekday, t timerange.TimeRange) *ReservationSlot {
	return r.FindSlotWithStrategy(day, t, FirstAvailable)
}
