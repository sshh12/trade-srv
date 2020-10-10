package indexers

import (
	"time"

	events "github.com/sshh12/trade-srv/events"
)

type IndexerOptions struct {
	PollRate time.Duration
}

var EventIndexers = map[string]func(*events.EventStream, *IndexerOptions) error{
	"marketwatch":      startMarketWatchIndexer,
	"globalnewswire":   startGlobalNewsWireIndexer,
	"finviz":           startFinVizIndexer,
	"prnewswire":       startPrNewsWireIndexer,
	"barrons":          startBarronsIndexer,
	"biopharmcatalyst": startBioPharmCatalystIndexer,
	"benzinga":         startBenzingaIndexer,
	"reuters":          startReutersIndexer,
	"cnn":              startCNNIndexer,
	"earningscast":     startEarningsCastIndexer,
	"businesswire":     startBusinessWireSourceIndexer,
	"seekingalpha":     startSeekingAlphaIndexer,
}
