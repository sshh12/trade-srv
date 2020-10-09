package main

import (
	"flag"

	events "github.com/sshh12/trade-srv/events"
	indexers "github.com/sshh12/trade-srv/indexers"
)

func main() {
	indexersSelected := make(map[string]*bool)
	for name := range indexers.AllIndexers {
		indexersSelected[name] = flag.Bool("run_"+name, false, "Run "+name+" indexer")
	}
	runAll := flag.Bool("run_all", false, "Run all indexers")
	flag.Parse()
	if *runAll {
		for name := range indexers.AllIndexers {
			indexersSelected[name] = runAll
		}
	}
	es := events.NewEventStream()
	for name, indexer := range indexers.AllIndexers {
		if *indexersSelected[name] {
			indexer(es, &indexers.IndexerOptions{})
		}
	}
}
