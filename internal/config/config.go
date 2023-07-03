package config

import (
	"log"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Init(env string) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")

	err = config.ReadInConfig()

	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	config.SetConfigName("secrets")
	err = config.MergeInConfig()
	if err != nil {
		log.Println("No secrets file found. Proceeding...")
	}
}

func GetConfig() *viper.Viper {
	return config
}