package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Beacon struct {
	BeaconID     string    `json:"beaconID"`
	Hostname     string    `json:"hostname"`
	BeaconIntIP  string    `json:"beaconIntIP"`
	BeaconExtIP  string    `json:"beaconExtIP"`
	OS           string    `json:"os"`
	EndpointProt string    `json:"endpointProt"`
	ProcessID    int       `json:"processID"`
	ProcessUser  string    `json:"processUser"`
	LastUpdate   time.Time `json:"lastUpdate"`
}

type BeaconFile struct {
	FileName string `json:"FileName"`
	FileType string `json:"FileType"`
	Output   string `json:"Output"`
	BeaconID string `json:"BeaconID"`
}

var mockBeacons = []Beacon{
	{
		BeaconID:     "beacon1",
		Hostname:     "host1.example.com",
		BeaconIntIP:  "192.168.1.1",
		BeaconExtIP:  "203.0.113.1",
		OS:           "Windows 10",
		EndpointProt: "AV1",
		ProcessID:    1234,
		ProcessUser:  "user1",
		LastUpdate:   time.Now().Add(-30 * time.Minute),
	},
	{
		BeaconID:     "beacon2",
		Hostname:     "host2.example.com",
		BeaconIntIP:  "192.168.1.2",
		BeaconExtIP:  "203.0.113.2",
		OS:           "Ubuntu 20.04",
		EndpointProt: "AV2",
		ProcessID:    5678,
		ProcessUser:  "user2",
		LastUpdate:   time.Now().Add(-7 * time.Hour),
	},
	{
		BeaconID:     "beacon3",
		Hostname:     "host3.example.com",
		BeaconIntIP:  "192.168.1.3",
		BeaconExtIP:  "203.0.113.3",
		OS:           "macOS 11",
		EndpointProt: "AV3",
		ProcessID:    9012,
		ProcessUser:  "user3",
		LastUpdate:   time.Now(),
	},
}

//go:embed frontendBuild/*
var frontendBuild embed.FS

func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// Ensure common MIME types are registered
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".html", "text/html")

	http.HandleFunc("/", serveFile)
	http.HandleFunc("/api/beacons", handleBeacons)
	http.HandleFunc("/api/beaconFiles", handleBeaconFiles)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.HandlerFunc(serveAsset)))

	fmt.Printf("Server is running on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func handleBeacons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockBeacons)
}

func handleBeaconFiles(w http.ResponseWriter, r *http.Request) {
	beaconID := r.URL.Query().Get("BeaconID")
	if beaconID == "" {
		http.Error(w, "BeaconID is required", http.StatusBadRequest)
		return
	}

	files, err := getFilesForBeacon(beaconID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Pretty-print the files slice for better readability
	filesJSON, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(filesJSON))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func getFilesForBeacon(beaconID string) ([]BeaconFile, error) {
	var files []BeaconFile
	basePath := filepath.Join("beaconFiles", beaconID)

	err := filepath.Walk(basePath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			files = append(files, BeaconFile{
				FileName: info.Name(),
				FileType: mime.TypeByExtension(filepath.Ext(info.Name())),
				Output:   string(content),
				BeaconID: beaconID,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	requestPath := filepath.Clean(r.URL.Path)
	if requestPath == "/" {
		requestPath = "/index.html"
	}

	file := path.Join("frontendBuild", requestPath)
	_, err := frontendBuild.Open(file)
	if os.IsNotExist(err) {
		// Serve index.html for all routes
		file = "frontendBuild/index.html"
	}

	// Set the correct content type
	ext := filepath.Ext(file)
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}

	data, err := frontendBuild.ReadFile(file)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Write(data)
}

func serveAsset(w http.ResponseWriter, r *http.Request) {
	requestPath := filepath.Clean(r.URL.Path)
	file := path.Join("frontendBuild/assets", requestPath)

	// Set the correct content type
	ext := filepath.Ext(file)
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}

	data, err := frontendBuild.ReadFile(file)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Write(data)
}
