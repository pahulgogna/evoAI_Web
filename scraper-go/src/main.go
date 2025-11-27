package main

import (
	"fmt"
	"os"
	"scraper/src/customTypes"
	"scraper/src/extra"
	"scraper/src/urlManager"
	"strings"
	"sync"
)

var (
	Output []customTypes.Page = []customTypes.Page{}
	OutMutex sync.Mutex
)

var Query string = "what is duckduckgo"

func main() {
	if len(os.Args) > 1 {
		Query = strings.Join(os.Args[1:], " ")
	}

	search(Query, 20, 5)
}

func search(query string, totalUrls int, maxUrlsPerPages int) {

	performDdgSearch(query, 10)

	for range totalUrls {
		
		link := urlManager.GetUrl()

		if link == nil {
			break
		}

		fmt.Printf("url: %s, priority: %d\n", link.Url, link.Priority)

		ScrapeWg.Add(1)
		go func(l *customTypes.StoreUrl) {
			defer ScrapeWg.Done()
			scrape(l, totalUrls, maxUrlsPerPages)
		}(link)

	}
	ScrapeWg.Wait()

	totalPagesFound := len(Output)

	fmt.Printf("saving %d file/s\n",totalPagesFound)

	for i, page := range Output {
		extra.DisplayProgress(i, totalPagesFound, false)
		extra.WritePageToFile(page)
	}
	extra.DisplayProgress(totalPagesFound, totalPagesFound, true)
}