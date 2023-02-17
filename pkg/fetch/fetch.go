package fetch

import (
	"fmt"
	"net/http"
)

func Fetch(baseURL string) error {
	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Error while getting the response:", err)
		return err
	}
	defer resp.Body.Close()

	return ParseHTML(baseURL, resp.Body)
}
