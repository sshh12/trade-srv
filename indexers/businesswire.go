package indexers

import (
	"log"
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const businessWireSource string = "businesswire"

func startBusinessWireSourceIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.businesswire.com/portal/site/home/news/", rate, func(body string) {
		onBusinessWireBody(es, body, scraper)
	})
	return nil
}

func parseBusinessWireArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	rg, _ := regexp.Compile("<p>([\\s\\S]+?)<\\/p>")
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

func onBusinessWireBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("<a class=\"bwTitleLink\"\\s*href=\"([^\"]+?)\"><span itemprop=\"headline\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://www.businesswire.com" + match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(businessWireSource, title, url, func(url string) string {
			return parseBusinessWireArticle(url, scraper)
		})
	}
}
