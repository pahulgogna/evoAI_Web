package utils

import (
	"strings"

	"github.com/pahulgogna/evoAI_Web/scraper/src/global"
	"github.com/pahulgogna/evoAI_Web/scraper/src/urlManager"

	"golang.org/x/net/html"
)

func ParseHtmlToContent(htmlString string) string {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return ""
	}

	var sb strings.Builder

	var f func(*html.Node)
	f = func(n *html.Node) {

		if isSkipableDiv(n) {
			return
		}

		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
			sb.WriteString(" ")
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return normalizeWhitespace(sb.String())
}

func FindLinks(n *html.Node, filterDDG bool, level int, query string, scraper *global.ScraperSession) {

	if isSkipableDiv(n) {
		return
	}

	if n.Type == html.ElementNode && (n.Data == "a") {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				isValid, url := urlManager.IsUrlValid(attr.Val, "")

				if filterDDG {
					url = UnwrapDuckDuckGoURL(url)
				}

				if isValid {
					scraper.Queue.AddUrl(url, query, level)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		FindLinks(c, filterDDG, level, query, scraper)
	}
}
