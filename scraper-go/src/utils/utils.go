package utils

import (
	"io"
	"scraper/src/customTypes"
	"scraper/src/global"
	"strings"

	"golang.org/x/net/html"
)

func sendRequest(url string, scraper *global.ScraperSession) (*html.Node, customTypes.Page) {

	resp, err := scraper.Client.Get(url)

	if err != nil {
		return nil, customTypes.Page{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, customTypes.Page{}
	}

	rootNode, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, customTypes.Page{}
	}

	var page customTypes.Page = customTypes.Page{
		Source: url,
		Body:   string(body),
	}

	return rootNode, page
}

func Scrape(link *customTypes.StoreUrl, query string, scraper *global.ScraperSession) bool {

	scraper.Queue.IncrementProcessCounter()
	rootNode, page := sendRequest(link.Url, scraper)

	if rootNode == nil {
		return false
	}
	
	page.Body = ParseHtmlToContent(page.Body)
	
	if len(page.Body) <= 200 {
		return false
	}

	scraper.Mutex.Lock()
	scraper.Output = append(scraper.Output, page)
	scraper.Mutex.Unlock()

	FindLinks(rootNode, false, link.Level+1, query, scraper)

	return true
}

func FromHtmlBytesToRoot(body []byte) (*html.Node, error) {
	rootNode, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	return rootNode, err
}

func FindDivByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if res := FindDivByID(c, id); res != nil {
			return res
		}
	}
	return nil
}

func isSkipableDiv(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "feedback-btn" {
				return true
			}
		}
	}

	if n.Type == html.ElementNode {
		if n.Data == "script" || n.Data == "style" || n.Data == "head" {
			return true
		}
	}

	return false
}

func normalizeWhitespace(input string) string {
	fields := strings.Fields(input)
	return strings.Join(fields, " ")
}


