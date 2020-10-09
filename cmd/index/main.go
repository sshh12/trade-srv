package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	events "github.com/sshh12/trade-srv/events"
	indexers "github.com/sshh12/trade-srv/indexers"
	"github.com/sshh12/trade-srv/scraping"
)

func main() {
	pgUser := flag.String("pg_user", "postgres", "Postgres username")
	pgPassword := flag.String("pg_pass", "password", "Postgres password")
	pgAddr := flag.String("pg_addr", "localhost:5432", "Postgres host and port")
	pgName := flag.String("pg_db", "tradesrv", "Postgres database name")
	indexersSelected := make(map[string]*bool)
	for name := range indexers.AllIndexers {
		indexersSelected[name] = flag.Bool("run_"+name, false, "Run "+name+" indexer")
	}
	runAll := flag.Bool("run_all", false, "Run all indexers")
	warmUp := flag.Int("warmup", 60, "Discard events that occur in this number of seconds")
	addSymbol := flag.String("add_sym", "", "Register symbol(s) in database")
	flag.Parse()
	if *runAll {
		for name := range indexers.AllIndexers {
			indexersSelected[name] = runAll
		}
	}
	if *warmUp < 0 {
		*warmUp = 0
	}
	db, err := events.NewPostgresDatabase(*pgUser, *pgPassword, *pgAddr, *pgName)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Connected to postgres://" + *pgAddr)
	if *addSymbol != "" {
		for _, sym := range strings.Split(*addSymbol, ",") {
			registerSymbol(sym, db)
		}
	}
	es := events.NewEventStream(db, time.Duration(*warmUp)*time.Second)
	indexerRunning := false
	for name, indexer := range indexers.AllIndexers {
		if *indexersSelected[name] {
			log.Println("Starting " + name)
			indexerRunning = true
			go indexer(es, &indexers.IndexerOptions{})
		}
	}
	if indexerRunning {
		for {
		}
	}
}

func registerSymbol(sym string, db *events.Database) {
	symClean := strings.TrimSpace(strings.ToUpper(sym))
	scraper := scraping.NewHTTPScraper()
	resp, err := scraper.Get(fmt.Sprintf("https://www.marketwatch.com/investing/stock/%s/profile", symClean))
	if err != nil {
		log.Fatal(err)
		return
	}
	nameRg := regexp.MustCompile("class=\"company__name\">([^<]+?)<")
	nameMatch := nameRg.FindStringSubmatch(resp)
	indRg := regexp.MustCompile("Industry<\\/small>\\s*?<span class=\"primary\\s*\">([^<]+?)<")
	indMatch := indRg.FindStringSubmatch(resp)
	secRg := regexp.MustCompile("Sector<\\/small>\\s*?<span class=\"primary\\s*\">([^<]+?)<")
	secMatch := secRg.FindStringSubmatch(resp)
	if len(nameMatch) == 0 || len(indMatch) == 0 || len(secMatch) == 0 {
		log.Fatal(symClean + " lookup failed")
		return
	}
	symbol := &events.Symbol{Sym: symClean, Name: nameMatch[1], Sector: secMatch[1], Industry: indMatch[1]}
	if err := db.AddSymbol(symbol); err != nil {
		log.Fatal(symClean + " registration failed")
	} else {
		log.Println(symClean + " added")
	}
}
