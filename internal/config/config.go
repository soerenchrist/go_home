package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

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
		log.Fatal().Err(err).Msg("Error reading config file")
	}

	config.SetConfigName("secrets")
	err = config.MergeInConfig()
	if err != nil {
		log.Warn().Msg("No secrets file found. Proceeding...")
	}
}

func GetConfig() *viper.Viper {
	return config
}
