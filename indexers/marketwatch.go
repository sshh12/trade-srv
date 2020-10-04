package indexers

import (
	"regexp"
	"strings"
	"fmt"

	scraping "github.com/sshh12/trade-srv/scraping"
	events "github.com/sshh12/trade-srv/events"
)

const source string = "marketwatch"

func startMarketWatchIndexer(es *events.EventStream, opts *IndexerOptions) error {
	scraper := scraping.NewHTTPGetScraper("https://www.marketwatch.com/latest-news")
	onGetContent := func(event *events.Event) {
		fmt.Println("lookup", event.URL)
		body, err := scraper.Get(event.URL)
		if err != nil {
			fmt.Println(err)
			return
		}
		rg, _ := regexp.Compile("<p>([\\s\\S]+?)<\\/p>")
		matches := rg.FindAllStringSubmatch(body, -1)
		for _, match := range matches {
			fmt.Println(match[1])
		}
	}
	scraper.OnBody = func(body string) {
		rg, _ := regexp.Compile("headline\"><a[ \"=\\w]+?href=\"(https:\\/\\/www.marketwatch.com[^\"]+?)\">([^<]+?)<")
		matches := rg.FindAllStringSubmatch(body, -1)
		for _, match := range matches {
			url := strings.ReplaceAll(match[1], "?mod=newsviewer_click", "")
			title := match[2]
			es.OnEventArticleResolveBody(source, title, url, onGetContent)
		}
	}
	scraper.Start()
	return nil
}
