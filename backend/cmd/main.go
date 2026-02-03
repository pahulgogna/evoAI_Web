package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pahulgogna/evoAI_Web/backend/cmd/api"
	"github.com/pahulgogna/evoAI_Web/backend/config"
	"github.com/pahulgogna/evoAI_Web/backend/db"
)

func main() {

	db, err := db.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	xdb := sqlx.NewDb(db, "postgres")

	initDb(xdb)

	apiServer := api.NewServer(fmt.Sprintf(":%s", config.Envs.Port), xdb)

	if err := apiServer.Run(); err != nil {
		log.Panic(err)
	}
}

func initDb(database *sqlx.DB) {
	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}

	// database.MustExec(db.Schema)

	log.Println("connected to the database")
}


// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/ollama/ollama/api"
// )

// func main() {
// 	client, err := api.ClientFromEnvironment()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	messages := []api.Message{
// 		{
// 			Role:    "system",
// 			Content: "Provide very brief, concise responses",
// 		},
// 		{
// 			Role:    "user",
// 			Content: "Name some unusual animals",
// 		},
// 	}

// 	ctx := context.Background()
// 	req := &api.ChatRequest{
// 		Model:    "gemma3:4b",
// 		Messages: messages,
// 	}

// 	respFunc := func(resp api.ChatResponse) error {
// 		fmt.Print(resp.Message.Content)
// 		return nil
// 	}

// 	err = client.Chat(ctx, req, respFunc)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }