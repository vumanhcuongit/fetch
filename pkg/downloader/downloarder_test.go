package downloader

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestAssetDownloader_DownloadAssets(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/allowed.jpg" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/disallowed.txt" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	logger, _ := zap.NewDevelopment()

	mockClient := mockServer.Client()

	assetDownloader := NewAssetDownloader(logger.Sugar(), mockClient, 1)

	tests := []struct {
		name       string
		urls       []string
		wantPaths  []string
		wantErrors bool
	}{
		{
			name:       "download allowed asset",
			urls:       []string{fmt.Sprintf("%s/allowed.jpg", mockServer.URL)},
			wantPaths:  []string{"assets/" + mockServer.URL[7:] + "/allowed.jpg"},
			wantErrors: false,
		},
		{
			name:       "skip disallowed asset",
			urls:       []string{fmt.Sprintf("%s/disallowed.txt", mockServer.URL)},
			wantPaths:  []string{},
			wantErrors: false,
		},
		{
			name:       "error on invalid URL",
			urls:       []string{"invalid/url"},
			wantPaths:  []string{},
			wantErrors: true,
		},
		{
			name:       "error on nonexistent asset",
			urls:       []string{fmt.Sprintf("%s/nonexistent.jpg", mockServer.URL)},
			wantPaths:  []string{},
			wantErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPaths := assetDownloader.DownloadAssets(tt.urls)

			if tt.wantErrors && len(gotPaths) > 0 {
				t.Errorf("DownloadAssets() returned unexpected paths when errors were expected: got %v, want %v", gotPaths, tt.wantPaths)
			}

			if !tt.wantErrors && !reflect.DeepEqual(gotPaths, tt.wantPaths) {
				t.Errorf("DownloadAssets() returned unexpected paths: got %v, want %v", gotPaths, tt.wantPaths)
			}
		})
	}

	// clean up assets
	err := os.RemoveAll("assets")
	if err != nil {
		log.Fatalf("failed to remove assets folder: %v", err)
	}
}
