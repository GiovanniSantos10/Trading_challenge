package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Prices reúne as cotações necessárias pra CalculateBRL numa única chamada.
type Prices struct {
	AskBinanceUSDTBRL float64
	AskOKXUSDTBRL     float64
	BidBinanceUSDTMXN float64
}

// PriceProvider busca as cotações atuais nas exchanges. A implementação real
// (Binance + OKX) vive em binance.go.
type PriceProvider interface {
	GetPrices() (Prices, error)
}

type QuoteHandler struct {
	Store  QuoteStore
	Prices PriceProvider
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

type createQuoteRequest struct {
	MXN float64 `json:"mxn"`
}

func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var req createQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, `corpo inválido, esperado {"mxn": number}`)
		return
	}
	if req.MXN <= 0 {
		writeError(w, http.StatusBadRequest, "mxn deve ser maior que zero")
		return
	}

	prices, err := h.Prices.GetPrices()
	if err != nil {
		writeError(w, http.StatusBadGateway, "não foi possível obter cotações das exchanges")
		return
	}

	brl := CalculateBRL(req.MXN, prices.AskBinanceUSDTBRL, prices.AskOKXUSDTBRL, prices.BidBinanceUSDTMXN)

	now := timeNow()
	q := Quote{
		ID:        newID(),
		MXN:       req.MXN,
		BRL:       brl,
		CreatedAt: now,
		ExpiresAt: now.Add(10 * time.Second),
		Status:    StatusPending,
	}
	h.Store.SavePending(q)

	writeJSON(w, http.StatusCreated, q)
}

func (h *QuoteHandler) ConfirmQuote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	q, err := h.Store.Confirm(id)
	switch {
	case errors.Is(err, ErrQuoteNotFound):
		writeError(w, http.StatusNotFound, "cotação não encontrada")
		return
	case errors.Is(err, ErrQuoteExpired):
		writeError(w, http.StatusGone, "cotação expirada")
		return
	case err != nil:
		writeError(w, http.StatusInternalServerError, "erro ao confirmar cotação")
		return
	}

	writeJSON(w, http.StatusOK, q)
}

func (h *QuoteHandler) History(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Store.ListConfirmed())
}
