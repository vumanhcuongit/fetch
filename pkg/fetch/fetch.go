package fetch

import (
	"fetch-go/pkg/utils"
	"fmt"
	"sync"
)

func Scrape(urls []string, getMetadata bool) error {
	var wg sync.WaitGroup
	wg.Add(len(urls))

	fmt.Println(urls)
	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			err := Fetch(url)
			if err != nil {
				fmt.Printf("error downloading %s: %v\n", url, err)
				return
			}
			fmt.Printf("downloaded %s\n", url)
		}(url)
	}

	wg.Wait()
	fmt.Println("all downloads complete.")
	return nil
}

func Fetch(url string) error {
	reader, err := utils.HttpGet(url)
	if err != nil {
		return err
	}

	return ParseHTML(url, reader)
}
