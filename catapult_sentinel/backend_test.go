package catapult_sentinel

import (
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestCatapultBackend_GetFile(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   File
	}{
		{
			name: "Test GetFile",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				filePath: "D:\\watch_folder\\MRC-Astral\\1000ngHeLa_180SPD_ES906_20240214_01.raw",
			},
			want: File{
				FilePath:               "D:\\watch_folder\\MRC-Astral\\1000ngHeLa_180SPD_ES906_20240214_01.raw",
				FolderWatchingLocation: 1,
				Size:                   3181910116,
				Experiment:             1,
				Processing:             false,
				ReadyForProcessing:     true,
				Id:                     5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got := c.GetFile(tt.args.filePath)
			if got.Id == 0 {
				t.Errorf("GetFile() Id = %v, want non-zero", got.Id)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCatapultBackend_GetFiles(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		filePaths []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test GetFiles",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				filePaths: []string{
					"D:\\watch_folder\\MRC-Astral\\1000ngHeLa_180SPD_ES906_20240214_01.raw",
					"D:\\watch_folder\\MRC-Astral\\1000ngHeLa_180SPD_ES906_20240214_02.raw",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got := c.GetFiles(tt.args.filePaths)
			if len(got) <= 1 {
				t.Errorf("GetFiles() length = %v, want > 1", len(got))
			}
			for _, file := range got {
				if file.Id == 0 {
					t.Errorf("GetFiles() Id = %v, want non-zero", file.Id)
				}
			}
		})
	}
}

func TestCatapultBackend_GetExperimentsByNames(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		experimentNames []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test GetExperimentsByNames",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				experimentNames: []string{"D:\\watch_folder\\MRC-Astral", "D:\\watch_folder\\MRC-Astral2", "Experiment3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got := c.GetExperimentsByNames(tt.args.experimentNames)
			if len(got) <= 1 {
				t.Errorf("GetExperimentsByNames() length = %v, want > 1", len(got))
			}
			for _, experiment := range got {
				if experiment.Id == 0 {
					t.Errorf("GetExperimentsByNames() Id = %v, want non-zero", experiment.Id)
				}
			}
		})
	}

}

func TestCatapultBackend_UpdateFile(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		file File
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test UpdateFile",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				file: File{
					Id:                 5,
					FilePath:           "D:\\watch_folder\\MRC-Astral\\1000ngHeLa_180SPD_ES906_20240214_01.raw",
					Size:               20000,
					Processing:         false,
					ReadyForProcessing: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got := c.UpdateFile(tt.args.file)
			if got.Id == 0 {
				t.Errorf("UpdateFile() Id = %v, want non-zero", got.Id)
			}
		})
	}
}

func TestCatapultBackend_UpdateExperiments(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		experiments []Experiment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test UpdateExperiments",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				experiments: []Experiment{
					{Id: 1, ExperimentName: "Experiment1"},
					{Id: 2, ExperimentName: "Experiment2"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got := c.UpdateExperiments(tt.args.experiments)
			if len(got) <= 1 {
				t.Errorf("UpdateExperiments() length = %v, want > 1", len(got))
			}
			for _, experiment := range got {
				if experiment.Id == 0 {
					t.Errorf("UpdateExperiments() Id = %v, want non-zero", experiment.Id)
				}
			}
		})
	}
}
func TestCatapultBackend_CreateCatapultRunConfig(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		config CatapultRunConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test CreateCatapultRunConfig",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				config: CatapultRunConfig{
					Id:             0,
					Experiment:     1,
					Content:        make(map[string]interface{}),
					ConfigFilePath: "catapult/management/commands/diann_config.cat.yml",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got := c.CreateCatapultRunConfig(tt.args.config)
			if got.Id == 0 {
				t.Errorf("CreateCatapultRunConfig() Id = %v, want non-zero", got.Id)
			}
		})
	}
}

func TestCatapultBackend_FilterCatapultRunConfig(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		prefix       string
		experimentId int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test FilterCatapultRunConfig",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				prefix:       "1",
				experimentId: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got := c.FilterCatapultRunConfig(tt.args.prefix, tt.args.experimentId)
			if len(got.Results) == 0 {
				t.Errorf("FilterCatapultRunConfig() length = %v, want > 0", len(got.Results))
			}
			for _, config := range got.Results {
				if config.Id == 0 {
					t.Errorf("FilterCatapultRunConfig() Id = %v, want non-zero", config.Id)
				}
			}
		})
	}
}

func TestCatapultBackend_GetFolderWatchingLocation(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		folderPath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test GetFolderWatchingLocation",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				folderPath: "D:\\watch_folder",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got := c.GetFolderWatchingLocation(tt.args.folderPath)
			if got.Id == 0 {
				t.Errorf("GetFolderWatchingLocation() Id = %v, want non-zero", got.Id)
			}
		})
	}
}

func TestCatapultBackend_GetAllFolderWatchingLocations(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test GetAllFolderWatchingLocations",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got := c.GetAllFolderWatchingLocations()
			if len(got) == 0 {
				t.Errorf("GetAllFolderWatchingLocations() length = %v, want > 0", len(got))
			}
			for _, location := range got {
				if location.Id == 0 {
					t.Errorf("GetAllFolderWatchingLocations() Id = %v, want non-zero", location.Id)
				}
			}
		})
	}
}

func TestCatapultBackend_GetFolderWatchingLocationById(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		folderId int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test GetFolderWatchingLocationById",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				folderId: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got := c.GetFolderWatchingLocationById(tt.args.folderId)
			if got.Id == 0 {
				t.Errorf("GetFolderWatchingLocationById() Id = %v, want non-zero", got.Id)
			}
		})
	}
}
