package requests

import (
	"net/http"
	"net/url"

	"github.com/joaquinrovira/upv-oos-reservations/internal/vars"
)

func GetReservationsData(c *http.Client) (res *http.Response, err error) {
	baseUrl := "https://intranet.upv.es/pls/soalu/sic_depact.HSemActividades"

	v := url.Values{}
	v.Set("p_campus", "V")
	v.Set("p_tipoact", vars.Get(vars.ActivityType))
	v.Set("p_codacti", vars.Get(vars.ActivityCode))
	v.Set("p_vista", "intranet")
	v.Set("p_idioma", "c")
	v.Set("p_solo_matricula_sn", "")
	v.Set("p_anc", "filtro_actividad")
	encodedURL := baseUrl + "?" + v.Encode()

	req, err := http.NewRequest("GET", encodedURL, nil)
	if err != nil {
		return
	}

	setCommonHeaders(req)

	// Specific headers
	req.Header.Set("Referer", "https://intranet.upv.es/pls/soalu/sic_menu.alumno?P_IDIOMA=c")

	return c.Do(req)
}
