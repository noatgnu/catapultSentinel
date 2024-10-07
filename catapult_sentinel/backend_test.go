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
				filePaths: []string{"D:\\watch_folder\\MRC-Astral\\1000ngHeLa_180SPD_ES906_20240214_01.raw"},
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

func TestCatapultBackend_GetFileById(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		fileId int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test GetFileById",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				fileId: 5,
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
			got := c.GetFileById(tt.args.fileId)
			if got.Id == 0 {
				t.Errorf("GetFileById() Id = %v, want non-zero", got.Id)
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
				experimentNames: []string{"Experiment1", "Experiment2"},
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

func TestCatapultBackend_GetExperimentById(t *testing.T) {
	type fields struct {
		Url    string
		Client *http.Client
		Token  string
	}
	type args struct {
		experimentId int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test GetExperimentById",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				experimentId: 1,
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
			got := c.GetExperimentById(tt.args.experimentId)
			if got.Id == 0 {
				t.Errorf("GetExperimentById() Id = %v, want non-zero", got.Id)
			}
		})
	}
}
