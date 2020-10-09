package events

import (
	"crypto/sha256"
	"fmt"
	"log"
	"sync"
	"time"
)

type Event struct {
	Source    string
	Title     string
	Content   string
	URL       string
	CacheHash string
}

type EventStream struct {
	cacheLock        sync.RWMutex
	cache            map[string]bool
	warmUpOver       int64
	warmUpInProgress bool
}

func NewEventStream(warmUp time.Duration) *EventStream {
	return &EventStream{
		cacheLock:        sync.RWMutex{},
		cache:            make(map[string]bool),
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
	fmt.Println("miss", evt.CacheHash)
}

func (es *EventStream) OnEventArticleResolveBody(source string, title string, url string, contentResolver func(string) string) {
	event := &Event{Source: source, Title: title, URL: url, CacheHash: hashKey(url)}
	if event.Content == "" && !es.hasCached(event) {
		event.Content = contentResolver(event.URL)
	}
	es.OnEvent(event)
}

func hashKey(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x\n", bs)
}
