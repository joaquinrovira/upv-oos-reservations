package requests

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/joaquinrovira/upv-oos-reservations/lib/model/login"
)

func Login(c *http.Client, username, password string, loginType login.LoginType) (err error) {
	switch loginType {
	case login.StudentLogin:
		return StudentLogin(c, username, password)
	case login.PersonnelLogin:
		return PersonnelLogin(c, username, password)
	default:
		return fmt.Errorf("invalid LoginType '%s'", loginType)
	}
}

func StudentLogin(c *http.Client, username, password string) (err error) {
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

func PersonnelLogin(c *http.Client, username, password string) (err error) {
	reader := strings.NewReader(fmt.Sprintf("id=c&X=500&vista=&param=&cua=miupv&dni=%s&rclau=%s", username, password))
	req, err := http.NewRequest("POST", "https://intranet.upv.es/pls/soalu/est_aute.intrapercomp", reader)

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
