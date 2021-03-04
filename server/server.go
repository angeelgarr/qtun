package server

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	qConn "github.com/marten-seemann/quic-conn"
	"github.com/net-byte/qtun/util"
)

// Start the server
func Start(config util.Config) {
	log.Printf("Proxy from %s to %s", config.From, config.To)
	tlsConf, err := getTLSConfig(config)
	if err != nil {
		panic(err)
	}

	l, err := qConn.Listen("udp", config.From, tlsConf)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		handleConn(conn, config)
	}
}

func handleConn(fromConn net.Conn, config util.Config) {
	// Connect the dest server
	toConn, err := net.DialTimeout("tcp", config.To, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		log.Println(err)
		return
	}

	// Copy data
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
