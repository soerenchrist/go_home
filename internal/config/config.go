package config

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("config")

var config *viper.Viper

func Init(env string) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../../config/")
	config.AddConfigPath("internal/config/")

	err = config.ReadInConfig()

	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	config.SetConfigName("secrets")
	err = config.MergeInConfig()
	if err != nil {
		log.Warning("No secrets file found. Proceeding...")
	}
}

func GetConfig() *viper.Viper {
	return config
}
