package client

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/lucas-clemente/quic-go"
	"github.com/net-byte/qtun/common"
	"github.com/net-byte/qtun/config"
)

var _sessionMap sync.Map

//Start udp proxy client
func StartUDP(config config.Config) {
	log.Printf("udp proxy from %s to %s", config.From, config.To)
	tlsConf, err := common.GetClientTLSConfig(config)
	if err != nil {
		log.Panic(err)
	}
	localAddr, err := net.ResolveUDPAddr("udp", config.From)
	if err != nil {
		log.Fatalln("failed to get udp socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed to listen on udp socket:", err)
	}
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		n, cliAddr, err := conn.ReadFromUDP(buf)
		if err != nil || n == 0 {
			continue
		}
		var _session quic.Session
		var _err error
		if v, ok := _sessionMap.Load(cliAddr.String()); ok {
			_session = v.(quic.Session)
		} else {
			_session, _err = quic.DialAddr(config.To, tlsConf, nil)
			if _session == nil || _err != nil {
				log.Println(_err)
				continue
			}
			_sessionMap.Store(cliAddr.String(), _session)
		}
		stream, err := _session.OpenStreamSync(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		go toClient(_session, stream, conn, cliAddr)
		stream.Write(buf[:n])
	}
}

func toClient(session quic.Session, stream quic.Stream, conn *net.UDPConn, cliAddr *net.UDPAddr) {
	defer stream.Close()
	defer session.CloseWithError(0, "bye")
	buf := make([]byte, 4096)
	for {
		n, err := stream.Read(buf)
		if n == 0 || err != nil {
			break
		}
		conn.WriteToUDP(buf[0:n], cliAddr)
	}
	_sessionMap.Delete(cliAddr.String())
}
