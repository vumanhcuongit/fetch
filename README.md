## Archiving Assets
When running the CLI with the makefile command, all assets (such as CSS, JS, and image files) are automatically downloaded and archived locally. This means that the web page can be loaded locally with all assets loading properly, providing a convenient way to view the web page offline or on a local network.

To access the archived assets, navigate to the `assets` folder that is created in the root directory after running the makefile command. Inside this folder, you will find a subfolder for each URL that was fetched, containing all of the downloaded assets for that URL.

## Concurrency and Performance
To improve the performance of URL fetching, concurrency features of Go such as Goroutines and channels are utilized. Instead of downloading URLs sequentially, which can cause blocking and longer wait times, they are downloaded concurrently. This allows multiple downloads to occur simultaneously, leading to better performance.

Similarly, assets such as CSS, JS, and image files are downloaded concurrently instead of sequentially to further improve performance. For each file to download, a new Goroutine is created to perform the download. A WaitGroup is utilized to ensure that all Goroutines complete their tasks before the program exits, and the number of Goroutines is limited to a reasonable number to avoid overloading the system.

By controlling the maximum number of Goroutines, optimal program performance is maintained even when downloading a large number of files. This approach results in faster response times and an improved user experience.

## Gzip Handling
To improve network performance and reduce the amount of data that needs to be transferred over the network, the application uses gzip compression for HTTP responses. To handle gzip-encoded content, we use the `gzip` package provided by Go standard library.

The HttpGet(url string) function makes an HTTP GET request to the given URL and sets the `Accept-Encoding` header to `gzip, deflate` to indicate that the application can handle compressed content.

By handling gzip-encoded content, the amount of data that needs to be transferred over the network is reduced, leading to faster response times and better user experience.

## Reusing HTTP Connections
One of the key ways to improve the performance of the application that downloads a large number of assets, such as CSS, JS, and image files, is to reuse http connections. This way, we can avoid the overhead of establishing a new TCP connection for each file to download.

To reuse connections in Go, we can create an `http.Client` with a custom `Transport` that has a connection pool. This way, we can reuse connections for multiple requests to the same host.

Let's take a look at an example of how this can be implemented in the app:

In the `GetHTTPClient` function, a singleton instance of the `http.Client` is created with a `Transport` that has a pool of connections. The `http.Client` is returned whenever `GetHTTPClient` is called.

```go
package utils

import (
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
```

## Running the CLI with makefile
1. To fetch a single URL, run the following command:

```bash
make fetch URLS=<URL>
```

2. To fetch multiple URLs, run the following command:

```bash
make fetch URLS="<URL1> <URL2> <URL3>"
```
Replace `<URL1>`, `<URL2>`, and `<URL3>` with the URLs that you want to fetch, **separated by a whitespace**.

3. If you want to include metadata in the output, add the `METADATA=true` argument:

```bash
make fetch URLS=<URL> METADATA=true
```

Example usage:

```bash
make fetch URLS="https://mholt.github.io/json-to-go" METADATA=true
```

## CLI Framework
The CLI tool is built using the [Cobra](https://github.com/spf13/cobra) library, which provides a powerful framework for building modern CLI applications in Go.
The `rootCmd` command takes two optional flags: `metadata` and `max_concurrent`. The `metadata` flag is used to indicate whether metadata should be included in the output, and the `max_concurrent` flag is used to set the maximum number of concurrent requests that can be made. The default value for `max_concurrent` is set to 100.