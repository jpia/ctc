package models

import (
	"ctc/logger"
	"sync"
	"time"
)

type Status string

const (
	PendingStatus  Status = "pending"
	DelayedStatus  Status = "delayed"
	ReleasedStatus Status = "released"
)

type ReleaseMethod string

const (
	OverrideReleaseMethod   ReleaseMethod = "override"
	StandardReleaseMethod   ReleaseMethod = "standard"
	ApiSickDayReleaseMethod ReleaseMethod = "api_sick_day"
)

type UpdatePriority int

const (
	HighUpdatePriority UpdatePriority = iota
	LowUpdatePriority
)

type URL struct {
	LongURL          string        `json:"long_url"`
	ReleaseDate      time.Time     `json:"release_date"`
	ReleaseDateOrig  time.Time     `json:"release_date_orig"`
	Shortcode        string        `json:"shortcode"`
	Status           Status        `json:"status"`
	Delays           int           `json:"delays"`
	ReleaseMethod    ReleaseMethod `json:"release_method"`
	ReleaseTimestamp time.Time     `json:"release_timestamp"`
}

type URLStore struct {
	store       map[string]URL
	mu          sync.RWMutex
	primaryCh   chan func(map[string]URL)
	secondaryCh chan func(map[string]URL)
}

var urlStoreInstance = &URLStore{
	store:       make(map[string]URL),
	primaryCh:   make(chan func(map[string]URL)),
	secondaryCh: make(chan func(map[string]URL)),
}

func GetURLStore() *URLStore {
	return urlStoreInstance
}

func (s *URLStore) run() {
	for {
		select {
		case update := <-s.primaryCh:
			s.mu.Lock()
			update(s.store)
			s.mu.Unlock()
		default:
			select {
			case update := <-s.primaryCh:
				s.mu.Lock()
				update(s.store)
				s.mu.Unlock()
			case update := <-s.secondaryCh:
				s.mu.Lock()
				update(s.store)
				s.mu.Unlock()
			}
		}
	}
}

func (s *URLStore) Set(shortcode string, url URL, priority UpdatePriority) {
	logger.DebugLog("Setting URL with shortcode %s, priority %d", shortcode, priority)
	update := func(store map[string]URL) {
		store[shortcode] = url
		logger.DebugLog("URL with shortcode %s set successfully", shortcode)
	}
	if priority == HighUpdatePriority {
		s.primaryCh <- update
	} else {
		s.secondaryCh <- update
	}
}

func (s *URLStore) Get(shortcode string) (URL, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, exists := s.store[shortcode]
	return url, exists
}

func (s *URLStore) Delete(shortcode string, priority UpdatePriority) {
	logger.DebugLog("Deleting URL with shortcode %s, priority %d", shortcode, priority)
	update := func(store map[string]URL) {
		delete(store, shortcode)
		logger.DebugLog("URL with shortcode %s deleted successfully", shortcode)
	}
	if priority == HighUpdatePriority {
		s.primaryCh <- update
	} else {
		s.secondaryCh <- update
	}
}

func (s *URLStore) GetAll() []URL {
	s.mu.RLock()
	defer s.mu.RUnlock()
	urls := make([]URL, 0, len(s.store))
	for _, url := range s.store {
		urls = append(urls, url)
	}
	return urls
}

func (s *URLStore) Reset() {
	// This function is only used for testing purposes
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store = make(map[string]URL)
	logger.DebugLog("URL store reset successfully")
}

func init() {
	go urlStoreInstance.run()
}
