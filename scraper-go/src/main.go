package main

import (
	"fmt"
	"sync"
)

var (
	seen map[string]int8 = map[string]int8{}
	MapMutex sync.Mutex
)

var (
	output []Page = []Page{}
	OutMutex sync.Mutex
)

func main() {
	search("", 2, 10)
}

func search(query string, level int, maxUrlsPerPages int) {

	links := performDdgSearch(query)[:10]

	if links == nil {
		fmt.Println("could not perform search")
		return
	}

	for _, link := range links {
		
		ScrapeWg.Add(1)
		go func(l string, lev int) {
			defer ScrapeWg.Done()
			scrape(l, lev, maxUrlsPerPages)
		}(link, level)

	}
	ScrapeWg.Wait()

	totalPagesFound := len(output)

	fmt.Printf("saving %d file/s\n",totalPagesFound)

	for i, page := range output {
		displayProgress(i, totalPagesFound, false)
		WritePageToFile(page)
	}
	displayProgress(totalPagesFound, totalPagesFound, true)
}