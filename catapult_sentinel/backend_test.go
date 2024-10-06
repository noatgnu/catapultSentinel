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
			name: "Test 1",
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
			if got := c.GetFile(tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFile() = %v, want %v", got, tt.want)
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
			name: "Test 1",
			fields: fields{
				Url:    "http://localhost:8000/",
				Client: &http.Client{},
				Token:  os.Getenv("API_TOKEN"),
			},
			args: args{
				prefix:       "1",
				experimentId: 0,
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
			query := c.FilterCatapultRunConfig(tt.args.prefix, tt.args.experimentId)
			if query.Count == 0 {
				t.Errorf("FilterCatapultRunConfig() = %v, want > %v", query.Count, 0)
			} else {
				t.Logf("FilterCatapultRunConfig() = %v, want > %v", query.Count, 0)
			}

		})
	}
}
