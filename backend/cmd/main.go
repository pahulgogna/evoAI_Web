package main

import (
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pahulgogna/evoAI_Web/backend/cmd/api"
	"github.com/pahulgogna/evoAI_Web/backend/config"
	"github.com/pahulgogna/evoAI_Web/backend/db"
)

func main() {

	db, err := db.NewMySqlStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	xdb := sqlx.NewDb(db, "mysql")

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
