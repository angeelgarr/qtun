//go:build linux || darwin || windows

package main

import (
	"flag"

	"github.com/net-byte/qtun/client"
	"github.com/net-byte/qtun/config"
	"github.com/net-byte/qtun/server"
)

func main() {
	config := config.Config{}
	flag.StringVar(&config.From, "from", ":1987", "from address")
	flag.StringVar(&config.To, "to", ":1080", "to address")
	flag.StringVar(&config.ClientKey, "ck", "../certs/client.key", "client key file path")
	flag.StringVar(&config.ClientPem, "cp", "../certs/client.pem", "client pem file path")
	flag.StringVar(&config.ServerKey, "sk", "../certs/server.key", "server key file path")
	flag.StringVar(&config.ServerPem, "sp", "../certs/server.pem", "server pem file path")
	flag.IntVar(&config.Timeout, "t", 30, "connection timeout in seconds")
	flag.BoolVar(&config.ServerMode, "S", false, "server mode")

	flag.Parse()
	if config.ServerMode {
		server.Start(config)
	} else {
		client.Start(config)
	}
}
