package config

import (
	"flag"
)

type Config struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	DBName   string `mapstructure:"dbname"`
	Port     int    `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

func Init() *Config {
	var cfg Config

	var (
		username string
		password string
		dbname   string
		port     int
		sslmode  string
		host     string
	)

	flag.StringVar(&username, "username", "postgres", "if required username is not postgres, then use this flag")
	flag.StringVar(&password, "password", "postgres", "if required password is not postgres, then use this flag")
	flag.StringVar(&dbname, "dbname", "postgres", "if required database is not postgres, then use this flag")
	flag.IntVar(&port, "port", 5432, "if required port is not 5432, then use this flag")
	flag.StringVar(&sslmode, "sslmode", "disable", "if required sslmode is not 'disabled', then use this flag")
	flag.StringVar(&host, "host", "localhost", "if required host is not localhost, then use this flag")

	flag.Parse()

	cfg.User = username
	cfg.Password = password
	cfg.DBName = dbname
	cfg.Port = port
	cfg.SSLMode = sslmode
	cfg.Host = host
	return &cfg
}
