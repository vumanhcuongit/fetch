## Concurrency and Performance

To improve the performance of our URL fetching and asset downloading, we use Go's concurrency features, such as Goroutines and channels. Instead of downloading URLs sequentially, which can result in blocking and longer wait times, we download them concurrently. This means that multiple downloads can occur at the same time, leading to better performance.

Similarly, we download assets such as CSS, JS, and image files concurrently rather than sequentially to further improve performance. For each file we want to download, we create a new Goroutine to perform the download. We use a WaitGroup to ensure that all Goroutines complete their tasks before the program exits, and we also limit the number of Goroutines to a reasonable number to avoid overloading the system.

By controlling the maximum number of Goroutines, we can ensure that the program performs well even when downloading a large number of files. This approach leads to faster response times and a better user experience.

## Gzip Handling

To improve network performance and reduce the amount of data that needs to be transferred over the network, the application uses gzip compression for HTTP responses. To handle gzip-encoded content, we use the `gzip` package provided by Go standard library.

The HttpGet(url string) function makes an HTTP GET request to the given URL and sets the `Accept-Encoding` header to `gzip, deflate` to indicate that the application can handle compressed content.

By handling gzip-encoded content, the amount of data that needs to be transferred over the network is reduced, leading to faster response times and better user experience.

## Running the CLI with Makefile

1. To fetch a single URL, run the following command:

```bash
make fetch URLS=<URL>
```

2. To fetch multiple URLs, run the following command:

```bash
make fetch URLS="<URL1> <URL2> <URL3>"
```
Replace <URL1>, <URL2>, and <URL3> with the URLs that you want to fetch, **separated by a space**.

3. If you want to include metadata in the output, add the `METADATA=true` argument:

```bash
make fetch URLS=<URL> METADATA=true
```

Example usage:

```bash
make fetch URLS="https://www.google.com https://mholt.github.io/json-to-go" METADATA=true
```
