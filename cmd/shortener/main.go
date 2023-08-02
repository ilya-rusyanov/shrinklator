package main

import (
	"github.com/ilya-rusyanov/shrinklator/internal/config"
	"github.com/ilya-rusyanov/shrinklator/internal/server"
)

func main() {
	config.Init()
	server.Run()
}
