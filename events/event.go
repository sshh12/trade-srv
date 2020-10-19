package events

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var hostname = ""

func init() {
	hostname, _ = os.Hostname()
}

// Event represents an events that occured in the world
type Event struct {
	ID               int
	Type             string                 `pg:"type:'varchar'"`
	Source           string                 `pg:"type:'varchar'"`
	Title            string                 `pg:"type:'varchar'"`
	Content          string                 `pg:"type:'text'"`
	URL              string                 `pg:"type:'varchar'"`
	ConfirmedSymbols []string               `pg:"type:'varchar',array"`
	Name             string                 `pg:"type:'varchar'"`
	Author           string                 `pg:"type:'varchar'"`
	PrevValue        string                 `pg:"type:'varchar'"`
	ExpectedValue    string                 `pg:"type:'varchar'"`
	ActualValue      string                 `pg:"type:'varchar'"`
	Impact           string                 `pg:"type:'varchar'"`
	TimeFor          string                 `pg:"type:'varchar'"`
	TimeReported     string                 `pg:"type:'varchar'"`
	Extras           map[string]interface{} `pg:"type:'json'"`
	CacheHash        string                 `pg:"type:'varchar',unique"`
	TimeLogged       string                 `pg:"type:'timestamptz'"`
	HostName         string                 `pg:"type:'varchar'"`
}

// TDAOHLCV open, high, low, close, volume
type TDAOHLCV struct {
	ID     string `pg:"type:'varchar',pk"`
	Date   int
	Symbol string `pg:"type:'varchar'"`
	Open   float64
	Close  float64
	Low    float64
	High   float64
	Volume float64
}

// GuruFin for fin data
type GuruFin struct {
	ID     string `pg:"type:'varchar',pk"`
	Date   int
	Symbol string                 `pg:"type:'varchar'"`
	Fin    map[string]interface{} `pg:"type:'json'"`
}

// EventStream is a stream of events
type EventStream struct {
	cacheLock        sync.RWMutex
	cache            map[string]bool
	db               *Database
	warmUpOver       int64
	warmUpInProgress bool
}

// NewEventStream creates an event stream
func NewEventStream(db *Database, warmUp time.Duration) *EventStream {
	return &EventStream{
		cacheLock:        sync.RWMutex{},
		cache:            make(map[string]bool),
		db:               db,
		warmUpInProgress: true,
		warmUpOver:       time.Now().Add(warmUp).Unix(),
	}
}

func (es *EventStream) hasCached(evt *Event) bool {
	es.cacheLock.Lock()
	defer es.cacheLock.Unlock()
	return es.cache[evt.CacheHash]
}

func (es *EventStream) hasCachedAddIfNot(evt *Event) bool {
	es.cacheLock.Lock()
	defer es.cacheLock.Unlock()
	val := es.cache[evt.CacheHash]
	if !val {
		es.cache[evt.CacheHash] = true
	}
	return val
}

// GetSymbols gets db of symbols
func (es *EventStream) GetSymbols() ([]Symbol, error) {
	return es.db.GetSymbols()
}

// OnEvent handles the occurance of an event
func (es *EventStream) OnEvent(evt *Event) {
	if evt.CacheHash == "" {
		panic(evt)
	}
	if es.hasCachedAddIfNot(evt) {
		return
	}
	if es.warmUpInProgress {
		if es.warmUpOver > time.Now().Unix() {
			log.WithField("hash", evt.CacheHash).WithField("reason", "warmup").Debug("Event Discarded")
			return
		}
		es.warmUpInProgress = false
	}
	evt.TimeLogged = time.Now().Format(time.RFC3339)
	evt.HostName = hostname
	if err := es.db.AddEvent(evt); err != nil {
		log.Print(err)
	}
	log.WithField("hash", evt.CacheHash).Debug("Event Mined")
}

// OnEventArticleResolveBody handles an article event, if appropriate it calls contentResolver
// for the content of the article
func (es *EventStream) OnEventArticleResolveBody(source string, title string, url string, contentResolver func(string) string) {
	event := &Event{Type: "article", Source: source, Title: title, URL: url, CacheHash: HashKey(url)}
	if event.Content == "" && !es.hasCached(event) && !es.warmUpInProgress {
		event.Content = contentResolver(event.URL)
	}
	es.OnEvent(event)
}

// OnMinOHLCVs handles several OHLCVs
func (es *EventStream) OnMinOHLCVs(ticks []TDAOHLCV) error {
	return es.db.AddMinOHLCVs(ticks)
}

// HashKey hashes the input string
func HashKey(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
