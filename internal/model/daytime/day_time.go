package daytime

import (
	"fmt"
	"strconv"
	"strings"
)

type DayTime struct {
	Hour   uint8
	Minute uint8
}

func (t DayTime) Value() uint16 {
	return uint16(t.Hour)*60 + uint16(t.Minute)
}

func (t DayTime) Validate() (err error) {
	_, err = New(t.Hour, t.Minute)
	return
}

func (t DayTime) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

func (t DayTime) Format(f fmt.State, c rune) {
	f.Write([]byte(t.String()))
}

func New(hours, minutes uint8) (t DayTime, err error) {
	// Parse hours
	if hours >= 24 {
		err = fmt.Errorf("hour '%d' out of range [0,24[", hours)
		return
	} else {
		t.Hour = uint8(hours)
	}

	// Parse minutes
	if minutes >= 60 {
		err = fmt.Errorf("minute '%d' out of range [0,60[", minutes)
		return
	} else {
		t.Minute = uint8(minutes)
	}

	return
}

func FromString(hours, minutes string) (t DayTime, err error) {
	var hrs, mins uint64

	if hrs, err = strconv.ParseUint(hours, 10, 8); err != nil {
		return t, fmt.Errorf("unable to parse hours: %v", err)
	}
	if mins, err = strconv.ParseUint(minutes, 10, 8); err != nil {
		return t, fmt.Errorf("unable to parse minutes: %v", err)
	}

	return New(uint8(hrs), uint8(mins))
}

// Expects string "HH:MM"
func Parse(s string) (t DayTime, err error) {
	units := strings.Split(s, ":")

	if len(units) != 2 {
		err = fmt.Errorf("invalid input '%s' expected format like HH:MM", s)
		return
	}

	t, err = FromString(units[0], units[1])

	return
}
