package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dovelus/corb4n-c2/server/comunication"
	_ "github.com/mattn/go-sqlite3"
)

// Embed the database schema
//
//go:embed schema.sql
var schema string

var dbConn *sql.DB
var dbPath string

// InitDB initializes the database connection and schema
func InitDB() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		comunication.Logger.Fatal("Error getting user home directory: ", err)
	}
	dbDir := filepath.Join(homeDir, "Programming", "corb4n-c2", ".corb4n-c2") // Change this to the correct path when in production
	dbPath = filepath.Join(dbDir, "corb4n-c2.db")
	fmt.Println(dbPath)

	// Ensure the directory exists
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		comunication.Logger.Fatal("Error creating database directory: ", err)
	}

	// Open the database connection
	dbConn, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		comunication.Logger.Fatal("Error opening database file: ", err)
	}

	// Check if the database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Executes the DB schema
		_, err = dbConn.Exec(schema)
		if err != nil {
			comunication.Logger.Fatal("Error executing schema: ", err)
		} else {
			comunication.Logger.Info("Schema Executed")
		}
	} else {
		comunication.Logger.Info("Database Already Exists")
	}
}

// Gracefully shuts down the database connection
func CloseDB() {
	comunication.Logger.Warn("Closing Database Connection")
	if dbConn != nil {
		dbConn.Close()
	}
}
