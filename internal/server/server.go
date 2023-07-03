package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/go_home/internal/background"
	"github.com/soerenchrist/go_home/internal/config"
	"github.com/soerenchrist/go_home/internal/db"
	"github.com/soerenchrist/go_home/internal/mqtt"
	"github.com/soerenchrist/go_home/internal/rules/evaluation"
	"github.com/soerenchrist/go_home/internal/value"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	g errgroup.Group
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
	outputBindings := value.NewOutputBindings()
	addRulesEngine(database, outputBindings)

	runHomeServer(config, database, outputBindings)
	runMqttBridge(config)

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func runHomeServer(config *viper.Viper, database db.Database, outputBindings *value.OutputBindings) {
	r := NewRouter(database, outputBindings)

	port := config.GetString("server.port")
	host := config.GetString("server.host")

	addr := fmt.Sprintf("%s:%s", host, port)

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		log.Printf("Starting home server on %s\n", addr)
		return server.ListenAndServe()
	})
}

func runMqttBridge(config *viper.Viper) {
	router, err := addMqttBridge(config)
	if err != nil {
		log.Println(err)
		return
	}

	port := config.GetString("mqtt.bridge.port")
	host := config.GetString("mqtt.bridge.host")

	addr := fmt.Sprintf("%s:%s", host, port)

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		log.Printf("Starting MQTT bridge on %s\n", addr)
		return server.ListenAndServe()
	})
}

func addRulesEngine(database db.Database, outputBindings *value.OutputBindings) {
	rulesEngine := evaluation.NewRulesEngine(database)

	sensorsChan := make(chan value.SensorValue)
	outputBindings.Register(sensorsChan)

	go rulesEngine.ListenForValues(sensorsChan)
	go background.PollSensorValues(database, outputBindings)
}

func addMqttBridge(config *viper.Viper) (*gin.Engine, error) {
	enabled := config.GetBool("mqtt.enabled")
	if !enabled {
		return nil, fmt.Errorf("MQTT bridge is not enabled")
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

	err := mqtt.ConnectToBroker(options, publishChannel, config)
	if err != nil {
		return nil, fmt.Errorf("failed to add MQTT binding: %v", err)
	}

	router := mqtt.NewMqttRouter(publishChannel)
	return router, nil
}

func openDatabase(path string) *gorm.DB {
	db := sqlite.Open(path)
	gdb, err := gorm.Open(db, &gorm.Config{})
	if err != nil {
		panic("failed to open database")
	}

	return gdb
}
