package requests

import (
	"fmt"
	"net/http"
	"strings"
)

func Login(c *http.Client, username, password string) (err error) {
	reader := strings.NewReader(fmt.Sprintf("id=c&estilo=500&vista=&param=&cua=miupv&dni=%s&clau=%s", username, password))
	req, err := http.NewRequest("POST", "https://intranet.upv.es/pls/soalu/est_aute.intraalucomp", reader)

	if err != nil {
		return
	}

	setCommonHeaders(req)

	// Specific headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://intranet.upv.es/pls/soalu/est_intranet.Ni_portal_n")

	res, err := c.Do(req)
	if err != nil {
		return
	}
	if res.Header.Get("Set-Cookie") == "" {
		return fmt.Errorf("unable to obtain authentication cookies")
	}

	return
}
