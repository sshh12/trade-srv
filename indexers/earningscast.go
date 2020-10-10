package indexers

import (
	"regexp"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const earningsCastSource string = "earningscast"

func startEarningsCastIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://earningscast.com/", rate, func(body string) {
		onEarningsCastIndex(es, body, scraper)
	})
	return nil
}

func onEarningsCastIndex(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("Act ([\\-\\.\\d]+) <br \\/>\\s*?Est ([\\-\\.\\d]+)\\s*?<\\/div>\\s*?<[\\w= \"]+?><a href=\"\\/([\\w\\.]+?)\\/\\d+\">([^<]+?)<")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		act := match[2]
		est := match[1]
		sym := match[3]
		title := scraping.CleanHTMLText(match[4])
		evt := &events.Event{
			Source:        earningsCastSource,
			Type:          "earnings",
			Title:         title,
			ActualValue:   act,
			ExpectedValue: est,
			CacheHash:     events.HashKey(sym + act + est + title),
		}
		es.OnEvent(evt)
	}
}
