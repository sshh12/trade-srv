package indexers

import (
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const cnnSource string = "cnn"

func startCNNIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	go scraper.StartGetHTML("https://money.cnn.com/data/markets/", rate, func(body string) {
		onCNNBody(es, body, scraper)
	})
	go scraper.StartGetHTML("https://www.cnn.com/business/tech", rate, func(body string) {
		onCNNBody(es, body, scraper)
	})
	go scraper.StartGetHTML("https://www.cnn.com/world", rate, func(body string) {
		onCNNBody(es, body, scraper)
	})
	scraper.StartGetHTML("https://www.cnn.com/business", rate, func(body string) {
		onCNNBody(es, body, scraper)
	})
	return nil
}

func parseCNNArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.WithField("source", cnnSource).Error(err)
		return ""
	}
	rg := regexp.MustCompile("<div class=\"zn-body__paragraph[\\w ]*\">([\\s\\S]+?)<\\/div>")
	matches := rg.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matches {
		paragraph := scraping.CleanHTMLText(match[1])
		paragraphs = append(paragraphs, paragraph)
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onCNNBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("\"uri\":\"([^\"]+?)\",\"headline\":\"([^\"]+?)\"")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://www.cnn.com" + match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(cnnSource, title, url, func(url string) string {
			return parseCNNArticle(url, scraper)
		})
	}
	rg2 := regexp.MustCompile("<a href=\"(\\/\\d+\\/\\d+\\/[^\"]+?)\"\\s*><span class=\"cd__headline-text\">([^<]+?)<\\/span>")
	matches2 := rg2.FindAllStringSubmatch(body, -1)
	for _, match := range matches2 {
		url := "https://www.cnn.com" + match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(cnnSource, title, url, func(url string) string {
			return parseCNNArticle(url, scraper)
		})
	}
}
