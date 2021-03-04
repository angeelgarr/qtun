package client

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net"

	qConn "github.com/marten-seemann/quic-conn"
	"github.com/net-byte/qtun/util"
)

//Start the client
func Start(config util.Config) {
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

func handleConn(fromConn net.Conn, config util.Config) {
	tlsConf, err := getTLSConfig(config)
	if err != nil {
		panic(err)
	}
	toConn, err := qConn.Dial(config.To, tlsConf)
	if err != nil {
		log.Printf("quic conn error:%s", err)
		return
	}
	go copy(toConn, fromConn)
	go copy(fromConn, toConn)
}

func copy(destination io.WriteCloser, source io.ReadCloser) {
	if destination == nil || source == nil {
		return
	}
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func getTLSConfig(config util.Config) (*tls.Config, error) {
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
