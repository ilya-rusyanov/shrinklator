package config

import (
	"flag"
	"os"
)

type Config struct {
	ListenAddr      string
	BasePath        string
	LogLevel        string
	FileStoragePath string
	StoreInFile     bool
}

func New() *Config {
	res := Config{}
	flag.StringVar(
		&res.ListenAddr, "a", ":8080",
		"address and port to listen on")
	flag.StringVar(&res.BasePath, "b", "http://localhost:8080", "base path")
	flag.StringVar(&res.LogLevel, "l", "info", "log level")
	flag.StringVar(&res.FileStoragePath, "f", "/tmp/short-url-db.json",
		"filepath to simple database")
	return &res
}

func (c *Config) Parse() {
	flag.Parse()

	if val, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.ListenAddr = val
	}

	if val, ok := os.LookupEnv("BASE_URL"); ok {
		c.BasePath = val
	}

	if val, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		c.FileStoragePath = val
	}

	if c.FileStoragePath != "" {
		c.StoreInFile = true
	}
}
