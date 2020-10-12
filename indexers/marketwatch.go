package indexers

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const marketWatchSource string = "marketwatch"

func startMarketWatchIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.marketwatch.com/latest-news", rate, func(body string) {
		onMarketWatchBody(es, body, scraper)
	})
	return nil
}

func parseMarketWatchArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	rg, _ := regexp.Compile("<p>([\\s\\S]+?)<\\/p>")
	matches := rg.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matches {
		paragraph := scraping.CleanHTMLText(match[1])
		if len(paragraph) > 30 {
			paragraphs = append(paragraphs, paragraph)
		}
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onMarketWatchBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg, _ := regexp.Compile("headline\"><a[ \"=\\w]+?href=\"(https:\\/\\/www.marketwatch.com[^\"]+?)\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := strings.ReplaceAll(match[1], "?mod=newsviewer_click", "")
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(marketWatchSource, title, url, func(url string) string {
			return parseMarketWatchArticle(url, scraper)
		})
	}
}
