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

// Start tcp proxy server
func StartTCP(config config.Config) {
	log.Printf("Proxy from %s to %s", config.From, config.To)
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
		stream, err := session.AcceptStream(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		go handleTCP(stream, config)
	}
}

func handleTCP(stream quic.Stream, config config.Config) {
	conn, err := net.DialTimeout("tcp", config.To, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		log.Println(err)
		return
	}

	go common.Copy(conn, stream)
	go common.Copy(stream, conn)
}
