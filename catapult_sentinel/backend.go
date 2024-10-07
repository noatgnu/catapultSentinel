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

func (c *CatapultBackend) GetFile(filePath string) File {
	baseUrl, err := url.Parse(c.Url + "api/files/get_exact_path/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
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
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var file File
	err = decoder.Decode(&file)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return file
}

func (c *CatapultBackend) GetFiles(filePaths []string) []File {
	baseUrl, err := url.Parse(c.Url + "api/files/get_exact_paths/")
	if len(filePaths) == 0 {
		return []File{}
	}
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
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
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var files []File
	err = decoder.Decode(&files)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return files

}

func (c *CatapultBackend) GetFileById(fileId int) File {
	baseUrl, err := url.Parse(c.Url + "api/files/" + strconv.Itoa(fileId) + "/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
	}

	params := url.Values{}
	baseUrl.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}
	log.Printf("Response: %v", resp.Body)

	decoder := json.NewDecoder(resp.Body)
	var file File
	err = decoder.Decode(&file)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return file
}

func (c *CatapultBackend) CreateFile(file File) File {
	baseUrl, err := url.Parse(c.Url + "api/files/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
	}

	bodyJson, err := json.Marshal(file)
	if err != nil {
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var newFile File
	err = decoder.Decode(&newFile)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return newFile
}

func (c *CatapultBackend) UpdateFile(file File) File {
	baseUrl, err := url.Parse(c.Url + "api/files/" + strconv.Itoa(file.Id) + "/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
	}

	params := url.Values{}
	baseUrl.RawQuery = params.Encode()

	bodyJson, err := json.Marshal(file)
	if err != nil {
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("PUT", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}
	return file
}

func (c *CatapultBackend) UpdateFiles(files []File) []File {
	baseUrl, err := url.Parse(c.Url + "api/files/update_multiple/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
	}

	if len(files) == 0 {
		return []File{}
	}

	body := struct {
		Files []File `json:"files"`
	}{
		Files: files,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("PUT", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var updatedFiles []File
	err = decoder.Decode(&updatedFiles)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return updatedFiles
}

func (c *CatapultBackend) CreateExperiment(experiment Experiment) Experiment {
	baseUrl, err := url.Parse(c.Url + "api/experiments/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
	}

	bodyJson, err := json.Marshal(experiment)
	if err != nil {
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var newExperiment Experiment
	err = decoder.Decode(&newExperiment)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return newExperiment
}

func (c *CatapultBackend) GetExperimentByName(experimentName string) Experiment {
	baseUrl, err := url.Parse(c.Url + "api/experiments/get_exact_name/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
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
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var experiment Experiment
	err = decoder.Decode(&experiment)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return experiment
}

func (c *CatapultBackend) GetExperimentsByNames(experimentNames []string) []Experiment {
	baseUrl, err := url.Parse(c.Url + "api/experiments/get_exact_names/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
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
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var experiments []Experiment
	err = decoder.Decode(&experiments)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return experiments
}

func (c *CatapultBackend) GetExperimentById(experimentId int) Experiment {
	baseUrl, err := url.Parse(c.Url + "api/experiments/" + string(rune(experimentId)) + "/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
	}

	params := url.Values{}
	baseUrl.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}

	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var experiment Experiment
	err = decoder.Decode(&experiment)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return experiment
}

func (c *CatapultBackend) UpdateExperiment(experiment Experiment) Experiment {
	baseUrl, err := url.Parse(c.Url + "api/experiments/" + strconv.Itoa(experiment.Id) + "/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
	}

	params := url.Values{}
	baseUrl.RawQuery = params.Encode()

	bodyJson, err := json.Marshal(experiment)
	if err != nil {
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("PUT", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var updatedExperiment Experiment
	err = decoder.Decode(&updatedExperiment)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return updatedExperiment
}

func (c *CatapultBackend) UpdateExperiments(experiments []Experiment) []Experiment {
	baseUrl, err := url.Parse(c.Url + "api/experiments/update_multiple/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
	}

	if len(experiments) == 0 {
		return []Experiment{}
	}

	body := struct {
		Experiments []Experiment `json:"experiments"`
	}{
		Experiments: experiments,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("PUT", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var updatedExperiments []Experiment
	err = decoder.Decode(&updatedExperiments)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}

	return updatedExperiments
}

func (c *CatapultBackend) CreateCatapultRunConfig(config CatapultRunConfig) CatapultRunConfig {
	baseUrl, err := url.Parse(c.Url + "api/catapultrunconfig/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
	}

	bodyJson, err := json.Marshal(config)
	if err != nil {
		log.Panicf("Error marshaling body: %s", err)
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Panicf("Error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}
	return config
}

func (c *CatapultBackend) FilterCatapultRunConfig(prefix string, experimentId int) CatapultRunConfigQuery {
	baseUrl, err := url.Parse(c.Url + "api/catapultrunconfig/")
	if err != nil {
		log.Panicf("Error parsing URL: %s", err)
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
		log.Panicf("Error creating request: %s", err)
	}
	req.Header.Set("Authorization", "Token "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Panicf("Error sending request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: %s", resp.Status)
	}
	log.Printf("Response: %v", resp.Body)

	decoder := json.NewDecoder(resp.Body)
	var configs CatapultRunConfigQuery
	err = decoder.Decode(&configs)
	if err != nil {
		log.Panicf("Error decoding response: %s", err)
	}
	return configs
}
