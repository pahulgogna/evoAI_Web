package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pahulgogna/evoAI_Web/scraper/config"
	"github.com/pahulgogna/evoAI_Web/scraper/customTypes"
	"github.com/pahulgogna/evoAI_Web/scraper/global"
	"github.com/pahulgogna/evoAI_Web/scraper/searching"
	"github.com/pahulgogna/evoAI_Web/scraper/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	config.Init()

	router := gin.Default()
	router.POST("/search", getSearchResults)

	router.Run(fmt.Sprintf("0.0.0.0:%s", config.PORT))
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

	scraper := global.NewScraperSession()
	scraper.Client = &http.Client{
		Timeout:   10 * time.Second,
		Transport: utils.GetTransportForRequest(req.DnsAddress),
	}

	c.IndentedJSON(200, searching.Search(req.Query, req.RequiredResults, scraper))

	scraper.Queue.ClearQueue()
	scraper.ClearQueue()
}
