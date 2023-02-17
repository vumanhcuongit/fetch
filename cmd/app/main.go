package main

import (
	"fetch-go/pkg/fetch"
	"fmt"
	"os"
	"sync"
)

func main() {
	// get the URLs from the command-line arguments
	urls := os.Args[1:]

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			err := fetch.Fetch(url)
			if err != nil {
				fmt.Printf("error downloading %s: %v\n", url, err)
				return
			}
			fmt.Printf("downloaded %s\n", url)
		}(url)
	}

	wg.Wait()
	fmt.Println("all downloads complete.")
}
