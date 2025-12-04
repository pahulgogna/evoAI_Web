package main

import (
	"fmt"
	"scraper/src/config"
	"scraper/src/customTypes"
	"scraper/src/searching"
	"scraper/src/urlManager"
	"scraper/src/utils"

	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    router.POST("/search", getSearchResults)

    router.Run(fmt.Sprintf("0.0.0.0:%d", config.PORT))
}

func getSearchResults(c *gin.Context) {

	var req customTypes.SearchRequest

	if err := c.BindJSON(&req); err != nil {
		fmt.Printf("Error: could not parse the request body: %v\n", err)
		c.AbortWithStatus(400)
		return
	}
	if req.RequiredResults <= 0 || req.Query == "" {
		c.AbortWithStatus(400)
		return
	}
	fmt.Println(req)
	if req.DnsAddress == "" {
		req.DnsAddress = config.DefaultDNSAddress
	}

	fmt.Printf("/search : query-> %s, results-> %d \n", req.Query, req.RequiredResults)

	c.IndentedJSON(200, searching.Search(req.Query, req.RequiredResults, req.DnsAddress))
	urlManager.ClearQueue()
	utils.OutMutex.Lock()
	utils.Output = []customTypes.Page{}
	utils.OutMutex.Unlock()
}
