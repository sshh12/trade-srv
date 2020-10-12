package indexers

import (
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const blsTEDSource string = "blsted"

func startBLSTEDIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.bls.gov/opub/ted/year.htm", rate, func(body string) {
		onBLSTEDBody(es, body, scraper)
	})
	return nil
}

func parseBLSTEDArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.WithField("source", blsTEDSource).Error(err)
		return ""
	}
	rg := regexp.MustCompile("<p>([\\s\\S]+?)<\\/p>")
	matches := rg.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matches {
		paragraph := scraping.CleanHTMLText(match[1])
		if len(paragraph) < 40 {
			continue
		}
		paragraphs = append(paragraphs, paragraph)
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onBLSTEDBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("href=\"(\\/opub\\/ted\\/\\d+[^\"]+?)\"[^>]+?>([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://www.bls.gov" + match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(blsTEDSource, title, url, func(url string) string {
			return parseBLSTEDArticle(url, scraper)
		})
	}
}
