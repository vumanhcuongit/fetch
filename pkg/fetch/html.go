package fetch

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	netURL "net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseHTML(baseURL string, body io.ReadCloser) error {
	// parse the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println("Error parsing document:", err)
		return err
	}

	// download and replace the URLs of the assets
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
			url = baseURL + url
		}

		// download the asset
		localPath, err := downloadAsset(url)
		if err != nil {
			fmt.Println("Error downloading asset:", err)
			return
		}

		// replace the URL in the HTML with the local path
		s.SetAttr("src", localPath)
		s.SetAttr("href", localPath)
	})

	// save the modified HTML to a file
	html, err := doc.Html()
	if err != nil {
		fmt.Println("Error generating HTML:", err)
		return err
	}

	filename := path.Base(baseURL) + ".html"
	err = ioutil.WriteFile(filename, []byte(html), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	fmt.Println("Done!")
	return nil
}

func downloadAsset(urlStr string) (string, error) {
	// parse the URL
	u, err := netURL.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// Check if the file extension is allowed
	ext := path.Ext(u.Path)
	if ext != ".css" && ext != ".js" && ext != ".png" && ext != ".jpg" && ext != ".gif" {
		return "", nil
	}

	// build the local path
	localPath := path.Join("assets", u.Host, u.Path)
	if localPath[len(localPath)-1] == '/' {
		localPath += "index.html"
	}

	// create the directories
	err = os.MkdirAll(path.Dir(localPath), 0755)
	if err != nil {
		return "", err
	}

	// download the asset
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// save the asset to disk
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(localPath, body, 0644)
	if err != nil {
		return "", err
	}

	return localPath, nil
}
