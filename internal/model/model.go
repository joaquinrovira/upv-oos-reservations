package model

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type slotTime float64 // Represents a time of day as float a within range 0.00 - 24.00

func NewSlotTime(t float64) (slotTime, error) {
	if t < 0 || t > 24 {
		return 0, fmt.Errorf("invalid slot time, received value %f should be range [0,24]", t)
	}

	return slotTime(t), nil
}

// var intToWeekday = []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}

// string(XX:YY) -> float64(XX + YY/60)
func ParseSlotTime(s string) (slotTime, error) {
	units := strings.Split(s, ":")

	hours, err := strconv.ParseFloat(units[0], 32)
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.ParseFloat(units[1], 32)
	if err != nil {
		return 0, err
	}

	return NewSlotTime(hours + minutes/60)
}

type SlotTimeRange struct {
	StartTime slotTime
	EndTime   slotTime
}

func (r SlotTimeRange) timeInRange(t slotTime) bool {
	return r.StartTime <= t && t < r.EndTime
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

func (r *ReservationsWeek) SlotAt(day time.Weekday, t slotTime) *ReservationSlot {
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
