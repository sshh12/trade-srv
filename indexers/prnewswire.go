package indexers

import (
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const prNewsWireSource string = "prnewswire"

func startPrNewsWireIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.prnewswire.com/news-releases/news-releases-list/", rate, func(body string) {
		onPrNewsWireBody(es, body, scraper)
	})
	return nil
}

func parsePrNewsWireArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.WithField("source", prNewsWireSource).Error(err)
		return ""
	}
	rg, _ := regexp.Compile("<article class=\"news-release carousel-template\">([\\s\\S]+?)<\\/article>")
	matches := rg.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matches {
		for _, group := range strings.Split(match[1], "</p>") {
			paragraph := scraping.CleanHTMLText(strings.ReplaceAll(group, "<br />", "\n"))
			paragraphs = append(paragraphs, paragraph)
		}
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onPrNewsWireBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("news-release\" href=\"([^\"]+?)\" title=\"[^\"]*?\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://www.prnewswire.com" + match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(prNewsWireSource, title, url, func(url string) string {
			return parsePrNewsWireArticle(url, scraper)
		})
	}
}
