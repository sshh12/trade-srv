package indexers

import (
	"log"
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const reutersSource string = "reuters"

func startReutersIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	go scraper.StartGetHTML("https://www.reuters.com/finance/markets", rate, func(body string) {
		onReutersBody(es, body, scraper)
	})
	scraper.StartGetHTML("https://www.reuters.com/news/archive/domesticNews", rate, func(body string) {
		onReutersBody(es, body, scraper)
	})
	return nil
}

func parseReutersArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	rg := regexp.MustCompile("<p[ \\w\"\\-=:/\\.]*>([\\s\\S]+?)<\\/p>")
	matches := rg.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matches {
		paragraph := scraping.CleanHTMLText(match[1])
		paragraphs = append(paragraphs, paragraph)
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onReutersBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("href=\"([^\"]+?)\">\\s*?<h3 class=\"story-title\">\\s*([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://www.reuters.com" + match[1]
		title := scraping.CleanHTMLText(match[2])
		title = scraping.RegexReplace(title, "([A-Z\\d])-([A-Z])", "$1 - $2")
		es.OnEventArticleResolveBody(reutersSource, title, url, func(url string) string {
			return parseReutersArticle(url, scraper)
		})
	}
}
