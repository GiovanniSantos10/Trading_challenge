package handlers

import "net/http"

func NewMux(qh *QuoteHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /quotes", qh.CreateQuote)
	mux.HandleFunc("POST /quotes/{id}/confirm", qh.ConfirmQuote)
	mux.HandleFunc("GET /history", qh.History)
	return withCORS(mux)
}
