package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"golang.org/x/net/html"
)

var ScrapeWg sync.WaitGroup

func sendRequest(url string) (*html.Node, Page) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, Page{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	
	if err != nil {
		return nil, Page{}
	}

	rootNode, err := html.Parse(strings.NewReader(string(body)))
    if err != nil {
        return nil, Page{}
    }

	var page Page = Page{
		source: url,
		body: string(body),
	}
	
	return rootNode, page
}

func findLinks(n *html.Node) (links []string) {

	if isSkipableDiv(n) {
		return []string{}
	}
    
	if n.Type == html.ElementNode && (n.Data == "a") {
        for _, attr := range n.Attr {
            if attr.Key == "href" {
				isValid, url := isUrlValid(attr.Val, "")
				if isValid {
					links = append(links, url)
				}
            }
        }
    }

    for c := n.FirstChild; c != nil; c = c.NextSibling {
        links = append(links, findLinks(c)...)
    }
    return links
}

func scrape(url string, level int, maxLinksPerPage int) bool {
	
	MapMutex.Lock()
	if _, ok := seen[url]; ok || level == -1 {
		MapMutex.Unlock()
		return !ok
	}

	seen[url] = 1
	MapMutex.Unlock()

	rootNode, page := sendRequest(url)

	if rootNode == nil {
		return false
	}

	OutMutex.Lock()
	output = append(output, page)
	OutMutex.Unlock()

	
	links := findLinks(rootNode)
	numLinks := len(links)

	links = links[(numLinks/5):(numLinks * 3)/5]
	if len(links) > 10 {
		links = links[:maxLinksPerPage]
	}

	for i, link := range links {
		
		if i == maxLinksPerPage - 1 {
			break
		}

		valid, link := isUrlValid(link, url)
		if !valid {
			continue
		}
		ScrapeWg.Add(1)

		go func(l string, lev int) {
			defer ScrapeWg.Done()
			scrape(l, lev, maxLinksPerPage)
		}(link, level-1)
	}
	fmt.Printf("scraped url %s, level: %d\n", url, level)
	return true
}

func isUrlValid(url string, root string) (bool, string) {

	url = strings.TrimSpace(url)

	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {
		return true, url
	} 
	if strings.HasPrefix(url, "#") {
		return false, url
	}
	if strings.HasPrefix(url, "//") {
		return true, fmt.Sprintf("https:%s", url)
	}
	if strings.HasPrefix(url, "/") {
		if root == "" {
			return false, ""
		}
		return true, fmt.Sprintf("%s%s", root, url)
	}

	return false, url
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