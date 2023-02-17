package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

// create a shared HTTP client for connection reuse
var httpClient *http.Client
var once sync.Once

func GetHTTPClient() *http.Client {
	once.Do(func() {
		transport := &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		}
		httpClient = &http.Client{
			Timeout:   time.Second * 30,
			Transport: transport,
		}
	})
	return httpClient
}

func HttpGet(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	resp, err := GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP response status code was %d", resp.StatusCode)
	}

	bodyCopy := &bytes.Buffer{}
	_, err = io.Copy(bodyCopy, resp.Body)
	if err != nil {
		return nil, err
	}

	// decode the response if it's compressed
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(bodyCopy)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	case "deflate":
		reader = flate.NewReader(bodyCopy)
		defer reader.Close()
	default:
		reader = ioutil.NopCloser(bodyCopy)
	}

	return reader, nil
}
