package fetch

import (
	"fetch-go/pkg/downloader"
	"fetch-go/pkg/utils"
	"fmt"
	"io"
	netURL "net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	maxConcurrent = 100
)

func ParseHTML(baseURL string, body io.ReadCloser) error {
	assetDownloader := downloader.NewAssetDownloader(
		utils.GetHTTPClient(),
		maxConcurrent,
	)

	// parse the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println("error parsing document:", err)
		return err
	}

	// download and replace the URLs of the assets
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
			fmt.Println("Error while parsing the image URL:", err)
			return
		}
		// relative URL
		if u.Host == "" {
			url = baseURL + "/" + url
		}

		assetURLs = append(assetURLs, url)
	})

	// download the assets
	localPaths := assetDownloader.DownloadAssets(assetURLs)
	// fmt.Println(localPaths)
	// replace the URLs in the HTML with the local paths
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
			// replace the URL in the HTML with the local path
			s.SetAttr("src", localPath)
			s.SetAttr("href", localPath)
		}
	})

	// save the modified HTML to a file
	html, err := doc.Html()
	if err != nil {
		fmt.Println("Error generating HTML:", err)
		return err
	}

	filename := path.Base(baseURL) + ".html"
	err = os.WriteFile(filename, []byte(html), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	fmt.Println("Done!")
	return nil
}
