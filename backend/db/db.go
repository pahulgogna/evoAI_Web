package db

import (
	"database/sql"
	"fmt"

	// "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pahulgogna/evoAI_Web/backend/config"
)

func NewPostgresStorage() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Envs.DBAddress,
		config.Envs.DBPort,
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// var mySqlConfig = mysql.Config{
// 	User:                 config.Envs.DBUser,
// 	Passwd:               config.Envs.DBPassword,
// 	Addr:                 fmt.Sprintf("%s:%s", config.Envs.DBAddress, config.Envs.DBPort),
// 	DBName:               config.Envs.DBName,
// 	Net:                  "tcp",
// 	AllowNativePasswords: true,
// 	ParseTime:            true,
// }

// func NewMySqlStorage() (*sql.DB, error) {

// 	db, err := sql.Open("mysql", mySqlConfig.FormatDSN())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	return db, nil
// }
