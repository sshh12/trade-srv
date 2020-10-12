package indexers

import (
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const federalReserveSource string = "federalreserve"

func startFederalReserveIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.federalreserve.gov/json/ne-press.json", rate, func(body string) {
		onFederalReserveBody(es, body, scraper)
	})
	return nil
}

func parseFederalReserveArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.WithField("source", federalReserveSource).Error(err)
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

func onFederalReserveBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("\"t\":\"([^\"]+?)\"[\\s\\S]+?\"l\":\"(\\/[^\"]+?)\"")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://www.federalreserve.gov" + match[2]
		title := scraping.CleanHTMLText(match[1])
		es.OnEventArticleResolveBody(federalReserveSource, title, url, func(url string) string {
			return parseFederalReserveArticle(url, scraper)
		})
	}
}
