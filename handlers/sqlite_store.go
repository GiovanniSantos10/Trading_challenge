package handlers

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB

	mu      sync.Mutex
	pending map[string]Quote
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("sqlite: abrindo %q: %w", path, err)
	}

	const schema = `
	CREATE TABLE IF NOT EXISTS quotes (
		id         TEXT PRIMARY KEY,
		mxn        REAL NOT NULL,
		brl        REAL NOT NULL,
		created_at TEXT NOT NULL,
		expires_at TEXT NOT NULL,
		status     TEXT NOT NULL
	);`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite: criando tabela quotes: %w", err)
	}

	return &SQLiteStore{
		db:      db,
		pending: make(map[string]Quote),
	}, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) SavePending(q Quote) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pending[q.ID] = q
}

func (s *SQLiteStore) GetPending(id string) (Quote, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	q, ok := s.pending[id]
	return q, ok
}

func (s *SQLiteStore) Confirm(id string) (Quote, error) {
	s.mu.Lock()
	q, ok := s.pending[id]
	if !ok {
		s.mu.Unlock()
		return Quote{}, ErrQuoteNotFound
	}
	if timeNow().After(q.ExpiresAt) {
		delete(s.pending, id)
		s.mu.Unlock()
		return Quote{}, ErrQuoteExpired
	}
	q.Status = StatusConfirmed
	delete(s.pending, id)
	s.mu.Unlock()

	const insert = `
	INSERT INTO quotes (id, mxn, brl, created_at, expires_at, status)
	VALUES (?, ?, ?, ?, ?, ?);`
	_, err := s.db.Exec(insert, q.ID, q.MXN, q.BRL,
		q.CreatedAt.Format(time.RFC3339Nano), q.ExpiresAt.Format(time.RFC3339Nano), q.Status)
	if err != nil {
		return Quote{}, fmt.Errorf("sqlite: gravando cotação confirmada: %w", err)
	}

	return q, nil
}

func (s *SQLiteStore) ListConfirmed() []Quote {
	const query = `
	SELECT id, mxn, brl, created_at, expires_at, status
	FROM quotes
	ORDER BY created_at;`
	rows, err := s.db.Query(query)
	if err != nil {
		return []Quote{}
	}
	defer rows.Close()

	out := []Quote{}
	for rows.Next() {
		var q Quote
		var createdAt, expiresAt string
		if err := rows.Scan(&q.ID, &q.MXN, &q.BRL, &createdAt, &expiresAt, &q.Status); err != nil {
			return []Quote{}
		}
		q.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
		q.ExpiresAt, _ = time.Parse(time.RFC3339Nano, expiresAt)
		out = append(out, q)
	}
	return out
}
