package indexers

import (
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const finvizSource string = "finviz"
const finvizCalendarRegex = "<td align=\"right\">([\\d:APM ]+?)<[\\s\\S]+?\"left\">([\\w\\d,\\-\\. ]+?)<[\\s\\S]+?impact_(\\d+).gif[\\s\\S]+?ft\">([^<]+?)<[\\/td><\\n align=\"rh]+>([^<]+?)<[\\/td><\\n align=\"rh]+>([^<]+?)<[\\/td><\\n align=\"rh]+>([^<]+?)<"
const finvizSignalRegex = "tab-link\">([A-Z\\.]+?)<[\\S\\s]+?tab-link-nw\">([\\w ]+?)<"
const finvizHeadlineRegex = "]\"><a href=\"([^\"]+?)\" target=\"_blank\" class=\"nn-tab-link\">([^<]+?)<"

var finvizSignalWords = []string{
	"Top",
	"New",
	"Oversold",
	"Most",
	"Downgrades",
	"Insider",
	"Wedge",
	"Triangle",
	"Unusual",
	"Overbought",
	"Insider",
	"Channel",
	"Double",
	"Multiple",
}

func startFinVizIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	go scraper.StartGetHTML("https://finviz.com/calendar.ashx", rate, func(body string) {
		onFinVizCalendar(es, body, scraper)
	})
	go scraper.StartGetHTML("https://finviz.com", rate, func(body string) {
		onFinVizIndex(es, body, scraper)
	})
	scraper.StartGetHTML("https://finviz.com/news.ashx", rate, func(body string) {
		onFinVizNews(es, body, scraper)
	})
	return nil
}

func decodeFinVizImpact(idx string) string {
	if idx == "1" {
		return "low"
	} else if idx == "2" {
		return "medium"
	} else if idx == "3" {
		return "high"
	}
	return "unk"
}

func isFinVizSignal(sig string) bool {
	sigSplit := strings.Split(sig, " ")
	for _, sigWord := range finvizSignalWords {
		if sigWord == sigSplit[0] {
			return true
		}
	}
	return false
}

func onFinVizCalendar(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	bodyClean := scraping.RegexReplace(body, "<span style=\"color:#[\\w\\d]+;\">", "")
	bodyClean = scraping.RegexReplace(bodyClean, "<\\/span>", "")
	rg, _ := regexp.Compile(finvizCalendarRegex)
	matches := rg.FindAllStringSubmatch(bodyClean, -1)
	for _, match := range matches {
		reportTime := strings.TrimSpace(match[1])
		name := strings.TrimSpace(match[2])
		impactIdx := strings.TrimSpace(match[3])
		dateFor := strings.TrimSpace(match[4])
		actual := strings.TrimSpace(match[5])
		expected := strings.TrimSpace(match[6])
		prev := strings.TrimSpace(match[7])
		if actual == "-" {
			continue
		}
		evt := &events.Event{
			Source:        finvizSource,
			Type:          "economic_calendar",
			Name:          name,
			ActualValue:   actual,
			ExpectedValue: expected,
			PrevValue:     prev,
			Impact:        decodeFinVizImpact(impactIdx),
			TimeFor:       dateFor,
			TimeReported:  reportTime,
			CacheHash:     events.HashKey(name + dateFor),
		}
		es.OnEvent(evt)
	}
}

func onFinVizIndex(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg, _ := regexp.Compile(finvizSignalRegex)
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		sym := match[1]
		signal := match[2]
		if !isFinVizSignal(signal) {
			continue
		}
		date := time.Now().Format(time.RFC3339)[:10]
		evt := &events.Event{
			Source:           finvizSource,
			Type:             "technical_signal",
			Name:             signal,
			ConfirmedSymbols: []string{sym},
			TimeFor:          date,
			CacheHash:        events.HashKey(sym + signal + date),
		}
		es.OnEvent(evt)
	}
}

func onFinVizNews(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	rg, _ := regexp.Compile(finvizHeadlineRegex)
	matches := rg.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		url := strings.TrimSpace(match[1])
		title := strings.TrimSpace(match[2])
		evt := &events.Event{
			Source:    finvizSource,
			Type:      "article",
			Title:     title,
			URL:       url,
			CacheHash: events.HashKey(url),
		}
		es.OnEvent(evt)
	}
}
