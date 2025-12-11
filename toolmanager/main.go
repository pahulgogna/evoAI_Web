package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/gin-gonic/gin"
	"github.com/pahulgogna/evoAI_Web/toolmanager/config"
	"github.com/pahulgogna/evoAI_Web/toolmanager/customTypes"
	"github.com/pahulgogna/evoAI_Web/toolmanager/utils"
)

func main() {

	utils.Init()
	config.Init()

	router := gin.Default()

	router.GET("/ping", ping)
	router.GET("/tools", getTools)
	router.POST("/tools/create", createTool)

	// TODO
	// router.DELETE("/tools/*name", deleteTool)
	// router.PATCH("/tools/*name", updateTool)

	router.Run(fmt.Sprintf("0.0.0.0:%s", "8080"))
}

func ping(c *gin.Context) {
	c.Writer.Write([]byte("pong"))
	c.Writer.Flush()
}

func getTools(c *gin.Context) {

	data, err := os.ReadFile(config.ToolFile)
	if err != nil {
		panic(fmt.Sprintf("could not open tools file. \n%s", err))
	}

	var store customtypes.Store
	if err := yaml.Unmarshal(data, &store); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(500)
		return
	}

	jsonBytes, err := json.Marshal(store)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(200, gin.H{
		"data": string(jsonBytes),
	})
}

func createTool(c *gin.Context) {

	var req customtypes.CreateRequestSchema
	if err := c.BindJSON(&req); err != nil {
		fmt.Printf("Error: could not parse the request body: %v\n", err)
		c.AbortWithStatus(400)
		return
	}

	content, err := os.ReadFile(config.ToolFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	var yamlData customtypes.Store
	if len(content) != 0 {
		if err := yaml.Unmarshal(content, &yamlData); err != nil {
			return
		}
	}
	if yamlData.Tools == nil {
		yamlData.Tools = make(map[string]customtypes.Snippet)
	}

	if _, ok := yamlData.Tools[req.Name]; ok {
		c.JSON(400, gin.H{
			"data": "tool name already exists",
		})
		return
	}
	yamlData.Tools[req.Name] = req.Tool

	newYAML, err := yaml.Marshal(yamlData)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.WriteFile(config.ToolFile, newYAML, 0644)
}


// func deleteTool(c *gin.Context) {

// }

// func updateTool(c *gin.Context) {

// }
