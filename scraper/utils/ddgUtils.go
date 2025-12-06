package utils

import (
	"net/http"
	"net/url"
)

func UnwrapDuckDuckGoURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if u.Host != "duckduckgo.com" {
		return raw
	}

	q := u.Query().Get("uddg")
	if q == "" {
		return raw
	}

	decoded, err := url.QueryUnescape(q)
	if err != nil {
		return raw
	}
	return decoded
}

func SetRequestHeaders(req *http.Request) {
	// Headers to look like a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("DNT", "1") // Do Not Track
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	req.Header.Set("sec-ch-ua", "'Chromium';v='142', 'Google Chrome';v='142', 'Not_A Brand';v='99'")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-ch-ua-platform", "'Linux'")

	req.Header.Set("sec-fetch-site", "same-origin")
}
