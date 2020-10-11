package indexers

import (
	"log"
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const benzingaSource string = "benzinga"

func startBenzingaIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	go scraper.StartGetHTML("https://www.benzinga.com/", rate, func(body string) {
		onBenzingaIndexBody(es, body, scraper)
	})
	scraper.StartGetHTML("https://www.benzinga.com/pressreleases/", rate, func(body string) {
		onBenzingaPRBody(es, body, scraper)
	})
	return nil
}

func parseBenzingaArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	rg := regexp.MustCompile("<p[ \\w\"=;\\-:/\\.]*>([\\s\\S]+?)<\\/p>")
	matches := rg.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matches {
		paragraph := scraping.CleanHTMLText(match[1])
		paragraphs = append(paragraphs, paragraph)
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onBenzingaIndexBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("href=\"(\\/[\\w\\-]+\\/[\\w\\-]+\\/\\d+\\/\\d+\\/\\d+\\/[^\"]+?)\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://www.benzinga.com" + match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(benzingaSource, title, url, func(url string) string {
			return parseBenzingaArticle(url, scraper)
		})
	}
}

func onBenzingaPRBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("href=\"(\\/pressreleases\\/\\d+[^\"]+?)\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://www.benzinga.com" + match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(benzingaSource, title, url, func(url string) string {
			return parseBenzingaArticle(url, scraper)
		})
	}
}
