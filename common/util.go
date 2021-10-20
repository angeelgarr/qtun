package common

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"

	"github.com/net-byte/qtun/config"
)

func Copy(destination io.WriteCloser, source io.ReadCloser) {
	if destination == nil || source == nil {
		return
	}
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func GetClientTLSConfig(config config.Config) (*tls.Config, error) {
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

func GetServerTLSConfig(config config.Config) (*tls.Config, error) {
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
