package main

import "sync"

type Store struct {
	SETs   map[string]string
	SETsMu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		SETs: map[string]string{},
	}
}

var store = NewStore()

func (s *Store) Set(key, value string) {
	s.SETsMu.Lock()
	defer s.SETsMu.Unlock()

	s.SETs[key] = value

	return
}

func (s *Store) Get(key string) (string, bool) {
	s.SETsMu.RLock()
	defer s.SETsMu.RUnlock()

	value, ok := s.SETs[key]

	return value, ok
}
