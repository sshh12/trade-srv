package main

import (
	events "github.com/sshh12/trade-srv/events"
	indexers "github.com/sshh12/trade-srv/indexers"
)

func main() {
	es := events.NewEventStream()
	for _, indexer := range indexers.AllIndexers {
		opts := &indexers.IndexerOptions{}
		indexer(es, opts)
	}
}
