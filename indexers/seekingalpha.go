package indexers

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const seekingAlphaSource string = "seekingalpha"

func startSeekingAlphaIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://seekingalpha.com/market-news", rate, func(body string) {
		onSeekingAlphaBody(es, body, scraper)
	})
	return nil
}

func parseSeekingAlphaArticle(url string, scraper *scraping.HTTPScraper) string {
	split := strings.Split(url, "/news/")
	path := "news"
	if len(split) != 2 {
		split = strings.Split(url, "/article/")
		path = "articles"
	}
	if len(split) != 2 {
		return ""
	}
	articleCode := strings.Split(split[1], "-")[0]
	apiURL := fmt.Sprintf("https://seekingalpha.com/api/v3/%s/%s?include=author%%2CprimaryTickers", path, articleCode)
	body, err := scraper.Get(apiURL)
	if err != nil {
		log.Println(err)
		return ""
	}
	var jsonBody map[string]interface{}
	if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
		log.Println(err)
		return ""
	}
	data, ok := jsonBody["data"].(map[string]interface{})
	if !ok {
		// Rate limited
		return ""
	}
	attr := data["attributes"].(map[string]interface{})
	content := attr["content"].(string)
	paragraphs := make([]string, 0)
	for _, line := range strings.Split(content, "<\\p>") {
		paragraphs = append(paragraphs, scraping.CleanHTMLText(line))
	}
	return strings.Join(paragraphs, "\n\n\n")
}

func onSeekingAlphaBody(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg, _ := regexp.Compile("href=\"([^\"]+?)\" class=\"[\\w-]+\" sasource=\"market_news\\w+\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := "https://seekingalpha.com" + match[1]
		title := scraping.CleanHTMLText(match[2])
		es.OnEventArticleResolveBody(seekingAlphaSource, title, url, func(url string) string {
			return parseSeekingAlphaArticle(url, scraper)
		})
	}
}
