## Concurrency and Performance

To improve the performance of URL fetching, concurrency features of Go such as Goroutines and channels are utilized. Instead of downloading URLs sequentially, which can cause blocking and longer wait times, they are downloaded concurrently. This allows multiple downloads to occur simultaneously, leading to better performance.

Similarly, assets such as CSS, JS, and image files are downloaded concurrently instead of sequentially to further improve performance. For each file to download, a new Goroutine is created to perform the download. A WaitGroup is utilized to ensure that all Goroutines complete their tasks before the program exits, and the number of Goroutines is limited to a reasonable number to avoid overloading the system.

By controlling the maximum number of Goroutines, optimal program performance is maintained even when downloading a large number of files. This approach results in faster response times and an improved user experience.

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
Replace `<URL1>`, `<URL2>`, and `<URL3>` with the URLs that you want to fetch, **separated by a space**.

3. If you want to include metadata in the output, add the `METADATA=true` argument:

```bash
make fetch URLS=<URL> METADATA=true
```

Example usage:

```bash
make fetch URLS="https://www.google.com https://mholt.github.io/json-to-go" METADATA=true
```
