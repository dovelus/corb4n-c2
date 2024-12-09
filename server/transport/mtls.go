package transport

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.RFC3339,
})

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

// logRequest is a middleware that logs the details of each request
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger.Infof("Received request: %s %s from %s", req.Method, req.URL.Path, req.RemoteAddr)
		handler.ServeHTTP(w, req)
	})
}

// StartServer starts the HTTPS server with the given TLS configuration
func StartServer(addr string, cfg *tls.Config) {
	srv := &http.Server{
		Addr: addr,
		Handler: logRequest(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			w.Write([]byte("This is an example server.\n"))
		})),
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	logger.Fatal(srv.ListenAndServeTLS("", ""))
}
