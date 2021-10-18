package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/net-byte/qtun/config"
)

// Start the proxy server
func Start(config config.Config) {
	log.Printf("Proxy from %s to %s", config.From, config.To)
	tlsConf, err := getTLSConfig(config)
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
		go handleConn(stream, config)
	}
}

func handleConn(stream quic.Stream, config config.Config) {
	conn, err := net.DialTimeout("tcp", config.To, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		log.Println(err)
		return
	}

	go copy(conn, stream)
	go copy(stream, conn)
}

func copy(destination io.WriteCloser, source io.ReadCloser) {
	if destination == nil || source == nil {
		return
	}
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func getTLSConfig(config config.Config) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(config.ServerPem, config.ServerKey)
	if err != nil {
		return nil, err
	}

	certBytes, err := ioutil.ReadFile(config.ClientPem)
	if err != nil {
		return nil, err
	}

	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse client certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
		NextProtos:   []string{"qtun/1.0"},
	}, nil
}
