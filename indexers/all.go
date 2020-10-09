package indexers

import (
	"time"

	events "github.com/sshh12/trade-srv/events"
)

type IndexerOptions struct {
	PollRate time.Duration
}

var AllIndexers = map[string]func(*events.EventStream, *IndexerOptions) error{
	"marketwatch":    startMarketWatchIndexer,
	"globalnewswire": startGlobalNewsWireIndexer,
	"finviz":         startFinVizIndexer,
}
