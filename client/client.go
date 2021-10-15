package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/lucas-clemente/quic-go"
	"github.com/net-byte/qtun/config"
)

//Start the proxy client
func Start(config config.Config) {
	log.Printf("Proxy from %s to %s", config.From, config.To)
	l, err := net.Listen("tcp", config.From)
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleConn(conn, config)
	}
}

func handleConn(conn net.Conn, config config.Config) {
	tlsConf, err := getTLSConfig(config)
	if err != nil {
		panic(err)
	}
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

	go copy(stream, conn)
	go copy(conn, stream)
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
	cert, err := tls.LoadX509KeyPair(config.ClientPem, config.ClientKey)
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
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		NextProtos:         []string{"qtun/1.0"},
	}, nil
}
