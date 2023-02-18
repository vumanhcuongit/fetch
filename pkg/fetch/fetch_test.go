package fetch

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestScrapeURLs(t *testing.T) {
	fetcher := &Fetcher{
		logger: zap.L().Sugar(),
	}

	testCases := []struct {
		urls        []string
		options     FetchOptions
		expectError bool
	}{
		{
			urls:        []string{"http://example.com"},
			options:     FetchOptions{ShouldFetchMetadata: true, MaxConcurrent: 5},
			expectError: false,
		},
		{
			urls:        []string{"invalid-url"},
			options:     FetchOptions{MaxConcurrent: 5},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		err := fetcher.ScrapeURLs(tc.urls, tc.options)

		if err != nil && !tc.expectError {
			t.Errorf("do not expect error but got %v", err)
		}
	}
}

func TestScrapeURLsWithMetadata(t *testing.T) {
	fetcher := &Fetcher{
		logger: zap.L().Sugar(),
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	urls := []string{"http://example.com"}
	options := FetchOptions{ShouldFetchMetadata: true, MaxConcurrent: 100}

	err := fetcher.ScrapeURLs(urls, options)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	w.Close()
	out, _ := io.ReadAll(r)
	// restore the stdout
	os.Stdout = old

	expected := fmt.Sprintf(`
site: %s
num_links: %d
images: %d
last_fetch: %s
`, "example.com", 1, 0, time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST"))
	if string(out) != expected {
		t.Errorf("unexpected output: got %q, want %q", string(out), expected)
	}
}
