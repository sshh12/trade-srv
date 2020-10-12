package indexers

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	events "github.com/sshh12/trade-srv/events"
	scraping "github.com/sshh12/trade-srv/scraping"
)

const bioPharmCatalystSource string = "biopharmcatalyst"

func startBioPharmCatalystIndexer(es *events.EventStream, opts *IndexerOptions) error {
	rate := opts.PollRate
	if rate == 0 {
		rate = 10 * time.Second
	}
	scraper := scraping.NewHTTPScraper()
	scraper.StartGetHTML("https://www.biopharmcatalyst.com/calendars/fda-calendar", rate, func(body string) {
		onFDACalendar(es, body, scraper)
	})
	return nil
}

func onFDACalendar(es *events.EventStream, body string, scraper *scraping.HTTPScraper) {
	bodyClean := strings.ReplaceAll(body, "&quot;", "\"")
	rg := regexp.MustCompile("<screener ticker-id=\"\" :pro=\"0\" :tabledata=\"([\\s\\S]+?)\">")
	matches := rg.FindAllStringSubmatch(bodyClean, -1)
	for _, match := range matches {
		var data [](map[string]interface{})
		if err := json.Unmarshal([]byte(match[1]), &data); err != nil {
			log.Println(err)
			continue
		}
		for _, item := range data {
			companies := item["companies"].(map[string]interface{})
			sym := companies["ticker"].(string)
			url := item["press_link"].(string)
			name := item["name"].(string)
			evt := &events.Event{
				Source:           bioPharmCatalystSource,
				Type:             "drug_update",
				Name:             name,
				URL:              url,
				ConfirmedSymbols: []string{sym},
				Extras:           item,
				CacheHash:        events.HashKey(name + url),
			}
			es.OnEvent(evt)
		}
	}
}
