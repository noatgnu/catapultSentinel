package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/noatgnu/catapultSentinel/catapult_sentinel"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var token *string
var backendURL *string
var catapultBackend *catapult_sentinel.CatapultBackend

func loadConfigYaml(filePath string, folderWatchingLocation int) {
	fmt.Printf("loading config file %s\n", filePath)
	configData := make(map[string]interface{})
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error loading yaml file %s", filePath)
		return
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&configData)
	if err != nil {
		log.Printf("Error decoding yaml file %s", filePath)
		return
	}
	if configData["cat_ready"] == true {
		fmt.Printf("config file %s is ready\n", filePath)
		parentFolder := filepath.Dir(filePath)

		exp := catapultBackend.GetExperimentByName(parentFolder)
		config := catapult_sentinel.CatapultRunConfig{
			ConfigFilePath:         filePath,
			FolderWatchingLocation: folderWatchingLocation,
			Experiment:             exp.Id,
			Content:                configData,
		}
		catapultBackend.CreateCatapultRunConfig(config)
		fmt.Printf("config file %s loaded successfully\n", filePath)
	}
}

func getFolderSize(folderPath string) int64 {
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

func initialScan(folderWatchingLocation catapult_sentinel.FolderWatchingLocation) {
	err := filepath.Walk(folderWatchingLocation.FolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(info.Name(), folderWatchingLocation.IgnoreTerm) {
			return nil
		}
		extension := filepath.Ext(info.Name())
		if strings.Contains(folderWatchingLocation.Extensions, extension) {
			if strings.HasSuffix(info.Name(), ".converted.mzML") {
				fileSize := info.Size()
				fileLocation := path
				if strings.HasSuffix(info.Name(), ".d") {
					fileLocation = filepath.Dir(path)
					fileSize = getFolderSize(fileLocation)
				}
				exp := catapultBackend.GetExperimentByName(filepath.Dir(fileLocation))
				file := catapult_sentinel.File{
					FilePath:               strings.Replace(fileLocation, folderWatchingLocation.FolderPath, "", 1),
					FolderWatchingLocation: folderWatchingLocation.Id,
					Size:                   fileSize,
					Experiment:             exp.Id,
				}
				catapultBackend.GetFile(file.FilePath)
				newFile := catapultBackend.CreateFile(file)

			} else if strings.HasSuffix(info.Name(), ".cat.yml") || strings.HasSuffix(info.Name(), ".cat.yaml") {
				loadConfigYaml(path, folderWatchingLocation.Id)
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}

func watchFolder(folderWatchingLocation catapult_sentinel.FolderWatchingLocation) {
	db, err := catapult_sentinel.InitDB("fileinfo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ticker := time.NewTicker(10 * time.Second) // Adjust the interval as needed
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentFiles := make(map[string]os.FileInfo)
			err := filepath.Walk(folderWatchingLocation.FolderPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && !strings.Contains(info.Name(), folderWatchingLocation.IgnoreTerm) {
					currentFiles[path] = info
				}
				if info.IsDir() && filepath.Ext(info.Name()) == ".d" {
					currentFiles[path] = info
				}
				return nil
			})
			if err != nil {
				log.Println(err)
				continue
			}

			for path, info := range currentFiles {
				var lastModified time.Time
				var size int64
				var isFolder bool

				if strings.HasSuffix(info.Name(), ".d") {
					isFolder = false
					size = getFolderSize(path)
					lastModified = info.ModTime()
				} else {
					isFolder = info.IsDir()
					size = info.Size()
					lastModified = info.ModTime()
				}

				var dbLastModified time.Time
				err := db.QueryRow("SELECT last_modified FROM files WHERE path = ?", path).Scan(&dbLastModified)
				if err == sql.ErrNoRows {
					// New file detected
					handleNewFile(path, info, folderWatchingLocation)
					_, err = db.Exec("INSERT INTO files (path, last_modified, size, is_folder) VALUES (?, ?, ?, ?)", path, lastModified, size, isFolder)
					if err != nil {
						log.Println(err)
					}
				} else if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func handleNewFile(path string, info os.FileInfo, folderWatchingLocation catapult_sentinel.FolderWatchingLocation) catapult_sentinel.File {
	extension := filepath.Ext(info.Name())
	if strings.Contains(folderWatchingLocation.Extensions, extension) {
		fileSize := info.Size()
		fileLocation := path
		if strings.HasSuffix(info.Name(), ".d") {
			fileLocation = filepath.Dir(path)
			fileSize = getFolderSize(fileLocation)
		}
		exp := catapultBackend.GetExperimentByName(filepath.Dir(fileLocation))
		file := catapult_sentinel.File{
			FilePath:               strings.Replace(fileLocation, folderWatchingLocation.FolderPath, "", 1),
			FolderWatchingLocation: folderWatchingLocation.Id,
			Size:                   fileSize,
			Experiment:             exp.Id,
		}
		return catapultBackend.CreateFile(file)
	} else if strings.HasSuffix(info.Name(), ".cat.yml") || strings.HasSuffix(info.Name(), ".cat.yaml") {
		loadConfigYaml(path, folderWatchingLocation.Id)
	}
	return catapult_sentinel.File{}
}

func main() {
	backendURL = flag.String("backend-url", "http://localhost:8080", "The backend URL")
	token = flag.String("token", "", "The token")
	flag.Parse()

	catapultBackend = &catapult_sentinel.CatapultBackend{
		Url:    *backendURL,
		Client: &http.Client{},
		Token:  *token,
	}

	folderWatchingLocations := catapultBackend.GetAllFolderWatchingLocations()

	for _, folder := range folderWatchingLocations {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in f", r)
				}
			}()
			initialScan(folder)
			watchFolder(folder)
		}()
	}

	select {}
}
