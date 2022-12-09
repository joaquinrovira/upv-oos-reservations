package timerange

import (
	"fmt"

	"github.com/joaquinrovira/upv-oos-reservations/lib/model/daytime"
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

func (tr TimeRange) String() string {
	return fmt.Sprintf("%v-%v", tr.Start, tr.End)
}

func (tr TimeRange) Format(f fmt.State, c rune) {
	f.Write([]byte(tr.String()))
}
