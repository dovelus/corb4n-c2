package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dovelus/corb4n-c2/server/db"
)

func main() {
	// Initialize the database
	db.InitDB()
	defer db.CloseDB()

	// Create sample implant data
	implant := &db.Implant_info{
		ID:          "implant1",
		Hostname:    "host1",
		IntIP:       "192.168.1.1",
		ExtIP:       "8.8.8.8",
		Os:          "Linux",
		ProcessID:   1234,
		ProcessUser: "user1",
		ProtName:    "None",
		LastCheckIn: time.Now().Unix(),
		Active:      true,
		KillDate:    0,
	}

	// Add implant to the database
	err := db.AddImplant(implant)
	if err != nil {
		fmt.Println("Error adding implant:", err)
		return
	}

	// Get all implants
	implants, err := db.GetAllImplants()
	if err != nil {
		fmt.Println("Error getting all implants:", err)
		return
	}
	fmt.Println("All Implants:", implants)

	// Get implant by ID
	implantByID, err := db.GetImplantByID("implant1")
	if err != nil {
		fmt.Println("Error getting implant by ID:", err)
		return
	}
	fmt.Println("Implant by ID:", implantByID)

	// Update implant status
	err = db.SetImplantStatus("implant1", false)
	if err != nil {
		fmt.Println("Error updating implant status:", err)
		return
	}

	// Get implant status
	status, err := db.GetImplantStatus("implant1")
	if err != nil {
		fmt.Println("Error getting implant status:", err)
		return
	}
	fmt.Println("Implant Status:", status)

	// Update implant kill date
	err = db.UpdateImplantKillDate("implant1")
	if err != nil {
		fmt.Println("Error updating implant kill date:", err)
		return
	}

	// Update implant check-in time
	err = db.UpdateImplantCheckin("implant1")
	if err != nil {
		fmt.Println("Error updating implant check-in time:", err)
		return
	}

	// Create sample task data
	task := &db.Implant_Task{
		TaskID:      "task1",
		ImplantID:   "implant1",
		TaskType:    1,
		TaskData:    []byte("Sample Task Data"),
		CreatedAt:   time.Now().Unix(),
		Completed:   false,
		CompletedAt: 0,
		TaskResult:  []byte(""),
	}

	// Add task to the database
	err = db.AddTask(task)
	if err != nil {
		fmt.Println("Error adding task:", err)
		return
	}

	// Get all tasks for an implant
	tasks, err := db.GetImplantTasks("implant1", false)
	if err != nil {
		fmt.Println("Error getting implant tasks:", err)
		return
	}
	fmt.Println("Implant Tasks:", tasks)

	// Complete a task
	err = db.CompleteTask("task1")
	if err != nil {
		fmt.Println("Error completing task:", err)
		return
	}

	// Get task by ID
	taskByID, err := db.GetTask("task1")
	if err != nil {
		fmt.Println("Error getting task by ID:", err)
		return
	}
	fmt.Println("Task by ID:", taskByID)

	// Remove a task
	// err = db.RemoveTask("task1")
	// if err != nil {
	// 	fmt.Println("Error removing task:", err)
	// 	return
	// }

	// Create sample listener data
	listener := &db.Listener_info{
		ListenerID: "listener1",
		Config:     []byte("Sample Config"),
		Host:       "localhost",
		Port:       8080,
		CreatedAt:  time.Now().Unix(),
		KillDate:   0,
	}

	// Add listener to the database
	err = db.AddListener(listener)
	if err != nil {
		fmt.Println("Error adding listener:", err)
		return
	}

	// Get all listeners
	listeners, err := db.GetAllListeners()
	if err != nil {
		fmt.Println("Error getting all listeners:", err)
		return
	}
	fmt.Println("All Listeners:", listeners)

	// Get listener by ID
	listenerByID, err := db.GetListenerByID("listener1")
	if err != nil {
		fmt.Println("Error getting listener by ID:", err)
		return
	}
	fmt.Println("Listener by ID:", listenerByID)

	// Update listener kill date
	err = db.UpdateListenerKillDate("listener1")
	if err != nil {
		fmt.Println("Error updating listener kill date:", err)
		return
	}

	// Remove a listener
	// err = db.RemoveListener("listener1")
	// if err != nil {
	// 	fmt.Println("Error removing listener:", err)
	// 	return
	// }

	// Create sample file data
	file := &db.File_info{
		ImplantID: "implant1",
		FileName:  "sample.txt",
		FileSize:  1024,
		FileType:  "text/plain",
		FilePath:  "files/sample.txt",
		CreatedAt: time.Now().Unix(),
	}

	// Ensure the files directory exists
	err = os.MkdirAll("files", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating files directory:", err)
		return
	}

	// Create a sample file
	fileContent := []byte("This is a sample file.")
	err = os.WriteFile(file.FilePath, fileContent, 0644)
	if err != nil {
		fmt.Println("Error creating sample file:", err)
		return
	}

	// Add file to the database
	err = db.AddFile(file)
	if err != nil {
		fmt.Println("Error adding file:", err)
		return
	}

	// Get all files
	files, err := db.GetAllFiles()
	if err != nil {
		fmt.Println("Error getting all files:", err)
		return
	}
	fmt.Println("All Files:", files)

	// Get file by implant ID
	filesByImplantID, err := db.GetFileByImplantID("implant1")
	if err != nil {
		fmt.Println("Error getting files by implant ID:", err)
		return
	}
	fmt.Println("Files by Implant ID:", filesByImplantID)

	// Get file by implant ID and name
	fileByIDAndName, err := db.GetFileByImplantIDAndName("implant1", "sample.txt")
	if err != nil {
		fmt.Println("Error getting file by implant ID and name:", err)
		return
	}
	fmt.Println("File by Implant ID and Name:", fileByIDAndName)
}
