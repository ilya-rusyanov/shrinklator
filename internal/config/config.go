package config

import (
	"flag"
	"os"
)

type Config struct {
	ListenAddr string
	BasePath   string
}

func New() *Config {
	res := Config{}
	flag.StringVar(
		&res.ListenAddr, "a", ":8080",
		"address and port to listen on")
	flag.StringVar(&res.BasePath, "b", "http://localhost:8080", "base path")
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
}
