package requests

import "net/http"

func setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0")
	req.Header.Set("Origin", "https://intranet.upv.es")
}
