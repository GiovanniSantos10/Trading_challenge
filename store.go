package main

import (
	"fmt"
	"sync"
	"time"
)

// timeNow é substituível nos testes, pra simular expiração sem sleep real.
var timeNow = time.Now

// QuoteStore guarda cotações pendentes (em memória) e confirmadas.
// A implementação de persistência (SQLite) vai satisfazer a mesma interface,
// então os handlers HTTP não precisam mudar quando ela entrar.
type QuoteStore interface {
	SavePending(q Quote)
	GetPending(id string) (Quote, bool)
	// Confirm move a cotação de pending para confirmed, se ela existir e
	// ainda não tiver expirado. Retorna a cotação confirmada.
	Confirm(id string) (Quote, error)
	ListConfirmed() []Quote
}

var (
	ErrQuoteNotFound = fmt.Errorf("quote not found")
	ErrQuoteExpired  = fmt.Errorf("quote expired")
)

// MemoryStore é a implementação em memória de QuoteStore, usada até o
// passo da persistência SQLite entrar.
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
