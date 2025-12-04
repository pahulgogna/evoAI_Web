package searching

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"scraper/src/utils"
	"strings"
	"time"
)

func httpGETFromDDG(url string, dnsAddress string) (*http.Response, error) {
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    return nil, err
  }

  utils.SetRequestHeaders(req)

  client := http.Client{
    Timeout: 10 * time.Second,
    Transport: utils.GetTransportForRequest(dnsAddress),
  }
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

func performDdgSearch(query string, dnsAddress string) {
  root := "https://duckduckgo.com/html/"

  url := fmt.Sprintf("%s?q=%s%s", root, url.QueryEscape(query), "&ia=web")


  res, err := httpGETFromDDG(url, dnsAddress)
  if err != nil || strings.Contains((*res).Status, "202") {
    fmt.Println("Error while searching the web", err)
    return
  }

  htmlBytes, err := debugPrintBody(res)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }

  rootHtml, err := utils.FromHtmlBytesToRoot(htmlBytes)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }

  resultsDiv := utils.FindDivByID(rootHtml, "links")
  if resultsDiv == nil {
    return
  }

  utils.FindLinks(resultsDiv, true, 0, query)
}
