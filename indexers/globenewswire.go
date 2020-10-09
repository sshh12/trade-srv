package indexers

import (
	"log"
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const globalnewswireSource string = "globalnewswire"

func startGlobalNewsWireIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.globenewswire.com/Index", rate, func(body string) {
		onGlobalNewsWireBody(es, body, scraper)
	})
	return nil
}

func parseGlobalNewsWireArticle(url string, scraper *scraping.HTTPScraper) string {
	body, err := scraper.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	rg, _ := regexp.Compile("itemprop=\"articleBody\">([\\s\\S]+?)<\\/span>")
	matches := rg.FindAllStringSubmatch(body, -1)
	paragraphs := make([]string, 0)
	for _, match := range matches {
		for _, group := range strings.Split(match[1], "<br /><br />") {
			paragraph := scraping.CleanHTMLText(strings.ReplaceAll(group, "<br />", "\n"))
			paragraphs = append(paragraphs, paragraph)
		}
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onGlobalNewsWireBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg, _ := regexp.Compile("title\">([^<]+?)<\\/p>[\\s\\S]+?title16px\"><a href=\"([^\"]+?)\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://www.globenewswire.com" + match[2]
		company := match[1]
		title := scraping.CleanHTMLText(company + " -- " + match[3])
		es.OnEventArticleResolveBody(globalnewswireSource, title, url, func(url string) string {
			return parseGlobalNewsWireArticle(url, scraper)
		})
	}
}
