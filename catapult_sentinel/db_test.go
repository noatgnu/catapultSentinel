package catapult_sentinel

import (
	"testing"
)

func TestInitDB(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB() error: %v", err)
	}
	defer db.Close()

	// Check if the table was created
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='files'").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to query table: %v", err)
	}
	if name != "files" {
		t.Fatalf("Expected table name 'files', got %s", name)
	}
}

func TestCheckFileExists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a test file
	_, err := db.Exec("INSERT INTO files (path, size, is_folder, last_modified, remote_id) VALUES (?, ?, ?, ?, ?)", "test.txt", 123, false, 1234567890, 1)
	if err != nil {
		t.Fatalf("Failed to insert test file: %v", err)
	}

	exists, err := CheckFileExists(db, "test.txt")
	if err != nil {
		t.Fatalf("CheckFileExists() error: %v", err)
	}
	if !exists {
		t.Fatalf("Expected file to exist")
	}

	exists, err = CheckFileExists(db, "nonexistent.txt")
	if err != nil {
		t.Fatalf("CheckFileExists() error: %v", err)
	}
	if exists {
		t.Fatalf("Expected file to not exist")
	}
}

func TestGetFile(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a test file
	_, err := db.Exec("INSERT INTO files (path, size, is_folder, last_modified, remote_id) VALUES (?, ?, ?, ?, ?)", "test.txt", 123, false, 1234567890, 1)
	if err != nil {
		t.Fatalf("Failed to insert test file: %v", err)
	}

	file, err := GetFile(db, "test.txt")
	if err != nil {
		t.Fatalf("GetFile() error: %v", err)
	}
	if file.Path != "test.txt" || file.Size != 123 || file.IsFolder != false || file.LastModified != 1234567890 || file.RemoteId != 1 {
		t.Fatalf("GetFile() returned unexpected result: %+v", file)
	}
}

func TestInsertFile(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	file := LocalFile{
		Path:         "test.txt",
		Size:         123,
		IsFolder:     false,
		LastModified: 1234567890,
		RemoteId:     1,
	}

	err := InsertFile(db, file)
	if err != nil {
		t.Fatalf("InsertFile() error: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM files WHERE path = ?", file.Path).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query inserted file: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 file, got %d", count)
	}
}

func TestUpdateFile(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a test file
	_, err := db.Exec("INSERT INTO files (path, size, is_folder, last_modified, remote_id) VALUES (?, ?, ?, ?, ?)", "test.txt", 123, false, 1234567890, 1)
	if err != nil {
		t.Fatalf("Failed to insert test file: %v", err)
	}

	file := LocalFile{
		Path:         "test.txt",
		Size:         456,
		IsFolder:     true,
		LastModified: 9876543210,
		RemoteId:     2,
	}

	err = UpdateFile(db, file)
	if err != nil {
		t.Fatalf("UpdateFile() error: %v", err)
	}

	var updatedFile LocalFile
	err = db.QueryRow("SELECT path, size, is_folder, last_modified, remote_id FROM files WHERE path = ?", file.Path).Scan(&updatedFile.Path, &updatedFile.Size, &updatedFile.IsFolder, &updatedFile.LastModified, &updatedFile.RemoteId)
	if err != nil {
		t.Fatalf("Failed to query updated file: %v", err)
	}
	if updatedFile != file {
		t.Fatalf("UpdateFile() returned unexpected result: %+v", updatedFile)
	}
}

func TestUpdateMultipleFiles(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test files
	files := []LocalFile{
		{"test1.txt", 123, false, 1234567890, 1},
		{"test2.txt", 456, true, 9876543210, 2},
	}
	for _, file := range files {
		_, err := db.Exec("INSERT INTO files (path, size, is_folder, last_modified, remote_id) VALUES (?, ?, ?, ?, ?)", file.Path, file.Size, file.IsFolder, file.LastModified, file.RemoteId)
		if err != nil {
			t.Fatalf("Failed to insert test file: %v", err)
		}
	}

	// Update files
	updatedFiles := []LocalFile{
		{"test1.txt", 789, true, 1111111111, 3},
		{"test2.txt", 101112, false, 2222222222, 4},
	}
	err := UpdateMultipleFiles(db, updatedFiles)
	if err != nil {
		t.Fatalf("UpdateMultipleFiles() error: %v", err)
	}

	for _, file := range updatedFiles {
		var updatedFile LocalFile
		err = db.QueryRow("SELECT path, size, is_folder, last_modified, remote_id FROM files WHERE path = ?", file.Path).Scan(&updatedFile.Path, &updatedFile.Size, &updatedFile.IsFolder, &updatedFile.LastModified, &updatedFile.RemoteId)
		if err != nil {
			t.Fatalf("Failed to query updated file: %v", err)
		}
		if updatedFile != file {
			t.Fatalf("UpdateMultipleFiles() returned unexpected result: %+v", updatedFile)
		}
	}
}
