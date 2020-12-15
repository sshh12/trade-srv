package scraping

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"

// HTTPScraper is a HTTP-based webscraper
type HTTPScraper struct {
	client *http.Client
}

// NewHTTPScraper creates a new HTTPScraper
func NewHTTPScraper() *HTTPScraper {
	return &HTTPScraper{client: &http.Client{}}
}

// Get does a HTTP Get request
func (hs *HTTPScraper) Get(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"86\", \"\"Not\\A;Brand\";v=\"99\", \"Google Chrome\";v=\"86\"")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	if err != nil {
		return "", err
	}
	resp, err := hs.client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return FixForUTF8(body), nil
}

// StartGetHTML starts a loop of sending GET requests to the given url
func (hs *HTTPScraper) StartGetHTML(url string, rate time.Duration, onBody func(string)) {
	ticker := time.NewTicker(rate)
	for {
		select {
		case <-ticker.C:
			body, err := hs.Get(url)
			if err != nil {
				fmt.Println(err)
			} else {
				onBody(body)
			}
		}
	}
}
