package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	events "github.com/sshh12/trade-srv/events"
	indexers "github.com/sshh12/trade-srv/indexers"
	"github.com/sshh12/trade-srv/scraping"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	pgURL := flag.String("pg_url", "", "Postgres url (use instead of individual pg options)")
	pgUser := flag.String("pg_user", "postgres", "Postgres username")
	pgPassword := flag.String("pg_pass", "password", "Postgres password")
	pgAddr := flag.String("pg_addr", "localhost:5432", "Postgres host and port")
	pgName := flag.String("pg_db", "tradesrv", "Postgres database name")
	twKey := flag.String("tw_key", "", "Twitter consumer key")
	twSecret := flag.String("tw_secret", "", "Twitter consumer secret")
	twToken := flag.String("tw_token", "", "Twitter access token")
	twTokenSecret := flag.String("tw_token_secret", "", "Twitter access token secret")
	twNames := flag.String("tw_names", "", "Twitter accounts to follow")
	indexersSelected := make(map[string]*bool)
	for name := range indexers.EventIndexers {
		indexersSelected[name] = flag.Bool("run_"+name, false, "Run "+name+" indexer")
	}
	runAllEvents := flag.Bool("run_all_events", false, "Run all event indexers")
	warmUp := flag.Int("warmup", 120, "Discard events that occur in this number of seconds")
	addSymbol := flag.String("add_sym", "", "Register symbol(s) in database")
	flag.Parse()
	if *runAllEvents {
		for name := range indexers.EventIndexers {
			indexersSelected[name] = runAllEvents
		}
	}
	if *warmUp < 0 {
		*warmUp = 0
	}
	var db *events.Database
	var err error
	var postgresName string
	if *pgURL != "" {
		db, err = events.NewPostgresDatabaseFromURL(*pgURL)
		postgresName = *pgURL
	} else {
		db, err = events.NewPostgresDatabase(*pgUser, *pgPassword, *pgAddr, *pgName)
		postgresName = "postgres://" + *pgAddr
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Info("Connected to " + postgresName)
	if *addSymbol != "" {
		for _, sym := range strings.Split(*addSymbol, ",") {
			registerSymbol(sym, db)
		}
	}
	es := events.NewEventStream(db, time.Duration(*warmUp)*time.Second)
	opts := &indexers.IndexerOptions{
		PollRate:              0,
		TwitterConsumerKey:    *twKey,
		TwitterConsumerSecret: *twSecret,
		TwitterAccessToken:    *twToken,
		TwitterAccessSecret:   *twTokenSecret,
		TwitterNames:          strings.Split(*twNames, ","),
	}
	indexersRunning := make([]string, 0)
	for name, indexer := range indexers.EventIndexers {
		if *indexersSelected[name] {
			indexersRunning = append(indexersRunning, name)
			go indexer(es, opts)
		}
	}
	if len(indexersRunning) > 0 {
		log.Infof("Running %v", indexersRunning)
		for {
		}
	}
}

func registerSymbol(sym string, db *events.Database) {
	symClean := strings.TrimSpace(strings.ToUpper(sym))
	scraper := scraping.NewHTTPScraper()
	resp, err := scraper.Get(fmt.Sprintf("https://www.marketwatch.com/investing/stock/%s/profile", symClean))
	if err != nil {
		log.Error(err)
		return
	}
	nameRg := regexp.MustCompile("class=\"company__name\">([^<]+?)<")
	nameMatch := nameRg.FindStringSubmatch(resp)
	indRg := regexp.MustCompile("Industry<\\/small>\\s*?<span class=\"primary\\s*\">([^<]+?)<")
	indMatch := indRg.FindStringSubmatch(resp)
	secRg := regexp.MustCompile("Sector<\\/small>\\s*?<span class=\"primary\\s*\">([^<]+?)<")
	secMatch := secRg.FindStringSubmatch(resp)
	if len(nameMatch) == 0 || len(indMatch) == 0 || len(secMatch) == 0 {
		log.Error(symClean + " lookup failed")
		return
	}
	symbol := &events.Symbol{Sym: symClean, Name: nameMatch[1], Sector: secMatch[1], Industry: indMatch[1]}
	if err := db.AddSymbol(symbol); err != nil {
		log.Error(symClean + " registration failed")
	} else {
		log.Info(symClean + " added")
	}
}
