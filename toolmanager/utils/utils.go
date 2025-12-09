package utils

import (
	"fmt"
	"os"

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