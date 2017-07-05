package server

import (
	"crypto/tls"
	"crypto/x509"

	"fmt"

	"github.com/biolee/gRPC-REST/cert"
)

const (
	ip   = "0.0.0.0"
	port = 10000
)

var (
	keyPair  = &tls.Certificate{}
	certPool = x509.NewCertPool()
	addr     = fmt.Sprintf("%s:%d", ip, port)
)

func init() {
	var err error
	*keyPair, err = tls.X509KeyPair(cert.Cert, cert.Key)
	if err != nil {
		panic(err)
	}

	ok := certPool.AppendCertsFromPEM(cert.Cert)
	if !ok {
		panic("bad certs")
	}
}
