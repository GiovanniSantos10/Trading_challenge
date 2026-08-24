package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateQuote_InvalidBody(t *testing.T) {
	h := &QuoteHandler{Store: NewMemoryStore()}
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBufferString(`not json`))
	rec := httptest.NewRecorder()

	h.CreateQuote(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestCreateQuote_NonPositiveMXN(t *testing.T) {
	h := &QuoteHandler{Store: NewMemoryStore()}
	body, _ := json.Marshal(map[string]float64{"mxn": 0})
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	h.CreateQuote(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestConfirmQuote_NotFound(t *testing.T) {
	h := &QuoteHandler{Store: NewMemoryStore()}
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodPost, "/quotes/nao-existe/confirm", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestConfirmQuote_Expired(t *testing.T) {
	restore := timeNow
	defer func() { timeNow = restore }()
	base := time.Now()
	timeNow = func() time.Time { return base }

	store := NewMemoryStore()
	store.SavePending(Quote{ID: "1", Status: StatusPending, ExpiresAt: base.Add(10 * time.Second)})
	h := &QuoteHandler{Store: store}
	mux := NewMux(h)

	timeNow = func() time.Time { return base.Add(11 * time.Second) }

	req := httptest.NewRequest(http.MethodPost, "/quotes/1/confirm", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusGone {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusGone)
	}
}

func TestConfirmQuote_Success(t *testing.T) {
	store := NewMemoryStore()
	store.SavePending(Quote{ID: "1", MXN: 100, BRL: 27.95, Status: StatusPending, ExpiresAt: timeNow().Add(10 * time.Second)})
	h := &QuoteHandler{Store: store}
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodPost, "/quotes/1/confirm", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got Quote
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Status != StatusConfirmed {
		t.Errorf("Status = %v, want %v", got.Status, StatusConfirmed)
	}
}

func TestHistory_ListsConfirmed(t *testing.T) {
	store := NewMemoryStore()
	store.SavePending(Quote{ID: "1", Status: StatusPending, ExpiresAt: timeNow().Add(10 * time.Second)})
	if _, err := store.Confirm("1"); err != nil {
		t.Fatalf("setup Confirm() error = %v", err)
	}
	h := &QuoteHandler{Store: store}
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got []Quote
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(got) != 1 || got[0].ID != "1" {
		t.Errorf("history = %v, want 1 item com ID 1", got)
	}
}
