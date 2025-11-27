package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"golang.org/x/net/html"
	"scraper/src/customTypes"
	"scraper/src/urlManager"
)

var ScrapeWg sync.WaitGroup

func sendRequest(url string) (*html.Node, customTypes.Page) {

	resp, err := http.Get(url)

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
		Body: string(body),
	}
	
	return rootNode, page
}

func findLinks(n *html.Node, filterDDG bool, level int) {

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
					urlManager.AddUrl(url, Query, level)
				}
            }
        }
    }

    for c := n.FirstChild; c != nil; c = c.NextSibling {
        findLinks(c, filterDDG, level)
    }
}

func scrape(link *customTypes.StoreUrl, totalLinksToProcess int, maxLinksPerPage int) bool {
	
	if urlManager.GetProcessedUrlsCount() >= totalLinksToProcess {
		return false
	}
	fmt.Println("Processed urls: ",  urlManager.GetProcessedUrlsCount())
	
	urlManager.IncrementProcessCounter()
	rootNode, page := sendRequest(link.Url)

	if rootNode == nil {
		return false
	}

	OutMutex.Lock()
	Output = append(Output, page)
	OutMutex.Unlock()

	findLinks(rootNode, false, link.Level + 1)

	// for i := 0; i < maxLinksPerPage; i++ {

	// 	link := urlManager.GetUrl()
	// 	if link == nil {
	// 		continue
	// 	}
		
	// 	if i == maxLinksPerPage - 1 {
	// 		break
	// 	}

	// 	valid, l := urlManager.IsUrlValid(link.Url, link.Url)
	// 	if !valid {
	// 		continue
	// 	}
	// 	link.Url = l

	// 	ScrapeWg.Add(1)

	// 	go func(l *customTypes.StoreUrl) {
	// 		defer ScrapeWg.Done()
	// 		scrape(l, totalLinksToProcess, maxLinksPerPage)
	// 	}(link)
	// }
	
	return true
}

func fromHtmlBytesToRoot(body []byte) (*html.Node, error) {
	rootNode, err := html.Parse(strings.NewReader(string(body)))
    if err != nil {
        return nil, err
    }
	return rootNode, err
}

func findDivByID(n *html.Node, id string) *html.Node {
    if n.Type == html.ElementNode && n.Data == "div" {
        for _, attr := range n.Attr {
            if attr.Key == "id" && attr.Val == id {
                return n
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if res := findDivByID(c, id); res != nil {
            return res
        }
    }
    return nil
}

func isSkipableDiv(n *html.Node) bool {
	if (n.Type == html.ElementNode && n.Data == "div") {
		for _, attr := range n.Attr {
            if attr.Key == "class" && attr.Val == "feedback-btn" {
                return true
            }
        }
	}
	return false
}