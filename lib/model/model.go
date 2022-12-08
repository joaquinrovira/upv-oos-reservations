package model

import (
	"net/url"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib/model/daytime"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/timerange"
	"github.com/joaquinrovira/upv-oos-reservations/lib/util"
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

type Reservation struct {
	Slot ReservationSlot
	At   timerange.TimeRange
}

func (r *ReservationsWeek) GetReservations() (reservations map[time.Weekday]Reservation) {
	reservations = make(map[time.Weekday]Reservation)

	for _, day := range util.DaysOfWeek {
		for idx, slot := range r.Slots[day] {
			if slot.Reserved {
				at := r.SlotTimes[idx]
				reservations[day] = Reservation{slot, at}
				break
			}
		}
	}

	return
}
