package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/noatgnu/catapultSentinel/catapult_sentinel"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var token *string
var backendURL *string
var interval *time.Duration
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
		exp := catapult_sentinel.Experiment{ExperimentName: parentFolder}
		config := catapult_sentinel.CatapultRunConfig{
			ConfigFilePath:         filePath,
			FolderWatchingLocation: folderWatchingLocation,
			Experiment:             exp.Id,
			Content:                configData,
		}
		// Send HTTP request to create CatapultRunConfig
		configJson, _ := json.Marshal(config)
		catapultBackend.Client.Post(fmt.Sprintf("%s/api/catapult_run_config/", catapultBackend.Url), "application/json", bytes.NewBuffer(configJson))
		fmt.Printf("config file %s loaded successfully\n", filePath)
	}
}

func main() {
	backendURL = flag.String("backend-url", "http://localhost:8080", "The backend URL")
	token = flag.String("token", "", "The token")
	interval = flag.Duration("interval", time.Minute, "The scan interval")
	flag.Parse()

	catapultBackend = &catapult_sentinel.CatapultBackend{
		Url:    *backendURL,
		Client: &http.Client{},
		Token:  *token,
	}

	db, err := catapult_sentinel.InitDB("fileinfo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	folderWatchingLocations, err := catapultBackend.GetAllFolderWatchingLocations()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, folder := range folderWatchingLocations {
				tasks, err := catapult_sentinel.ScanFolder(folder, db)
				if err != nil {
					log.Println(err)
					continue
				}
				// get list of paths from tasks.NewFile
				var newPaths []string
				for _, file := range tasks.NewFile {
					newPaths = append(newPaths, file.FilePath)
				}
				if len(newPaths) > 0 {
					results, err := catapultBackend.GetFiles(newPaths)
					if err != nil {
						log.Println(err)
						continue
					}

				}

			}
		}
	}
}
