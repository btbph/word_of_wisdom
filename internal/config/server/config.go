package server

import (
	"github.com/btbph/word_of_wisdom/internal/config/server/challenge"
	"github.com/btbph/word_of_wisdom/internal/config/server/server"
	"github.com/spf13/viper"
)

type Config struct {
	Server    server.Config
	Challenge challenge.Config
}

func New() (*Config, error) {
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigName("server")
	vp.AddConfigPath("config")
	vp.AutomaticEnv()
	vp.SetEnvPrefix("APP")

	var config *Config
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := vp.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
