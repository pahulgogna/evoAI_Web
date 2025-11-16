package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func httpGETFromDDG(url string) (*http.Response, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    // Headers to look like a real browser
    req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
    req.Header.Set("Accept-Language", "en-US,en;q=0.9")
    // req.Header.Set("Referer", "https://duckduckgo.com/")
    req.Header.Set("DNT", "1") // Do Not Track
    req.Header.Set("Upgrade-Insecure-Requests", "1")
    // req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
    
    req.Header.Set("sec-ch-ua", "'Chromium';v='142', 'Google Chrome';v='142', 'Not_A Brand';v='99'")
    req.Header.Set("priority", "u=0, i")
    req.Header.Set("cache-control", "max-age=0")
    req.Header.Set("sec-fetch-dest", "document")
    req.Header.Set("sec-ch-ua-platform", "'Linux'")

    req.Header.Set("sec-fetch-site", "same-origin")

    client := http.Client{}
    fmt.Println("sending req")
    return client.Do(req)
}


func debugPrintBody(res *http.Response) ([]byte, error) {
    if res == nil {
        return nil, fmt.Errorf("nil response")
    }
    defer res.Body.Close()

    var reader io.Reader = res.Body

    switch res.Header.Get("Content-Encoding") {
    case "gzip":
        gz, err := gzip.NewReader(res.Body)
        if err != nil {
            return nil, fmt.Errorf("creating gzip reader: %w", err)
        }
        defer gz.Close()
        reader = gz
    }

    bodyBytes, err := io.ReadAll(reader)
    if err != nil {
        return nil, fmt.Errorf("read body: %w", err)
    }

    return bodyBytes, nil
}

func performDdgSearch(query string) []string {

	root := "https://duckduckgo.com/html/"

	url := fmt.Sprintf("%s?q=%s%s", root, url.QueryEscape(query), "&ia=web")

	
	res, err := httpGETFromDDG(url)
	if err != nil || strings.Contains((*res).Status, "202") {
		fmt.Println("Error while searching the web", err)
		return []string{}
	}

	htmlBytes, err := debugPrintBody(res)
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}
	}

	rootHtml, err := fromHtmlBytesToRoot(htmlBytes)
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}
	}

    resultsDiv := findDivByID(rootHtml, "links")
    if resultsDiv == nil {
        return nil
    }

    links := []string{}

	for _, link := range findLinks(resultsDiv) {
        links = append(links, unwrapDuckDuckGoURL(link))
    }

    return links
}

func unwrapDuckDuckGoURL(raw string) string {
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