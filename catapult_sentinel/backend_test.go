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
		name    string
		fields  fields
		args    args
		want    File
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got, err := c.GetFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got, err := c.GetFiles(tt.args.filePaths)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got, err := c.GetExperimentsByNames(tt.args.experimentNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExperimentsByNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got, err := c.UpdateFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got, err := c.UpdateExperiments(tt.args.experiments)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateExperiments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got, err := c.CreateCatapultRunConfig(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCatapultRunConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got, err := c.FilterCatapultRunConfig(tt.args.prefix, tt.args.experimentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterCatapultRunConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got, err := c.GetFolderWatchingLocation(tt.args.folderPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFolderWatchingLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test GetAllFolderWatchingLocations",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got, err := c.GetAllFolderWatchingLocations()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllFolderWatchingLocations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultBackend{
				Url:    tt.fields.Url,
				Client: tt.fields.Client,
				Token:  tt.fields.Token,
			}
			got, err := c.GetFolderWatchingLocationById(tt.args.folderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFolderWatchingLocationById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Id == 0 {
				t.Errorf("GetFolderWatchingLocationById() Id = %v, want non-zero", got.Id)
			}
		})
	}
}
