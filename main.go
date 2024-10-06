package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

func loadConfigYaml(filePath string, folderWatchingLocation string, backendURL string) {
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
		exp := Experiment{ExperimentName: parentFolder}
		config := CatapultRunConfig{
			ConfigFilePath:         filePath,
			FolderWatchingLocation: folderWatchingLocation,
			Experiment:             exp.ExperimentName,
			Content:                configData,
		}
		// Send HTTP request to create CatapultRunConfig
		configJson, _ := json.Marshal(config)
		http.Post(fmt.Sprintf("%s/api/catapult_run_config/", backendURL), "application/json", bytes.NewBuffer(configJson))
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

func initialScan(folderWatchingLocation FolderWatchingLocation, backendURL string) {
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
				exp := Experiment{ExperimentName: filepath.Dir(fileLocation)}
				file := File{
					FilePath:               strings.Replace(fileLocation, folderWatchingLocation.FolderPath, "", 1),
					FolderWatchingLocation: folderWatchingLocation.FolderPath,
					Size:                   fileSize,
					Experiment:             exp.ExperimentName,
				}
				// Send HTTP request to create File
				fileJson, _ := json.Marshal(file)
				http.Post(fmt.Sprintf("%s/api/file/", backendURL), "application/json", bytes.NewBuffer(fileJson))
			} else if strings.HasSuffix(info.Name(), ".cat.yml") || strings.HasSuffix(info.Name(), ".cat.yaml") {
				loadConfigYaml(path, folderWatchingLocation.FolderPath, backendURL)
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}

func watchFolder(folderWatchingLocation FolderWatchingLocation, backendURL string) {
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
					handleNewFile(path, info, folderWatchingLocation, backendURL)
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

func handleNewFile(path string, info os.FileInfo, folderWatchingLocation FolderWatchingLocation, backendURL string) {
	extension := filepath.Ext(info.Name())
	if strings.Contains(folderWatchingLocation.Extensions, extension) {
		fileSize := info.Size()
		fileLocation := path
		if strings.HasSuffix(info.Name(), ".d") {
			fileLocation = filepath.Dir(path)
			fileSize = getFolderSize(fileLocation)
		}
		exp := Experiment{ExperimentName: filepath.Dir(fileLocation)}
		file := File{
			FilePath:               strings.Replace(fileLocation, folderWatchingLocation.FolderPath, "", 1),
			FolderWatchingLocation: folderWatchingLocation.FolderPath,
			Size:                   fileSize,
			Experiment:             exp.ExperimentName,
		}
		// Send HTTP request to create File
		fileJson, _ := json.Marshal(file)
		http.Post(fmt.Sprintf("%s/api/file/", backendURL), "application/json", bytes.NewBuffer(fileJson))
	} else if strings.HasSuffix(info.Name(), ".cat.yml") || strings.HasSuffix(info.Name(), ".cat.yaml") {
		loadConfigYaml(path, folderWatchingLocation.FolderPath, backendURL)
	}
}

func main() {
	backendURL = flag.String("backend-url", "http://localhost:8080", "The backend URL")
	token = flag.String("token", "", "The token")
	flag.Parse()

	// Replace with actual API call to get FolderWatchingLocation objects
	folderWatchingLocations := []FolderWatchingLocation{
		{
			FolderPath:    "/path/to/watch",
			Extensions:    ".mzML,.yml,.yaml",
			IgnoreTerm:    "ignore",
			NetworkFolder: false,
		},
	}

	for _, folder := range folderWatchingLocations {
		go initialScan(folder, *backendURL)
		go watchFolder(folder, *backendURL)
	}

	// Keep the main function running
	select {}
}
