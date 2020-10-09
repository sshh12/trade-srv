package main

import (
	"flag"
	"log"
	"time"

	events "github.com/sshh12/trade-srv/events"
	indexers "github.com/sshh12/trade-srv/indexers"
)

func main() {
	pgUser := flag.String("pg_user", "postgres", "Postgres username")
	pgPassword := flag.String("pg_pass", "password", "Postgres password")
	pgAddr := flag.String("pg_addr", "localhost:5432", "Postgres host address")
	pgName := flag.String("pg_db", "tradesrv", "Postgres database name")
	indexersSelected := make(map[string]*bool)
	for name := range indexers.AllIndexers {
		indexersSelected[name] = flag.Bool("run_"+name, false, "Run "+name+" indexer")
	}
	runAll := flag.Bool("run_all", false, "Run all indexers")
	warmUp := flag.Int("warmup", 60, "Discard events that occur in this number of seconds")
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
	es := events.NewEventStream(db, time.Duration(*warmUp)*time.Second)
	for name, indexer := range indexers.AllIndexers {
		if *indexersSelected[name] {
			log.Println("Starting " + name)
			go indexer(es, &indexers.IndexerOptions{})
		}
	}
	for {
	}
}
