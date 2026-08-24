package handlers

import (
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func TestSQLiteStore_ConfirmWithinExpiry(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	s, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("NewSQLiteStore() error = %v", err)
	}
	defer s.Close()

	q := Quote{ID: "1", MXN: 100, BRL: 27.95, Status: StatusPending, ExpiresAt: timeNow().Add(10 * time.Second)}
	s.SavePending(q)

	got, err := s.Confirm("1")
	if err != nil {
		t.Fatalf("Confirm() error = %v, want nil", err)
	}
	if got.Status != StatusConfirmed {
		t.Errorf("Status = %v, want %v", got.Status, StatusConfirmed)
	}

	if _, ok := s.GetPending("1"); ok {
		t.Error("cotação confirmada ainda aparece como pending")
	}

	confirmed := s.ListConfirmed()
	if len(confirmed) != 1 || confirmed[0].ID != "1" {
		t.Errorf("ListConfirmed() = %v, want 1 item com ID 1", confirmed)
	}
}

func TestSQLiteStore_ConfirmExpired(t *testing.T) {
	restore := timeNow
	defer func() { timeNow = restore }()

	base := time.Now()
	timeNow = func() time.Time { return base }

	path := filepath.Join(t.TempDir(), "test.db")
	s, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("NewSQLiteStore() error = %v", err)
	}
	defer s.Close()

	s.SavePending(Quote{ID: "1", Status: StatusPending, ExpiresAt: base.Add(10 * time.Second)})

	timeNow = func() time.Time { return base.Add(11 * time.Second) }

	_, err = s.Confirm("1")
	if !errors.Is(err, ErrQuoteExpired) {
		t.Fatalf("Confirm() error = %v, want %v", err, ErrQuoteExpired)
	}

	if len(s.ListConfirmed()) != 0 {
		t.Error("cotação expirada não deveria ir para confirmadas")
	}
}

func TestSQLiteStore_ConfirmNotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	s, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("NewSQLiteStore() error = %v", err)
	}
	defer s.Close()

	_, err = s.Confirm("nao-existe")
	if !errors.Is(err, ErrQuoteNotFound) {
		t.Fatalf("Confirm() error = %v, want %v", err, ErrQuoteNotFound)
	}
}

func TestSQLiteStore_PersistsAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	s1, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("NewSQLiteStore() error = %v", err)
	}
	s1.SavePending(Quote{ID: "1", MXN: 50, BRL: 13.98, Status: StatusPending, ExpiresAt: timeNow().Add(10 * time.Second)})
	if _, err := s1.Confirm("1"); err != nil {
		t.Fatalf("Confirm() error = %v, want nil", err)
	}
	if err := s1.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	s2, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("NewSQLiteStore() (reopen) error = %v", err)
	}
	defer s2.Close()

	confirmed := s2.ListConfirmed()
	if len(confirmed) != 1 || confirmed[0].ID != "1" || confirmed[0].BRL != 13.98 {
		t.Errorf("ListConfirmed() após reabrir = %v, want 1 item com ID 1 e BRL 13.98", confirmed)
	}
}
