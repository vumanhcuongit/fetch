package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	resp, err := http.DefaultClient.Do(req)
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
