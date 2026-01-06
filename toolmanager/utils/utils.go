package utils

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pahulgogna/evoAI_Web/toolmanager/config"
)

func Init() {
	if _, err := os.Stat(config.ToolFile); os.IsNotExist(err) {
		f, err := os.Create(config.ToolFile)
		if err != nil {
			fmt.Println("error creating tools file: ", err)
		}

		f.WriteString("tools: ")

	}
}

func SendError(status int, data string, c *gin.Context) {
	c.AbortWithError(status, fmt.Errorf("{data:%s}", data))
}
