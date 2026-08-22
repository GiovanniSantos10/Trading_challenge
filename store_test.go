package main

import (
	"errors"
	"testing"
	"time"
)

func TestMemoryStore_ConfirmWithinExpiry(t *testing.T) {
	s := NewMemoryStore()
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

func TestMemoryStore_ConfirmExpired(t *testing.T) {
	restore := timeNow
	defer func() { timeNow = restore }()

	base := time.Now()
	timeNow = func() time.Time { return base }

	s := NewMemoryStore()
	s.SavePending(Quote{ID: "1", Status: StatusPending, ExpiresAt: base.Add(10 * time.Second)})

	timeNow = func() time.Time { return base.Add(11 * time.Second) }

	_, err := s.Confirm("1")
	if !errors.Is(err, ErrQuoteExpired) {
		t.Fatalf("Confirm() error = %v, want %v", err, ErrQuoteExpired)
	}

	if len(s.ListConfirmed()) != 0 {
		t.Error("cotação expirada não deveria ir para confirmadas")
	}
}

func TestMemoryStore_ConfirmNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.Confirm("nao-existe")
	if !errors.Is(err, ErrQuoteNotFound) {
		t.Fatalf("Confirm() error = %v, want %v", err, ErrQuoteNotFound)
	}
}
