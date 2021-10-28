package server

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/net-byte/qtun/common"
	"github.com/net-byte/qtun/config"
)

// Start proxy server
func Start(config config.Config) {
	network := "tcp"
	if config.UDPMode {
		network = "udp"
	}
	log.Printf("%s proxy from %s to %s", network, config.From, config.To)
	tlsConf, err := common.GetServerTLSConfig(config)
	if err != nil {
		log.Panic(err)
	}
	l, err := quic.ListenAddr(config.From, tlsConf, nil)
	if err != nil {
		log.Panic(err)
	}
	for {
		session, err := l.Accept(context.Background())
		if err != nil {
			continue
		}
		go handleConn(network, session, config)
	}
}

func handleConn(network string, session quic.Session, config config.Config) {
	defer session.CloseWithError(0, "bye")
	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := net.DialTimeout(network, config.To, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		log.Println(err)
		return
	}

	go common.Copy(conn, stream)
	common.Copy(stream, conn)
}
