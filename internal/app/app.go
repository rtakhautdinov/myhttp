package app

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultTimeoutInSeconds default timeout that client could wait while wait for response
	DefaultTimeoutInSeconds = 10
	// DefaultScheme default scheme for url is it is not specified in request
	DefaultScheme = "https://"
)

// HTTPClient interface
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// WorkerResult this structure to specify format for worker's
// result
type WorkerResult struct {
	Url       string
	IsSuccess bool
	Md5Data   string
}

func prepareUrl(link string) (string, error) {
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = DefaultScheme + link
	}

	if _, err := url.ParseRequestURI(link); err != nil {
		return "", err
	}

	return link, nil
}

// Worker is computation unit that could be run in parallel with other workers
// to get url md5 calculation concurrently
func Worker(client HTTPClient, urls chan string, results chan WorkerResult) {
	for u := range urls {
		u, err := prepareUrl(u)

		if err != nil {
			results <- WorkerResult{
				u,
				false,
				"",
			}

			continue
		}

		resp, err := client.Get(u)
		if err != nil {
			results <- WorkerResult{
				u,
				false,
				"",
			}

			continue
		}

		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			results <- WorkerResult{
				u,
				false,
				"",
			}

			continue
		}

		results <- WorkerResult{
			u,
			true,
			fmt.Sprintf("%x", md5.Sum(data)),
		}
	}
}

// Run just simple endpoint to this tool that accept
// list/slice of urls to be processed and amount of parallel
// workers
func Run(urls []string, numOfParallel int) {
	client := &http.Client{
		Timeout: DefaultTimeoutInSeconds * time.Second,
	}

	buf := make(chan string, numOfParallel)
	results := make(chan WorkerResult)

	for i := 0; i < cap(buf); i++ {
		go Worker(client, buf, results)
	}

	go func() {
		for _, u := range urls {
			buf <- u
		}
	}()

	for i := 0; i < len(urls); i++ {
		urlData := <-results

		if urlData.IsSuccess {
			fmt.Printf("%s %s\n", urlData.Url, urlData.Md5Data)
		} else {
			fmt.Printf("%s ERROR\n", urlData.Url)
		}
	}
}
