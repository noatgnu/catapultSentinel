package catapult_sentinel

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"testing"
)

type LocalFile struct {
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	IsFolder     bool   `json:"is_folder"`
	LastModified int64  `json:"last_modified"`
	RemoteId     int64  `json:"remote_id"`
}

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
	  last_modified TIMESTAMP,
	  remote_id INTEGER
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

func CheckFileExists(db *sql.DB, path string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM files WHERE path = ?)", path).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetFile(db *sql.DB, path string) (LocalFile, error) {
	var file LocalFile
	err := db.QueryRow("SELECT path, size, is_folder, last_modified, remote_id FROM files WHERE path = ?", path).Scan(&file.Path, &file.Size, &file.IsFolder, &file.LastModified, &file.RemoteId)
	if err != nil {
		return LocalFile{}, err
	}
	return file, nil
}

func InsertFile(db *sql.DB, file LocalFile) error {
	_, err := db.Exec("INSERT INTO files (path, size, is_folder, last_modified, remote_id) VALUES (?, ?, ?, ?, ?)", file.Path, file.Size, file.IsFolder, file.LastModified, file.RemoteId)
	return err
}

func UpdateFile(db *sql.DB, file LocalFile) error {
	_, err := db.Exec("UPDATE files SET size = ?, is_folder = ?, last_modified = ?, remote_id = ? WHERE path = ?", file.Size, file.IsFolder, file.LastModified, file.RemoteId, file.Path)
	return err
}

func UpdateMultipleFiles(db *sql.DB, files []LocalFile) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE files SET size = ?, is_folder = ?, last_modified = ?, remote_id = ? WHERE path = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, file := range files {
		_, err = stmt.Exec(file.Size, file.IsFolder, file.LastModified, file.RemoteId, file.Path)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
