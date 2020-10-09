package scraping

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"

type HTTPScraper struct {
	client *http.Client
}

func NewHTTPScraper() *HTTPScraper {
	return &HTTPScraper{client: &http.Client{}}
}

func (hs *HTTPScraper) Get(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", defaultUserAgent)
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
	return string(body), nil
}

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
