package timerange

import (
	"fmt"

	"github.com/joaquinrovira/upv-oos-reservations/internal/model/daytime"
)

type TimeRange struct {
	Start daytime.DayTime
	End   daytime.DayTime
}

func (tr TimeRange) Validate() error {
	if tr.Start.Value() > tr.End.Value() {
		return fmt.Errorf("start time (%v) greater than end time (%v)", tr.Start, tr.End)
	}
	return nil
}

func (tr TimeRange) Contains(t daytime.DayTime) bool {
	return tr.Start.Value() <= t.Value() && t.Value() < tr.End.Value()
}
