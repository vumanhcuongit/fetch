package fetch

import (
	"fetch-go/pkg/utils"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Fetcher struct {
	logger *zap.SugaredLogger
}

type FetchOptions struct {
	ShouldFetchMetadata bool
	MaxConcurrent       int
}

func NewFetcher(logger *zap.SugaredLogger) *Fetcher {
	return &Fetcher{
		logger: logger,
	}
}

func (f *Fetcher) ScrapeURLs(urls []string, options FetchOptions) error {
	var wg sync.WaitGroup
	wg.Add(len(urls))

	var returnErr error
	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			metadata, err := f.FetchURL(url, options)
			if err != nil {
				f.logger.Infof("error downloading %s: %v\n", url, err)
				if returnErr == nil {
					returnErr = err
				}
				return
			}
			if options.ShouldFetchMetadata {
				fmt.Printf(`
site: %s
num_links: %d
images: %d
last_fetch: %s
`, metadata.Site, metadata.NumLinks, metadata.NumImages, metadata.LastFetch.Format("Mon, 02 Jan 2006 15:04:05 MST"))
			}
		}(url)
	}
	wg.Wait()

	if returnErr != nil {
		return returnErr
	}

	return nil
}

func (f *Fetcher) FetchURL(url string, options FetchOptions) (*Metadata, error) {
	reader, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	input := ParseHTMLContentInput{
		URL:                 url,
		Body:                reader,
		ShouldFetchMetadata: options.ShouldFetchMetadata,
		MaxConcurrent:       options.MaxConcurrent,
	}

	return f.ParseHTMLContent(input)
}
