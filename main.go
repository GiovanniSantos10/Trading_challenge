package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func newMux(qh *QuoteHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /quotes", qh.CreateQuote)
	mux.HandleFunc("POST /quotes/{id}/confirm", qh.ConfirmQuote)
	mux.HandleFunc("GET /history", qh.History)
	return mux
}

// notImplementedPrices é um placeholder até o PriceProvider real (Binance +
// OKX, em binance.go) estar pronto.
//
// TODO(você): trocar por um PriceProvider real.
type notImplementedPrices struct{}

func (notImplementedPrices) GetPrices() (Prices, error) {
	return Prices{}, fmt.Errorf("price provider not implemented")
}

func main() {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	} else {
		addr = ":" + addr
	}

	qh := &QuoteHandler{
		Store:  NewMemoryStore(),
		Prices: notImplementedPrices{},
	}

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, newMux(qh)); err != nil {
		log.Fatal(err)
	}
}
