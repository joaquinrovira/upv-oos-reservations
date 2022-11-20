package model

import (
	"fmt"
	"strconv"
)

type TimeOfDay struct {
	Hours   uint8
	Minutes uint8
}

func NewTimeOfDay(hours, minutes uint8) (t TimeOfDay, err error) {
	// Parse hours
	if hours >= 24 {
		err = fmt.Errorf("hour '%d' out of range [0,24[", hours)
		return
	} else {
		t.Hours = uint8(hours)
	}

	// Parse minutes
	if minutes >= 60 {
		err = fmt.Errorf("minute '%d' out of range [0,60[", minutes)
		return
	} else {
		t.Minutes = uint8(minutes)
	}

	return
}

func NewTimeOfDayString(hours, minutes string) (t TimeOfDay, err error) {
	var hrs, mins uint64

	if hrs, err = strconv.ParseUint(hours, 10, 8); err != nil {
		return t, fmt.Errorf("unable to parse hours: %v", err)
	}
	if mins, err = strconv.ParseUint(minutes, 10, 8); err != nil {
		return t, fmt.Errorf("unable to parse minutes: %v", err)
	}

	return NewTimeOfDay(uint8(hrs), uint8(mins))
}

func (t TimeOfDay) String() string {
	return fmt.Sprintf("%2d:%2d", t.Hours, t.Minutes)
}

func (t TimeOfDay) Format(f fmt.State, c rune) {
	f.Write([]byte(t.String()))
}

func (t TimeOfDay) Value() uint16 {
	return uint16(t.Hours)*60 + uint16(t.Minutes)
}

func (t TimeOfDay) Validate() (err error) {
	_, err = NewTimeOfDay(t.Hours, t.Minutes)
	return
}
