package config

import "flag"

// app configuration
var Values struct {
	ListenAddr string
	BasePath   string
}

func Init() {
	flag.StringVar(&Values.ListenAddr, "a", ":8080", "address and port to listen on")
	flag.StringVar(&Values.BasePath, "b", "http://localhost:8080", "base path")
	flag.Parse()
}
