package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config - app configuration
type Config struct {
	ListenAddr      string `json:"server_address"`
	BasePath        string `json:"base_url"`
	LogLevel        string `json:"log_level"`
	FileStoragePath string `json:"file_storage_path"`
	StoreInFile     bool   `json:"store_in_file"`
	DSN             string `json:"database_dsn"`
	StoreInDB       bool   `json:"store_in_db"`
	DelBufSize      int    `json:"del_buf_size"`
	Secure          bool   `json:"enable_https"`
	configFile      string
}

// New - constructor
func New() *Config {
	res := Config{}

	flag.StringVar(
		&res.ListenAddr, "a", ":8080",
		"address and port to listen on")
	flag.StringVar(&res.BasePath, "b", "http://localhost:8080", "base path")
	flag.StringVar(&res.LogLevel, "l", "info", "log level")
	flag.StringVar(&res.FileStoragePath, "f", "/tmp/short-url-db.json",
		"filepath to simple database")
	flag.StringVar(&res.DSN, "d",
		"",
		"data source name")
	flag.IntVar(&res.DelBufSize, "delbuf", 10,
		"how many delete requests to buffer")
	flag.BoolVar(&res.Secure, "s", false,
		"enable HTTPS")
	flag.StringVar(&res.configFile, "c", "", "configuration file")
	flag.StringVar(&res.configFile, "config", "", "configuration file")
	return &res
}

// MustParse - parse configuration or panic
func (c *Config) MustParse() {
	flag.Parse()

	if val := os.Getenv("CONFIG"); len(val) > 0 {
		c.configFile = val
	}

	c.readConfigFile()

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

	if val := os.Getenv("DATABASE_DSN"); val != "" {
		c.DSN = val
	}

	if val := os.Getenv("DELETE_BUF_SIZE"); len(val) > 0 {
		var err error
		c.DelBufSize, err = strconv.Atoi(val)
		if err != nil {
			panic(fmt.Errorf("failed to parse delete buf size: %w", err))
		}
	}

	if val := os.Getenv("ENABLE_HTTPS"); len(val) > 0 {
		c.Secure = true
	}

	if c.DSN != "" {
		c.StoreInDB = true
	}
}

func (c *Config) readConfigFile() {
	if len(c.configFile) == 0 {
		return
	}

	b, err := os.ReadFile(c.configFile)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(b, c)
}
