package indexers

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const nasdaqSource string = "nasdaq"

func startNasdaqIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.nasdaq.com/news-and-insights", rate, func(body string) {
		onNasdaqBody(es, body, scraper)
	})
	return nil
}

func parseNasdaqArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	rg := regexp.MustCompile("<p>([\\s\\S]+?)<\\/p>")
	matches := rg.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matches {
		paragraph := scraping.CleanHTMLText(match[1])
		paragraphs = append(paragraphs, paragraph)
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onNasdaqBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("card-title-link\" href=\"(\\/articles[^\"]+?)\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://www.nasdaq.com" + match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(nasdaqSource, title, url, func(url string) string {
			return parseNasdaqArticle(url, scraper)
		})
	}
	rg2 := regexp.MustCompile("related-item__link\" href=\"(\\/articles[^\"]+?)\"[^<]+?<p[\\w =\"\\-]+?>([^<]+?)<")
	matches2 := rg2.FindAllStringSubmatch(body, -1)
	for _, match := range matches2 {
		url := "https://www.nasdaq.com" + match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(nasdaqSource, title, url, func(url string) string {
			return parseNasdaqArticle(url, scraper)
		})
	}
}
