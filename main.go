package main

import (
	"flag"
	"qtun/client"
	"qtun/server"
	"qtun/util"
)

func main() {
	config := util.Config{}
	flag.StringVar(&config.From, "from", "127.0.0.1:1987", "From address")
	flag.StringVar(&config.To, "to", "127.0.0.1:1080", "To address")
	flag.StringVar(&config.ClientKey, "ck", "../certs/client.key", "Client key")
	flag.StringVar(&config.ClientPem, "cp", "../certs/client.pem", "Client pem")
	flag.StringVar(&config.ServerKey, "sk", "../certs/server.key", "Server key")
	flag.StringVar(&config.ServerPem, "sp", "../certs/server.pem", "Server pem")
	flag.IntVar(&config.Timeout, "t", 10, "Connection timeout of seconds")
	flag.BoolVar(&config.ServerMode, "s", false, "Server mode")

	flag.Parse()
	if config.ServerMode {
		server.Start(config)
	} else {
		client.Start(config)
	}
}
