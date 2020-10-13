package indexers

import (
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const forbesSource string = "forbes"

func startForbesIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.forbes.com/", rate, func(body string) {
		onForbesBody(es, body, scraper)
	})
	return nil
}

func parseForbesArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.WithField("source", forbesSource).Error(err)
		return ""
	}
	paragraphs := make([]string, 0)
	rg := regexp.MustCompile("<p>([\\s\\S]+?)<\\/p>")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		paragraph := scraping.CleanHTMLText(match[1])
		paragraphs = append(paragraphs, paragraph)
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onForbesBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile(" href=\"(https:\\/\\/www.forbes.com\\/sites[^\"]+?)\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(forbesSource, title, url, func(url string) string {
			return parseForbesArticle(url, scraper)
		})
	}
}
