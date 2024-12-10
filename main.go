package main

import (
	"log"

	"github.com/dovelus/corb4n-c2/server/transport"
)



func main() {
	// Load server certificate and key
	serverCert, err := transport.LoadServerCertificate("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("failed to load server certificate and key: %v", err)
	}

	// Load client CA certificate
	clientCAPool, err := transport.LoadClientCACertificate("certs/client.crt")
	if err != nil {
		log.Fatalf("failed to read client CA certificate: %v", err)
	}

	// Create TLS configuration
	cfg := transport.CreateTLSConfig(serverCert, clientCAPool)

	// Start HTTPS server with TLS
	transport.StartServer(":8443", cfg)
}
