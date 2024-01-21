package client

import (
	"github.com/btbph/word_of_wisdom/internal/config/client/client"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Client client.Config
}

func New() (*Config, error) {
	vp := viper.New()
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vp.SetConfigType("yaml")
	vp.SetConfigName("client")
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
