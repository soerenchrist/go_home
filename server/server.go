package server

import (
	"database/sql"
	"fmt"

	"github.com/soerenchrist/mini_home/background"
	"github.com/soerenchrist/mini_home/config"
	"github.com/soerenchrist/mini_home/db"
)

func Init() {

	config := config.GetConfig()
	databasePath := config.GetString("database.path")

	sqlite := openDatabase(databasePath)

	database, err := db.NewDevicesDatabase(sqlite)

	seed := config.GetBool("database.seed")
	if seed {
		database.SeedDatabase()
	}

	if err != nil {
		panic(err)
	}

	go background.PollSensorValues(database)
	r := NewRouter(database)

	port := config.GetString("server.port")
	host := config.GetString("server.host")

	addr := fmt.Sprintf("%s:%s", host, port)
	r.Run(addr)
}

func openDatabase(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	return db
}
