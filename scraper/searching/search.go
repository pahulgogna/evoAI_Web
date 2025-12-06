package searching

import (
	"github.com/pahulgogna/evoAI_Web/scraper/src/customTypes"
	"github.com/pahulgogna/evoAI_Web/scraper/src/extra"
	"github.com/pahulgogna/evoAI_Web/scraper/src/global"
	"github.com/pahulgogna/evoAI_Web/scraper/src/utils"
)

func Search(query string, totalResults int, scraper *global.ScraperSession) string {

	performDdgSearch(query, scraper)

	callsCount := 0

	for totalResults > len(scraper.Output) {

		link := scraper.Queue.GetUrl()

		if link == nil {
			break
		}

		scraper.Wg.Add(1)

		scraper.Mutex.Lock()
		callsCount += 1
		scraper.Mutex.Unlock()

		go func(l *customTypes.StoreUrl) {
			defer scraper.Wg.Done()
			utils.Scrape(l, query, scraper)

			scraper.Mutex.Lock()
			callsCount -= 1
			scraper.Mutex.Unlock()

		}(link)

		if callsCount == totalResults {
			scraper.Wg.Wait()
		}
	}
	scraper.Wg.Wait()

	return extra.GetJSON(scraper.Output)
}
