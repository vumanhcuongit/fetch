package fetch

import (
	"errors"
	"fetch-go/pkg/downloader"
	"fetch-go/pkg/utils"
	"io"
	"net/url"
	netURL "net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Metadata struct {
	Site      string
	NumLinks  int
	NumImages int
	LastFetch time.Time
}

type ParseHTMLContentInput struct {
	URL                 string
	Body                io.ReadCloser
	ShouldFetchMetadata bool
	MaxConcurrent       int
}

func (f *Fetcher) ParseHTMLContent(
	input ParseHTMLContentInput,
) (*Metadata, error) {
	assetDownloader := downloader.NewAssetDownloader(
		f.logger,
		utils.GetHTTPClient(),
		input.MaxConcurrent,
	)

	// parse the HTML document
	if input.Body == nil {
		return nil, errors.New("invalid body")
	}
	doc, err := goquery.NewDocumentFromReader(input.Body)
	if err != nil {
		f.logger.Errorf("error parsing document: %+v", err)
		return nil, err
	}

	assetURLs := f.getAssetURLs(input.URL, doc)

	// download the assets
	localPaths := assetDownloader.DownloadAssets(assetURLs)

	// replace the URLs in the HTML with the local paths
	f.updateHTMLWithLocalPaths(doc, localPaths)

	// save the modified HTML to a file
	filename := path.Base(input.URL) + ".html"
	if err := f.saveHTMLToFile(doc, filename); err != nil {
		f.logger.Errorf("Error writing file: %+v", err)
		return nil, err
	}

	if input.ShouldFetchMetadata {
		parsedURL, err := url.Parse(input.URL)
		if err != nil {
			return nil, err
		}

		metadata := &Metadata{
			Site:      parsedURL.Hostname(),
			NumLinks:  doc.Find("a").Length(),
			NumImages: doc.Find("img").Length(),
			LastFetch: time.Now(),
		}
		return metadata, nil
	}

	return nil, nil
}

func (f *Fetcher) getAssetURLs(baseURL string, doc *goquery.Document) []string {
	assetURLs := []string{}
	doc.Find("[src], [href]").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("src")
		if !ok {
			url, ok = s.Attr("href")
			if !ok {
				return
			}
		}

		// skip URLs that are already local
		if strings.HasPrefix(url, "file:///") {
			return
		}

		// check if it's relative URL
		u, err := netURL.Parse(url)
		if err != nil {
			f.logger.Errorf("error while parsing the image URL: %+v", err)
			return
		}

		// relative URL
		if u.Host == "" {
			url = baseURL + "/" + url
		}

		assetURLs = append(assetURLs, url)
	})

	return assetURLs
}

func (f *Fetcher) updateHTMLWithLocalPaths(doc *goquery.Document, localPaths []string) {
	doc.Find("[src], [href]").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("src")
		if !ok {
			url, ok = s.Attr("href")
			if !ok {
				return
			}
		}
		// find the local path for this URL
		var localPath string
		for _, path := range localPaths {
			if strings.Contains(path, url) || strings.Contains(url, strings.Replace(path, "assets/", "", 1)) {
				localPath = path
				break
			}
		}

		if localPath != "" {
			s.SetAttr("src", localPath)
			s.SetAttr("href", localPath)
		}
	})
}

func (f *Fetcher) saveHTMLToFile(doc *goquery.Document, filename string) error {
	html, err := doc.Html()
	if err != nil {
		f.logger.Errorf("Error generating HTML: %+v", err)
		return err
	}

	err = os.WriteFile(filename, []byte(html), 0644)
	if err != nil {
		f.logger.Errorf("Error writing file: %+v", err)
		return err
	}

	return nil
}
