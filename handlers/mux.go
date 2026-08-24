package handlers

import "net/http"

// NewMux monta as rotas HTTP da API sobre o QuoteHandler dado.
func NewMux(qh *QuoteHandler) *http.ServeMux {
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
