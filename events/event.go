package events

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var hostname = ""

func init() {
	hostname, _ = os.Hostname()
}

type Event struct {
	ID               int
	Type             string   `pg:"type:'varchar'"`
	Source           string   `pg:"type:'varchar'"`
	Title            string   `pg:"type:'varchar'"`
	Content          string   `pg:"type:'text'"`
	URL              string   `pg:"type:'varchar'"`
	ConfirmedSymbols []string `pg:"type:'varchar',array"`
	Name             string   `pg:"type:'varchar'"`
	PrevValue        string   `pg:"type:'varchar'"`
	ExpectedValue    string   `pg:"type:'varchar'"`
	ActualValue      string   `pg:"type:'varchar'"`
	Impact           string   `pg:"type:'varchar'"`
	TimeFor          string   `pg:"type:'varchar'"`
	TimeReported     string   `pg:"type:'varchar'"`
	CacheHash        string   `pg:"type:'varchar',unique"`
	TimeLogged       string   `pg:"type:'timestamptz'"`
	HostName         string   `pg:"type:'varchar'"`
}

type EventStream struct {
	cacheLock        sync.RWMutex
	cache            map[string]bool
	db               *Database
	warmUpOver       int64
	warmUpInProgress bool
}

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

func (es *EventStream) OnEvent(evt *Event) {
	if evt.CacheHash == "" {
		panic(evt)
	}
	if es.hasCachedAddIfNot(evt) {
		return
	}
	if es.warmUpInProgress {
		if es.warmUpOver > time.Now().Unix() {
			log.Println("Discarded during warmup", evt.CacheHash)
			return
		}
		es.warmUpInProgress = false
	}
	evt.TimeLogged = time.Now().Format(time.RFC3339)
	evt.HostName = hostname
	fmt.Println(evt.TimeLogged)
	if err := es.db.AddEvent(evt); err != nil {
		log.Fatal(err)
	}
	fmt.Println("miss", evt.CacheHash)
}

func (es *EventStream) OnEventArticleResolveBody(source string, title string, url string, contentResolver func(string) string) {
	event := &Event{Type: "article", Source: source, Title: title, URL: url, CacheHash: HashKey(url)}
	if event.Content == "" && !es.hasCached(event) {
		event.Content = contentResolver(event.URL)
	}
	es.OnEvent(event)
}

func HashKey(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x\n", bs)
}
