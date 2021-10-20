package client

import (
	"context"
	"crypto/tls"
	"log"
	"net"

	"github.com/lucas-clemente/quic-go"
	"github.com/net-byte/qtun/common"
	"github.com/net-byte/qtun/config"
)

//Start tcp proxy client
func StartTCP(config config.Config) {
	log.Printf("Proxy from %s to %s", config.From, config.To)
	tlsConf, err := common.GetClientTLSConfig(config)
	if err != nil {
		log.Panic(err)
	}
	l, err := net.Listen("tcp", config.From)
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go handleTCP(conn, config, tlsConf)
	}
}

func handleTCP(conn net.Conn, config config.Config, tlsConf *tls.Config) {
	session, err := quic.DialAddr(config.To, tlsConf, nil)
	if err != nil {
		log.Println(err)
		return
	}
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	go common.Copy(stream, conn)
	go common.Copy(conn, stream)
}
