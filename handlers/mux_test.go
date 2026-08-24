package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMux_CORSPreflight(t *testing.T) {
	h := &QuoteHandler{Store: NewMemoryStore()}
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodOptions, "/quotes", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "*")
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("Access-Control-Allow-Methods não foi setado")
	}
}

func TestNewMux_CORSHeaderOnNormalRequest(t *testing.T) {
	h := &QuoteHandler{Store: NewMemoryStore()}
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "*")
	}
}
