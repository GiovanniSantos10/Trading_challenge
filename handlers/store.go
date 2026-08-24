package handlers

import (
	"fmt"
	"sync"
	"time"
)

var timeNow = time.Now

type QuoteStore interface {
	SavePending(q Quote)
	GetPending(id string) (Quote, bool)
	Confirm(id string) (Quote, error)
	ListConfirmed() []Quote
}

var (
	ErrQuoteNotFound = fmt.Errorf("quote not found")
	ErrQuoteExpired  = fmt.Errorf("quote expired")
)

type MemoryStore struct {
	mu        sync.Mutex
	pending   map[string]Quote
	confirmed []Quote
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		pending: make(map[string]Quote),
	}
}

func (s *MemoryStore) SavePending(q Quote) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pending[q.ID] = q
}

func (s *MemoryStore) GetPending(id string) (Quote, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	q, ok := s.pending[id]
	return q, ok
}

func (s *MemoryStore) Confirm(id string) (Quote, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	q, ok := s.pending[id]
	if !ok {
		return Quote{}, ErrQuoteNotFound
	}
	if timeNow().After(q.ExpiresAt) {
		delete(s.pending, id)
		return Quote{}, ErrQuoteExpired
	}

	q.Status = StatusConfirmed
	delete(s.pending, id)
	s.confirmed = append(s.confirmed, q)
	return q, nil
}

func (s *MemoryStore) ListConfirmed() []Quote {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Quote, len(s.confirmed))
	copy(out, s.confirmed)
	return out
}
