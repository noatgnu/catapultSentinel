package catapult_sentinel

import (
	"database/sql"
	"testing"
)

// setupTestDB sets up an in-memory SQLite database for testing purposes.
// It creates the necessary tables for storing file sizes and copied status.
//
// Parameters:
// - t: The testing object.
//
// Returns:
// - *sql.DB: The initialized in-memory SQLite database.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", "file:/foobar?vfs=memdb")
	if err != nil {
		t.Fatalf("sql.Open() error: %v", err)
	}

	createTableSQL := `
	 CREATE TABLE IF NOT EXISTS files (
	  path TEXT PRIMARY KEY,
	  size INTEGER,
	  is_folder BOOLEAN,
	  last_modified TIMESTAMP
	 );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	return db
}

// InitDB initializes the SQLite database with the given file path.
// It creates the necessary tables for storing file sizes and copied status if they do not exist.
//
// Parameters:
// - dbPath: The file path for the SQLite database.
//
// Returns:
// - *sql.DB: The initialized SQLite database.
// - error: An error object if there was an issue initializing the database.
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	 CREATE TABLE IF NOT EXISTS files (
	  path TEXT PRIMARY KEY,
	  size INTEGER,
	  is_folder BOOLEAN,
	  last_modified TIMESTAMP,
	  remote_id INTEGER
	 );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
