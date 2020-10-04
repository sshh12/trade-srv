package index

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPGetScraper struct {
	URL    string
	OnBody func(string)
}

func NewHTTPGetScraper(url string) *HTTPGetScraper {
	return &HTTPGetScraper{URL: url}
}

func (hg *HTTPGetScraper) Get(url string) (string, error) {
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


func (hg *HTTPGetScraper) Start() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			body, err := hg.Get(hg.URL)
			if err != nil {
				fmt.Println(err)
			} else {
				hg.OnBody(body)
			}
		}
	}
}
