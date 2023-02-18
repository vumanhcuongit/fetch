package fetch

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestParseHTMLContent(t *testing.T) {
	fetcher := &Fetcher{
		logger: zap.L().Sugar(),
	}

	testHTML := `<html><head><title>Test Page</title></head><body><p>Hello, world!</p><a href="/test.html"><img src="file:///test.jpg"></a></body></html>`
	tempFile, err := ioutil.TempFile("", "test*.html")
	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.Write([]byte(testHTML)); err != nil {
		t.Fatalf("error writing to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("error closing temp file: %v", err)
	}

	input := ParseHTMLContentInput{
		URL:                 "http://localhost",
		Body:                ioutil.NopCloser(bytes.NewReader([]byte(testHTML))),
		ShouldFetchMetadata: true,
		MaxConcurrent:       10,
	}

	metadata, err := fetcher.ParseHTMLContent(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if metadata.Site != "localhost" {
		t.Errorf("unexpected site: got %s, want localhost", metadata.Site)
	}
	if metadata.NumLinks != 1 {
		t.Errorf("unexpected num_links: got %d, want 0", metadata.NumLinks)
	}
	if metadata.NumImages != 1 {
		t.Errorf("unexpected num_images: got %d, want 0", metadata.NumImages)
	}
	if metadata.LastFetch.IsZero() {
		t.Errorf("unexpected last_fetch: got zero time")
	}

	input.Body = nil
	_, err = fetcher.ParseHTMLContent(input)
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}
