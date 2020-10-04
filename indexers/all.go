package indexers

import (
	events "github.com/sshh12/trade-srv/events"
	"time"
)

type IndexerOptions struct {
	PollRate time.Duration
}

var AllIndexers = map[string]func(*events.EventStream, *IndexerOptions) error{
	"marketwatch": startMarketWatchIndexer,
}
