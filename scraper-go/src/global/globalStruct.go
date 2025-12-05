package global

import (
	"net/http"
	"scraper/src/customTypes"
	"scraper/src/urlManager"
	"sync"
)

type ScraperSession struct {
	Output []customTypes.Page
	Queue  *urlManager.UrlQueue
	Seen   map[string]struct{}
	Mutex  sync.Mutex
	Wg     sync.WaitGroup
	Client *http.Client
}

func NewScraperSession() *ScraperSession {
	return &ScraperSession{
		Output: make([]customTypes.Page, 0),
		Seen:   make(map[string]struct{}),
		Mutex:  sync.Mutex{},
		Queue:  urlManager.NewUrlQueue(),
		Wg: sync.WaitGroup{},
	}
}

func (s *ScraperSession) ClearQueue() {
	s.Queue = nil
}
