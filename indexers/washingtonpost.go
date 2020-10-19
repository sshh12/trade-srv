package indexers

import (
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const washingtonPostSource string = "washingtonpost"

func startWashingtonPostIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.washingtonpost.com/", rate, func(body string) {
		onWashintonPostBody(es, body, scraper)
	})
	return nil
}

func parseWashingtonPostArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.WithField("source", washingtonPostSource).Error(err)
		return ""
	}
	rg := regexp.MustCompile("<p[ =\\w\"\\-]+?>([\\s\\S]+?)<\\/p>")
	matches := rg.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matches {
		paragraph := scraping.CleanHTMLText(match[1])
		paragraphs = append(paragraphs, paragraph)
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onWashintonPostBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("href=\"(https:\\/\\/www.washingtonpost.com[^\"]+?story.html)\"><span>([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(washingtonPostSource, title, url, func(url string) string {
			return parseWashingtonPostArticle(url, scraper)
		})
	}
}
