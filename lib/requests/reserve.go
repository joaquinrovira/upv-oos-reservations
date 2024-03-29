package requests

import (
	"net/http"

	"github.com/joaquinrovira/upv-oos-reservations/lib/model"
	"github.com/joaquinrovira/upv-oos-reservations/lib/util"
)

func ReserveSlot(c util.RLHTTPClient, s *model.ReservationSlot) (res *http.Response, err error) {
	req, err := http.NewRequest("GET", s.URL.String(), nil)
	if err != nil {
		return
	}
	setCommonHeaders(req)
	return c.Do(req)
}
