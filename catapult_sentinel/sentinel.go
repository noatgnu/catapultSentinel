package catapult_sentinel

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Task struct {
	NewFile     []File
	ChangedFile []File
}

func GetFolderSize(folderPath string) int64 {
	var totalSize int64
	err := filepath.Walk(folderPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		totalSize += info.Size()
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return totalSize
}

func ScanFolder(location FolderWatchingLocation, db *sql.DB) (Task, error) {

	currentFiles := make(map[string]os.FileInfo)
	err := filepath.Walk(location.FolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (!strings.Contains(info.Name(), location.IgnoreTerm) || strings.HasSuffix(info.Name(), ".cat.yml")) {
			currentFiles[path] = info
		}
		if info.IsDir() && filepath.Ext(info.Name()) == ".d" {
			currentFiles[path] = info
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return Task{}, err
	}
	task := Task{
		NewFile:     []File{},
		ChangedFile: []File{},
	}

	for path, info := range currentFiles {
		var localFile LocalFile
		exists, _ := CheckFileExists(db, path)
		if !exists {
			localFile = LocalFile{
				IsFolder:     info.IsDir(),
				Size:         info.Size(),
				LastModified: info.ModTime().Unix(),
				RemoteId:     0,
				Path:         path,
			}
			err := InsertFile(db, localFile)
			if err != nil {
				log.Println(err)
			}

			task.NewFile = append(task.NewFile, File{
				FilePath:               localFile.Path,
				FolderWatchingLocation: location.Id,
				Size:                   localFile.Size,
			})
		} else {
			localFile, err = GetFile(db, path)
			if err != nil {
				log.Println(err)
			}
			if localFile.Size != info.Size() || localFile.LastModified != info.ModTime().Unix() {
				task.ChangedFile = append(task.ChangedFile, File{
					FilePath:               localFile.Path,
					FolderWatchingLocation: location.Id,
					Size:                   localFile.Size,
				})

			}

		}

	}
	return task, nil
}
