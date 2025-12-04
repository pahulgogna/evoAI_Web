package searching

import (
	"scraper/src/customTypes"
	"scraper/src/extra"
	"scraper/src/urlManager"
	"scraper/src/utils"
)


func Search(query string, totalResults int, dnsAddress string) string {

	performDdgSearch(query, dnsAddress)

	callsCount := 0

	for totalResults > len(utils.Output) {
		
		link := urlManager.GetUrl()


		if link == nil {
			break
		}

		utils.ScrapeWg.Add(1)
		callsCount += 1
		go func(l *customTypes.StoreUrl) {
			defer utils.ScrapeWg.Done()
			utils.Scrape(l, query, dnsAddress)
			callsCount -= 1
		}(link)
		
		if callsCount == totalResults {
			utils.ScrapeWg.Wait()
		}
	}
	utils.ScrapeWg.Wait()

	return extra.GetJSON(utils.Output)
}