package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/soerenchrist/go_home/internal/background"
	"github.com/soerenchrist/go_home/internal/config"
	"github.com/soerenchrist/go_home/internal/db"
	"github.com/soerenchrist/go_home/internal/mqtt"
	"github.com/soerenchrist/go_home/internal/rules/evaluation"
	"github.com/soerenchrist/go_home/pkg/output"
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
	setupLogging(config)
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
	outputBindings := output.NewManager()
	go background.CleanupExpiredSensorValues(sqlite)
	addRulesEngine(database, outputBindings)

	runHomeServer(config, database, outputBindings)
	runMqttBridge(config)

	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start")
	}
}

func setupLogging(config *viper.Viper) {
	loglevel := config.GetString("logging.level")
	switch loglevel {
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func runHomeServer(config *viper.Viper, database db.Database, outputBindings *output.OutputBindingsManager) {
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
		log.Info().Str("address", addr).Msg("Starting home server")
		return server.ListenAndServe()
	})
}

func runMqttBridge(config *viper.Viper) {
	router, err := addMqttBridge(config)
	if err != nil {
		log.Warn().Err(err)
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
		log.Info().Str("address", addr).Msg("Starting MQTT bridge")
		return server.ListenAndServe()
	})
}

func addRulesEngine(database db.Database, outputBindings *output.OutputBindingsManager) {
	rulesEngine := evaluation.NewRulesEngine(database)

	rulesOutput := output.NewChannelOutput()
	outputBindings.Register(rulesOutput)

	go rulesEngine.ListenForValues(rulesOutput)
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
