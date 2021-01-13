package indexers

import (
	"time"

	events "github.com/sshh12/trade-srv/events"
)

// IndexerOptions are options that can be passed to indexers
type IndexerOptions struct {
	PollRate              time.Duration
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
	TwitterNames          []string
}

// EventIndexers is all the events indexers
var EventIndexers = map[string]func(*events.EventStream, *IndexerOptions) error{
	marketWatchSource:      startMarketWatchIndexer,
	globalNewsWireSource:   startGlobalNewsWireIndexer,
	finvizSource:           startFinVizIndexer,
	prNewsWireSource:       startPrNewsWireIndexer,
	barronsSource:          startBarronsIndexer,
	bioPharmCatalystSource: startBioPharmCatalystIndexer,
	benzingaSource:         startBenzingaIndexer,
	reutersSource:          startReutersIndexer,
	cnnSource:              startCNNIndexer,
	earningsCastSource:     startEarningsCastIndexer,
	businessWireSource:     startBusinessWireSourceIndexer,
	seekingAlphaSource:     startSeekingAlphaIndexer,
	nasdaqSource:           startNasdaqIndexer,
	investingComSource:     startInvestingComIndexer,
	twitterSource:          startTwitterIndexer,
	federalReserveSource:   startFederalReserveIndexer,
	cnbcSource:             startCNBCIndexer,
	blsTEDSource:           startBLSTEDIndexer,
	streetInsiderSource:    startStreetInsiderIndexer,
	businessInsiderSource:  startBusinessInsiderIndexer,
	forbesSource:           startForbesIndexer,
	zacksSource:            startZacksIndexer,
	washingtonPostSource:   startWashingtonPostIndexer,
}

// OtherIndexers are non-event indexers
var OtherIndexers = map[string]func(*events.EventStream, *IndexerOptions) error{}
