package indexers

import (
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const zacksSource string = "zacks"

func startZacksIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.zacks.com/research/earnings/z2_earnings_tab_data.php", rate, func(body string) {
		onZacksEarningsTabData(es, body, scraper)
	})
	return nil
}

func onZacksEarningsTabData(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg := regexp.MustCompile("symbol\\\\\">(\\w+)<\\/span>[\\s\\S]+?\", \"(\\d+:\\d+)\", \"([\\-NA\\d\\.]+?)\", \"([\\-NA\\d\\.]+?)\"")
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		sym := strings.TrimSpace(match[1])
		time := strings.TrimSpace(match[2])
		est := strings.TrimSpace(match[3])
		act := strings.TrimSpace(match[4])
		if act == "--" || act == "NA" || act == "" {
			continue
		}
		evt := &events.Event{
			Source:           zacksSource,
			Type:             "earnings",
			ActualValue:      act,
			ExpectedValue:    est,
			TimeReported:     time,
			ConfirmedSymbols: []string{sym},
			CacheHash:        events.HashKey(zacksSource + sym + act),
		}
		es.OnEvent(evt)
	}
}
