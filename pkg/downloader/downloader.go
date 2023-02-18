package downloader

import (
	"fetch-go/pkg/utils"
	"io"
	"net/http"
	netURL "net/url"
	"os"
	"path"
	"sync"

	"go.uber.org/zap"
)

var whitelistExtensions = []string{".css", ".js", ".jpg", ".jpeg", ".png", ".gif", ".bmp"}

// Asset represents a downloaded asset.
type Asset struct {
	URL       string
	LocalPath string // the local path of the downloaded asset
	Error     error
}

// AssetDownloader holds the data for the downloading process.
type AssetDownloader struct {
	logger        *zap.SugaredLogger
	client        *http.Client
	guardChan     chan struct{}   // the channel used for limiting the download goroutines
	wg            *sync.WaitGroup // a WaitGroup used to wait for all download goroutines to complete
	maxConcurrent int             // the maximum number of concurrent download goroutines
}

func NewAssetDownloader(logger *zap.SugaredLogger, client *http.Client, maxConcurrent int) *AssetDownloader {
	return &AssetDownloader{
		logger:        logger,
		client:        client,
		guardChan:     make(chan struct{}, maxConcurrent),
		wg:            &sync.WaitGroup{},
		maxConcurrent: maxConcurrent,
	}
}

// DownloadAsset downloads an asset and returns its local path.
func (d *AssetDownloader) DownloadAsset(url string) (*Asset, error) {
	// parse the URL
	u, err := netURL.Parse(url)
	if err != nil {
		d.logger.Errorf("failed to parse URL +%v: %+v", url, err)
		return nil, err
	}

	// check if the file extension is allowed
	ext := path.Ext(u.Path)
	isFileKindsupported := false
	for _, whitelistExtension := range whitelistExtensions {
		if ext == whitelistExtension {
			isFileKindsupported = true
		}
	}
	if !isFileKindsupported {
		return &Asset{}, nil
	}

	// build the local path
	localPath := path.Join("assets", u.Host, u.Path)
	if localPath[len(localPath)-1] == '/' {
		localPath += ".html"
	}

	// create the directories
	err = os.MkdirAll(path.Dir(localPath), 0755)
	if err != nil {
		d.logger.Errorf("failed to mkdirall, err: %+v", err)
		return nil, err
	}

	// download the asset
	reader, err := utils.HttpGet(url)
	if err != nil {
		d.logger.Errorf("failed to get content from URL, err: %+v", err)
		return nil, err
	}

	// save the asset to disk
	body, err := io.ReadAll(reader)
	if err != nil {
		d.logger.Errorf("failed to read content, err: %+v", err)
		return nil, err
	}
	err = os.WriteFile(localPath, body, 0644)
	if err != nil {
		d.logger.Errorf("failed to save the asset to disk, err: %+v", err)
		return nil, err
	}

	return &Asset{URL: url, LocalPath: localPath}, nil
}

// DownloadAssets downloads all the assets in the given URLs and returns a slice of the local paths.
func (d *AssetDownloader) DownloadAssets(urls []string) []string {
	localPaths := make([]string, 0, len(urls))
	d.wg.Add(len(urls))
	for i, url := range urls {
		// would block if guard channel is already filled
		d.guardChan <- struct{}{}
		go func(url string, i int) {
			defer d.wg.Done()
			asset, err := d.DownloadAsset(url)
			if err == nil && asset.LocalPath != "" {
				localPaths = append(localPaths, asset.LocalPath)
			}
			<-d.guardChan
		}(url, i)
	}

	d.wg.Wait()
	return localPaths
}
