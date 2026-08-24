package exchanges

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// withMockBinance sobe um httptest.Server e aponta binanceBaseURL pra ele
// durante o teste, restaurando no fim.
func withMockBinance(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	restore := binanceBaseURL
	binanceBaseURL = srv.URL
	t.Cleanup(func() { binanceBaseURL = restore })

	return srv
}

// Formato real da resposta da Binance: bidPrice/askPrice vêm como STRING.
const realBinanceResponse = `{
	"symbol": "BTCUSDT",
	"bidPrice": "61345.12000000",
	"bidQty": "0.50000000",
	"askPrice": "61346.00000000",
	"askQty": "0.30000000"
}`

func TestGetBookTicker_PricesAreParsedCorrectly(t *testing.T) {
	withMockBinance(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(realBinanceResponse))
	})

	got, err := GetBookTicker("BTCUSDT")
	if err != nil {
		t.Fatalf("GetBookTicker() error = %v, want nil", err)
	}

	if got.Symbol != "BTCUSDT" {
		t.Errorf("Symbol = %q, want %q", got.Symbol, "BTCUSDT")
	}
	const wantBid = float32(61345.12)
	if got.BidPrice != wantBid {
		t.Errorf("BidPrice = %v, want %v", got.BidPrice, wantBid)
	}
	const wantAsk = float32(61346.00)
	if got.AskPrice != wantAsk {
		t.Errorf("AskPrice = %v, want %v", got.AskPrice, wantAsk)
	}
}

func TestGetBookTicker_ReturnsErrorOnRequestFailure(t *testing.T) {
	restore := binanceBaseURL
	// porta sem nada escutando -> http.Get retorna erro de conexão, sem precisar de rede real
	binanceBaseURL = "http://127.0.0.1:1"
	defer func() { binanceBaseURL = restore }()

	_, err := GetBookTicker("BTCUSDT")
	if err == nil {
		t.Fatal("GetBookTicker() error = nil, want non-nil em falha de conexão")
	}
}

func TestGetBookTicker_ReturnsErrorOnHTTPErrorStatus(t *testing.T) {
	withMockBinance(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code":-1121,"msg":"Invalid symbol."}`))
	})

	_, err := GetBookTicker("SIMBOLO_INVALIDO")
	if err == nil {
		t.Fatal("GetBookTicker() error = nil, want non-nil em status 400")
	}
}

func TestGetBookTicker_ReturnsErrorOnMalformedPrice(t *testing.T) {
	withMockBinance(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"symbol":"BTCUSDT","bidPrice":"not-a-number","bidQty":"0","askPrice":"0","askQty":"0"}`))
	})

	_, err := GetBookTicker("BTCUSDT")
	if err == nil {
		t.Fatal("GetBookTicker() error = nil, want non-nil quando bidPrice não é numérico")
	}
}
