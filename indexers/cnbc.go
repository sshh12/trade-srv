package indexers

import (
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const cnbcSource string = "cnbc"

func startCNBCIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	go scraper.StartGetHTML("https://www.cnbc.com/business/", rate, func(body string) {
		onCNBCBody(es, body, scraper)
	})
	go scraper.StartGetHTML("https://www.cnbc.com/economy/", rate, func(body string) {
		onCNBCBody(es, body, scraper)
	})
	go scraper.StartGetHTML("https://www.cnbc.com/technology/", rate, func(body string) {
		onCNBCBody(es, body, scraper)
	})
	scraper.StartGetHTML("https://www.cnbc.com/markets/", rate, func(body string) {
		onCNBCBody(es, body, scraper)
	})
	return nil
}

func parseCNBCArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.WithField("source", cnbcSource).Error(err)
		return ""
	}
	rgList := regexp.MustCompile("<li>([\\s\\S]+?)<\\/li>")
	matchesList := rgList.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matchesList {
		paragraph := scraping.CleanHTMLText(match[1])
		if len(paragraph) < 40 {
			continue
		}
		paragraphs = append(paragraphs, paragraph)
	}
	rg := regexp.MustCompile("<p>([\\s\\S]+?)<\\/p>")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		paragraph := scraping.CleanHTMLText(match[1])
		if len(paragraph) < 40 {
			continue
		}
		paragraphs = append(paragraphs, paragraph)
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onCNBCBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("<a href=\"(https:\\/\\/www.cnbc.com\\/\\d+\\/\\d+\\/\\d+\\/[^\"]+?)\" class=\"Card-title\" target=\"\"><div>([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(cnbcSource, title, url, func(url string) string {
			return parseCNBCArticle(url, scraper)
		})
	}
}
