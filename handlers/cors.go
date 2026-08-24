package handlers

import "net/http"

// withCORS libera chamadas de qualquer origem. A API não usa cookies/sessão
// (sem credentials), então um Allow-Origin permissivo não abre brecha de
// segurança real — só evita ter que sincronizar a porta do front (Vite,
// CRA, etc.) com a API depois.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
