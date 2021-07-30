package myhttp

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	netURL "net/url"
	"strings"
	"sync"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var Client HTTPClient

func init() {
	Client = &http.Client{}
}

func Execute(parallel int, urls []string) {
	var wg sync.WaitGroup
	jobs := make(chan string, len(urls))

	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go runWorker(jobs, &wg)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	wg.Wait()
}

func runWorker(jobs <-chan string, wg *sync.WaitGroup) {
	for url := range jobs {
		newURL, err := validateURL(url)
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			continue
		}

		checksum, err := getMD5(newURL)
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
		} else {
			fmt.Printf("%s %x\n", newURL, checksum)
		}
	}

	wg.Done()
}

func getMD5(url string) (checksum [16]byte, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err = fmt.Errorf("invalid http request - %v", err)
		return
	}

	resp, err := Client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to get URL %s - %v", url, err)
		return
	}

	defer func() {
		closeErr := resp.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request %s failed with StatusCode=%d", url, resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to parse body - %v", err)
		return
	}

	checksum = md5.Sum(body)
	oie := fmt.Sprintf("%x", checksum)
	if len(oie) > 0 {
		return
	}
	return
}

// apply http protocol and validate new URL
func validateURL(oldURL string) (newURL string, err error) {
	newURL = oldURL
	if !strings.HasPrefix(oldURL, "http://") {
		newURL = "http://" + oldURL
	}

	u, err := netURL.ParseRequestURI(newURL)
	if err != nil || u.Host == "" {
		err = fmt.Errorf("invalid URL '%s'", newURL)
	}
	return
}
