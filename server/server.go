package server

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/soerenchrist/go_home/background"
	"github.com/soerenchrist/go_home/config"
	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/models"
	"github.com/soerenchrist/go_home/mqtt"
	"github.com/soerenchrist/go_home/rules/evaluation"
	"github.com/spf13/viper"
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
	outputBindings := make(chan models.SensorValue, 10)
	// TODO: Refactor outputbindings channel passing around
	addRulesEngine(database, outputBindings)
	addMqttBinding(config, database, outputBindings)

	r := NewRouter(database, outputBindings)

	port := config.GetString("server.port")
	host := config.GetString("server.host")

	addr := fmt.Sprintf("%s:%s", host, port)
	r.Run(addr)
}

func addRulesEngine(database db.Database, outputBindings chan models.SensorValue) {
	rulesEngine := evaluation.NewRulesEngine(database)

	go rulesEngine.ListenForValues(outputBindings)
	go background.PollSensorValues(database, outputBindings)
}

func addMqttBinding(config *viper.Viper, database db.Database, outputBindings chan models.SensorValue) {
	enabled := config.GetBool("mqtt.enabled")
	if !enabled {
		return
	}

	host := config.GetString("mqtt.host")
	port := config.GetInt("mqtt.port")
	username := config.GetString("mqtt.username")
	password := config.GetString("mqtt.password")
	clientId := config.GetString("mqtt.clientId")

	options := mqtt.MqttConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		ClientId: clientId,
	}

	publishChannel := make(chan mqtt.Message, 10)

	err := mqtt.AddMqttBinding(options, publishChannel, database, outputBindings)
	if err != nil {
		log.Println("Failed to add MQTT binding: ", err)
	}
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
