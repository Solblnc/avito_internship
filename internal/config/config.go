package config

import (
	"avito_internship/internal/database"
	"avito_internship/internal/server"
	"github.com/spf13/viper"
)

type Env struct {
	Config database.Config
	Port   server.Server
}

func Init() (*Env, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Env

	if err := viper.UnmarshalKey("pg_config", &cfg.Config); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("server_config", &cfg.Port); err != nil {
		return nil, err
	}
	return &cfg, nil
}
