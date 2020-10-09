package indexers

import (
	"log"
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const barronsSource string = "barrons"

func startBarronsIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.barrons.com/real-time?mod=hp_LATEST&mod=hp_LATEST", rate, func(body string) {
		onBarronsBody(es, body, scraper)
	})
	return nil
}

func parseBarronsArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	rg := regexp.MustCompile("<p>([\\s\\S]+?)<\\/p>")
	matches := rg.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matches {
		paragraph := scraping.CleanHTMLText(strings.ReplaceAll(match[1], "<br />", "\n"))
		paragraphs = append(paragraphs, paragraph)
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onBarronsBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("headline-link--[\\w\\d ]+\" href=\"([^\"]+?)\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := strings.ReplaceAll(match[1], "?mod=RTA", "")
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(barronsSource, title, url, func(url string) string {
			return parseBarronsArticle(url, scraper)
		})
	}
}
