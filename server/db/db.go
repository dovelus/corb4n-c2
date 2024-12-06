package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"

	_ "github.com/mattn/go-sqlite3"
)

// Embed the database schema
//
//go:embed sql.schema
var schema string

var dbConn *sql.DB
var dbPath string

var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.RFC3339,
})

// InitDB initializes the database connection and schema
func InitDB() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Fatal("Error getting user home directory: ", err)
	}
	dbDir := filepath.Join(homeDir, "Programming", "corb4n-c2", ".corb4n-c2")
	dbPath = filepath.Join(dbDir, "corb4n-c2.db")
	fmt.Println(dbPath)

	// Ensure the directory exists
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		logger.Fatal("Error creating database directory: ", err)
	}

	// Open the database connection
	dbConn, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatal("Error opening database file: ", err)
	}

	// Check if the database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Executes the DB schema
		_, err = dbConn.Exec(schema)
		if err != nil {
			logger.Fatal("Error executing schema: ", err)
		} else {
			logger.Info("Schema Executed")
		}
	} else {
		logger.Info("Database Already Exists")
	}
}

// Gracefully shuts down the database connection
func CloseDB() {
	logger.Warn("Closing Database Connection")
	if dbConn != nil {
		dbConn.Close()
	}
}
