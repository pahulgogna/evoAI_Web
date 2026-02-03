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