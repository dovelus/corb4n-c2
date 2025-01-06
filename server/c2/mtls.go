package c2

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

// LoadServerCertificate loads the server certificate and key
func LoadServerCertificate(certFile, keyFile string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certFile, keyFile)
}

// LoadClientCACertificate loads the client CA certificate
func LoadClientCACertificate(caCertFile string) (*x509.CertPool, error) {
	clientCACert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}
	clientCAPool := x509.NewCertPool()
	clientCAPool.AppendCertsFromPEM(clientCACert)
	return clientCAPool, nil
}

// CreateTLSConfig creates a TLS configuration
func CreateTLSConfig(serverCert tls.Certificate, clientCAPool *x509.CertPool) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    clientCAPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
}
