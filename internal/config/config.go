package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	DBName   string `mapstructure:"dbname"`
	Port     int    `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

func Init(path string) (*Config, error) {
	op := "config.Init()"
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%s: failed to read in config: %w", op, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("%s: failed to unmarshal config: %w", op, err)
	}

	return &cfg, nil
}
