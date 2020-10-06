package events

import "fmt"

type Event struct {
	Source    string
	Title     string
	Content   string
	URL       string
	CacheKey  string
	CacheHash string
}

type EventStream struct{}

func NewEventStream() *EventStream {
	return &EventStream{}
}

func (es *EventStream) OnEvent(evt *Event) {
	fmt.Println(evt)
}

func (es *EventStream) OnEventArticle(source string, title string, url string, content string) {
	es.OnEvent(&Event{Source: source, Title: title, URL: url, Content: content, CacheKey: url})
}

func (es *EventStream) OnEventArticleResolveBody(source string, title string, url string, contentResolver func(string) string) {
	event := &Event{Source: source, Title: title, URL: url, CacheKey: url}
	if event.Content == "" {
		event.Content = contentResolver(event.URL)
	}
	es.OnEvent(event)
}
