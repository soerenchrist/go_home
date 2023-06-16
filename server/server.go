package server

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/soerenchrist/go_home/background"
	"github.com/soerenchrist/go_home/config"
	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/models"
	"github.com/soerenchrist/go_home/rules/evaluation"
)

func Init() {
	config := config.GetConfig()
	databasePath := config.GetString("database.path")

	sqlite := openDatabase(databasePath)

	database, err := db.NewDevicesDatabase(sqlite)
	if err != nil {
		panic(err)
	}

	seed := config.GetBool("database.seed")
	if seed {
		database.SeedDatabase()
	}

	sensorsChan := make(chan models.Sensor)
	rulesEngine := evaluation.NewRulesEngine(database)

	go rulesEngine.ListenForValues(sensorsChan)
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
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		panic(err)
	}

	return db
}
