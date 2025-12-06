package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)


func main() {

	router := gin.Default()

	router.GET("/ping", ping)

	router.Run(fmt.Sprintf("0.0.0.0:%s", "8080"))

}

func ping(c *gin.Context) {
	c.Writer.Write([]byte("pong"))
	c.Writer.Flush()
}

func getTools(c *gin.Context) {

}

func createTool(c *gin.Context) {
	
}

func deleteTool(c *gin.Context) {

}

func updateTool(c *gin.Context) {
	
}
