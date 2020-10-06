package scraping

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPScraper struct {
}

func NewHTTPScraper() *HTTPScraper {
	return &HTTPScraper{}
}

func (hs *HTTPScraper) Get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (hs *HTTPScraper) StartGetHTML(url string, onBody func(string)) {
	ticker := time.NewTicker(5 * time.Second)
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
