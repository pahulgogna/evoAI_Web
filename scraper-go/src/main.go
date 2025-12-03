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

	search(Query, 5)
}

func search(query string, totalResults int) {

	performDdgSearch(query)

	callsCount := 0

	for totalResults > len(Output) {
		
		link := urlManager.GetUrl()

		if link == nil {
			break
		}

		ScrapeWg.Add(1)
		callsCount += 1
		go func(l *customTypes.StoreUrl) {
			defer ScrapeWg.Done()
			scrape(l)
			callsCount -= 1
		}(link)
		
		if callsCount == totalResults {
			ScrapeWg.Wait()
		}
	}
	ScrapeWg.Wait()

	jsonData:= extra.GetJSON("out.json", Output)

	fmt.Println(jsonData)
}