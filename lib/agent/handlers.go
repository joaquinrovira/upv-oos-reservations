package agent

import (
	"fmt"
	"strings"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib/logging"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/timerange"
	"github.com/joaquinrovira/upv-oos-reservations/lib/requests"
	"github.com/joaquinrovira/upv-oos-reservations/lib/vars"
)

func (a *Agent) handleTargetList(reservations *model.ReservationsWeek, day time.Weekday, targetList []timerange.TimeRange) (err error) {
	for _, target := range targetList {
		err = a.handleTarget(reservations, day, target)

		if err == nil {
			return
		}
	}

	return err
}

func (a *Agent) handleTarget(reservations *model.ReservationsWeek, day time.Weekday, timerange timerange.TimeRange) error {
	slot := reservations.FindSlotWithStrategy(day, timerange, model.FirstAvailable)

	// Validate slot
	if slot == nil {
		return fmt.Errorf("no empty time slot available")
	}
	if slot.Availability == 0 {
		return fmt.Errorf("no availability on this slot")
	}

	if vars.Has(vars.TestMode) {
		logging.Out().Debug().Msg("test mode enabled - skipping sending reservation request")
		return nil
	}

	// Send reservation request
	res, err := requests.ReserveSlot(a.client, slot)
	if err != nil {
		return err
	}
	
	if res.StatusCode < 200 || 400 <= res.StatusCode {
		return fmt.Errorf("response with status code %v", res.Status)
	}
	// TODO: Refactor error detection
	if 300 < res.StatusCode && res.StatusCode < 400 {
		// find redirect header
		if params, isErr := strings.CutPrefix(res.Header.Get("Location"), "sic_depact.pagina_de_error?"); isErr {
			urlParams := strings.Split(params, "&")
			for _, param := range urlParams {
				if err, isErr := strings.CutPrefix(param, "p_msg_error="); isErr {
					return fmt.Errorf("%v", err)
				}
			}
		}
	}

	return nil
}
