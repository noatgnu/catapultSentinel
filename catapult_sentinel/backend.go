package catapult_sentinel

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type CatapultBackend struct {
	Url    string
	Client *http.Client
	Token  string
}

type File struct {
	FilePath               string `json:"file_path"`
	FolderWatchingLocation int    `json:"folder_watching_location"`
	Size                   int64  `json:"size"`
	Experiment             int    `json:"experiment"`
	Processing             bool   `json:"processing"`
	ReadyForProcessing     bool   `json:"ready_for_processing"`
	Id                     int    `json:"id"`
}

type FolderWatchingLocation struct {
	FolderPath    string `json:"folder_path"`
	Extensions    string `json:"extensions"`
	IgnoreTerm    string `json:"ignore_term"`
	NetworkFolder bool   `json:"network_folder"`
	Id            int    `json:"id"`
}

type Experiment struct {
	ExperimentName string `json:"experiment_name"`
	Id             int    `json:"id"`
	Vendor         string `json:"vendor"`
	SampleCount    int    `json:"sample_count"`
}

type CatapultRunConfig struct {
	ConfigFilePath          string                 `json:"config_file_path"`
	FolderWatchingLocation  int                    `json:"folder_watching_location"`
	Experiment              int                    `json:"experiment"`
	Content                 map[string]interface{} `json:"content"`
	Id                      int                    `json:"id"`
	FastaReady              bool                   `json:"fasta_ready"`
	FastaRequired           bool                   `json:"fasta_required"`
	SpectralLibraryReady    bool                   `json:"spectral_library_ready"`
	SpectralLibraryRequired bool                   `json:"spectral_library_required"`
}

type CatapultRunConfigQuery struct {
	Results  []CatapultRunConfig `json:"results"`
	Next     string              `json:"next"`
	Previous string              `json:"previous"`
	Count    int                 `json:"count"`
}

func NewCatapultBackend(url string, token string) *CatapultBackend {
	return &CatapultBackend{Url: url, Client: &http.Client{}, Token: token}
}

func (c *CatapultBackend) GetUrl() string {
	return c.Url
}

func (c *CatapultBackend) GetFile(filePath string) (File, error) {
	baseUrl, err := url.Parse(c.Url + "api/files/get_exact_path/")
	if err != nil {
		return File{}, err
	}

	params := url.Values{}
	baseUrl.RawQuery = params.Encode()

	body := struct {
		FilePath string `json:"file_path"`
		Create   bool   `json:"create"`
	}{
		FilePath: filePath,
		Create:   true,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return File{}, err
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return File{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return File{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return File{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var file File
	err = decoder.Decode(&file)
	if err != nil {
		return File{}, err
	}
	return file, nil
}

func (c *CatapultBackend) GetFiles(filePaths []string) ([]File, error) {
	baseUrl, err := url.Parse(c.Url + "api/files/get_exact_paths/")
	if len(filePaths) == 0 {
		return []File{}, nil
	}
	if err != nil {
		return []File{}, err
	}
	params := url.Values{}
	baseUrl.RawQuery = params.Encode()
	body := struct {
		FilePaths []string `json:"file_paths"`
		Create    bool     `json:"create"`
	}{
		FilePaths: filePaths,
		Create:    true,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return []File{}, err
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return []File{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return []File{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var files []File
	err = decoder.Decode(&files)
	if err != nil {
		return []File{}, err
	}
	return files, nil

}

func (c *CatapultBackend) GetFileById(fileId int) (File, error) {
	baseUrl, err := url.Parse(c.Url + "api/files/" + strconv.Itoa(fileId) + "/")
	if err != nil {
		return File{}, err
	}

	params := url.Values{}
	baseUrl.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		return File{}, err
	}
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return File{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return File{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var file File
	err = decoder.Decode(&file)
	if err != nil {
		return File{}, err
	}
	return file, nil
}

func (c *CatapultBackend) CreateFile(file File) (File, error) {
	baseUrl, err := url.Parse(c.Url + "api/files/")
	if err != nil {
		return File{}, err
	}

	bodyJson, err := json.Marshal(file)
	if err != nil {
		return File{}, err
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return File{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return File{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return File{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var newFile File
	err = decoder.Decode(&newFile)
	if err != nil {
		return File{}, err
	}
	return newFile, nil
}

func (c *CatapultBackend) UpdateFile(file File) (File, error) {
	baseUrl, err := url.Parse(c.Url + "api/files/" + strconv.Itoa(file.Id) + "/")
	if err != nil {
		return File{}, err
	}

	params := url.Values{}
	baseUrl.RawQuery = params.Encode()

	bodyJson, err := json.Marshal(file)
	if err != nil {
		return File{}, err
	}

	req, err := http.NewRequest("PUT", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return File{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return File{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return File{}, err
	}

	decoder := json.NewDecoder(resp.Body)

	var updatedFile File
	err = decoder.Decode(&updatedFile)
	if err != nil {
		return File{}, err
	}

	return updatedFile, nil
}

func (c *CatapultBackend) UpdateFiles(files []File) ([]File, error) {
	baseUrl, err := url.Parse(c.Url + "api/files/update_multiple/")
	if err != nil {
		return []File{}, err
	}

	if len(files) == 0 {
		return []File{}, nil
	}

	body := struct {
		Files []File `json:"files"`
	}{
		Files: files,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return []File{}, err
	}

	req, err := http.NewRequest("PUT", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return []File{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return []File{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []File{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var updatedFiles []File
	err = decoder.Decode(&updatedFiles)
	if err != nil {
		return []File{}, err
	}
	return updatedFiles, nil
}

func (c *CatapultBackend) CreateExperiment(experiment Experiment) (Experiment, error) {
	baseUrl, err := url.Parse(c.Url + "api/experiments/")
	if err != nil {
		return Experiment{}, err
	}

	bodyJson, err := json.Marshal(experiment)
	if err != nil {
		return Experiment{}, err
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return Experiment{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return Experiment{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return Experiment{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var newExperiment Experiment
	err = decoder.Decode(&newExperiment)
	if err != nil {
		return Experiment{}, err
	}
	return newExperiment, nil
}

func (c *CatapultBackend) GetExperimentByName(experimentName string) (Experiment, error) {
	baseUrl, err := url.Parse(c.Url + "api/experiments/get_exact_name/")
	if err != nil {
		return Experiment{}, err
	}

	body := struct {
		ExperimentName string `json:"experiment_name"`
		Create         bool   `json:"create"`
	}{
		ExperimentName: experimentName,
		Create:         true,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return Experiment{}, err
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return Experiment{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return Experiment{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Experiment{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var experiment Experiment
	err = decoder.Decode(&experiment)
	if err != nil {
		return Experiment{}, err
	}
	return experiment, nil
}

func (c *CatapultBackend) GetExperimentsByNames(experimentNames []string) ([]Experiment, error) {
	baseUrl, err := url.Parse(c.Url + "api/experiments/get_exact_names/")
	if err != nil {
		return []Experiment{}, err
	}

	body := struct {
		ExperimentNames []string `json:"experiment_names"`
		Create          bool     `json:"create"`
	}{
		ExperimentNames: experimentNames,
		Create:          true,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return []Experiment{}, err
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return []Experiment{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return []Experiment{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Experiment{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var experiments []Experiment
	err = decoder.Decode(&experiments)
	if err != nil {
		return []Experiment{}, err
	}
	return experiments, nil
}

func (c *CatapultBackend) GetExperimentById(experimentId int) (Experiment, error) {
	baseUrl, err := url.Parse(c.Url + "api/experiments/" + string(rune(experimentId)) + "/")
	if err != nil {
		return Experiment{}, err
	}

	params := url.Values{}
	baseUrl.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		return Experiment{}, err
	}

	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return Experiment{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Experiment{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var experiment Experiment
	err = decoder.Decode(&experiment)
	if err != nil {
		return Experiment{}, err
	}
	return experiment, nil
}

func (c *CatapultBackend) UpdateExperiment(experiment Experiment) (Experiment, error) {
	baseUrl, err := url.Parse(c.Url + "api/experiments/" + strconv.Itoa(experiment.Id) + "/")
	if err != nil {
		return Experiment{}, err
	}

	params := url.Values{}
	baseUrl.RawQuery = params.Encode()

	bodyJson, err := json.Marshal(experiment)
	if err != nil {
		return Experiment{}, err
	}

	req, err := http.NewRequest("PUT", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return Experiment{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return Experiment{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Experiment{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var updatedExperiment Experiment
	err = decoder.Decode(&updatedExperiment)
	if err != nil {
		return Experiment{}, err
	}
	return updatedExperiment, nil
}

func (c *CatapultBackend) UpdateExperiments(experiments []Experiment) ([]Experiment, error) {
	baseUrl, err := url.Parse(c.Url + "api/experiments/update_multiple/")
	if err != nil {
		return []Experiment{}, err
	}

	if len(experiments) == 0 {
		return []Experiment{}, nil
	}

	body := struct {
		Experiments []Experiment `json:"experiments"`
	}{
		Experiments: experiments,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return []Experiment{}, err
	}

	req, err := http.NewRequest("PUT", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return []Experiment{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return []Experiment{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Experiment{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var updatedExperiments []Experiment
	err = decoder.Decode(&updatedExperiments)
	if err != nil {
		return []Experiment{}, err
	}

	return updatedExperiments, nil
}

func (c *CatapultBackend) CreateCatapultRunConfig(config CatapultRunConfig) (CatapultRunConfig, error) {
	baseUrl, err := url.Parse(c.Url + "api/catapultrunconfig/")
	if err != nil {
		return CatapultRunConfig{}, err
	}

	bodyJson, err := json.Marshal(config)
	if err != nil {
		return CatapultRunConfig{}, err
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return CatapultRunConfig{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return CatapultRunConfig{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return CatapultRunConfig{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var newConfig CatapultRunConfig
	err = decoder.Decode(&newConfig)
	if err != nil {
		return CatapultRunConfig{}, err
	}
	return newConfig, nil
}

func (c *CatapultBackend) FilterCatapultRunConfig(prefix string, experimentId int) (CatapultRunConfigQuery, error) {
	baseUrl, err := url.Parse(c.Url + "api/catapultrunconfig/")
	if err != nil {
		return CatapultRunConfigQuery{}, err
	}

	params := url.Values{}
	if prefix != "" {
		params.Add("prefix", prefix)
	}
	if experimentId != 0 {
		params.Add("experiment", strconv.Itoa(experimentId))
	}

	baseUrl.RawQuery = params.Encode()
	log.Printf("URL: %s", baseUrl.String())

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		return CatapultRunConfigQuery{}, err
	}
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return CatapultRunConfigQuery{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CatapultRunConfigQuery{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var configs CatapultRunConfigQuery
	err = decoder.Decode(&configs)
	if err != nil {
		return CatapultRunConfigQuery{}, err
	}
	return configs, nil
}

func (c *CatapultBackend) GetFolderWatchingLocation(folderPath string) (FolderWatchingLocation, error) {
	baseUrl, err := url.Parse(c.Url + "api/folderlocations/get_exact_path/")
	if err != nil {
		return FolderWatchingLocation{}, err
	}

	body := struct {
		FolderPath string `json:"folder_path"`
		Create     bool   `json:"create"`
	}{
		FolderPath: folderPath,
		Create:     true,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return FolderWatchingLocation{}, err
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return FolderWatchingLocation{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return FolderWatchingLocation{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FolderWatchingLocation{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var folderWatchingLocation FolderWatchingLocation
	err = decoder.Decode(&folderWatchingLocation)
	if err != nil {
		return FolderWatchingLocation{}, err
	}
	return folderWatchingLocation, nil
}

func (c *CatapultBackend) GetAllFolderWatchingLocations() ([]FolderWatchingLocation, error) {
	baseUrl, err := url.Parse(c.Url + "api/folderlocations/get_all_paths/")
	if err != nil {
		return []FolderWatchingLocation{}, err
	}

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		return []FolderWatchingLocation{}, err
	}

	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return []FolderWatchingLocation{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []FolderWatchingLocation{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var folderWatchingLocations []FolderWatchingLocation
	err = decoder.Decode(&folderWatchingLocations)
	if err != nil {
		return []FolderWatchingLocation{}, err
	}
	return folderWatchingLocations, nil
}

func (c *CatapultBackend) GetFolderWatchingLocationById(folderId int) (FolderWatchingLocation, error) {
	baseUrl, err := url.Parse(c.Url + "api/folderlocations/" + strconv.Itoa(folderId) + "/")
	if err != nil {
		return FolderWatchingLocation{}, err
	}

	params := url.Values{}
	baseUrl.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		return FolderWatchingLocation{}, err
	}

	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return FolderWatchingLocation{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FolderWatchingLocation{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var folderWatchingLocation FolderWatchingLocation
	err = decoder.Decode(&folderWatchingLocation)
	if err != nil {
		return FolderWatchingLocation{}, err
	}
	return folderWatchingLocation, nil
}
