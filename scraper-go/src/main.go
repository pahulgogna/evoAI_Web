package main

import (
	"fmt"
	"scraper/src/customTypes"
	"scraper/src/extra"
	"scraper/src/urlManager"
	"sync"
	"os"
)

var (
	Output []customTypes.Page = []customTypes.Page{}
	OutMutex sync.Mutex
)

func main() {

	query := "what is duckduckgo"

	if len(os.Args) > 1 {
		query = fmt.Sprint(os.Args[1:])
	}

	search(query, 20, 10)
}

func search(query string, totalUrls int, maxUrlsPerPages int) {

	performDdgSearch(query, 10)

	for i := 0; i < maxUrlsPerPages; i++ {
		
		link := urlManager.GetUrl()

		if link == "" {
			break
		}

		ScrapeWg.Add(1)
		go func(l string) {
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