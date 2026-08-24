package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"KX_trading_challenge/exchanges"
	"KX_trading_challenge/handlers"
)

// LivePriceProvider implementa handlers.PriceProvider buscando as cotações
// reais na Binance e na OKX. Busca exatamente os 3 tickers que a fórmula
// (handlers.CalculateBRL) usa — o "menor ask" entre as exchanges é decidido
// lá dentro, não aqui.
type LivePriceProvider struct{}

func (LivePriceProvider) GetPrices() (handlers.Prices, error) {
	binanceUSDTBRL, err := exchanges.GetBookTicker("USDTBRL")
	if err != nil {
		return handlers.Prices{}, fmt.Errorf("binance USDT/BRL: %w", err)
	}

	okxUSDTBRL, err := exchanges.Get_mxn_okx("USDT-BRL")
	if err != nil {
		return handlers.Prices{}, fmt.Errorf("okx USDT/BRL: %w", err)
	}

	binanceUSDTMXN, err := exchanges.GetBookTicker("USDTMXN")
	if err != nil {
		return handlers.Prices{}, fmt.Errorf("binance USDT/MXN: %w", err)
	}

	return handlers.Prices{
		AskBinanceUSDTBRL: float64(binanceUSDTBRL.AskPrice),
		AskOKXUSDTBRL:     float64(okxUSDTBRL.AskPrice),
		BidBinanceUSDTMXN: float64(binanceUSDTMXN.BidPrice),
	}, nil
}

func main() {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	} else {
		addr = ":" + addr
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "quotes.db"
	}
	store, err := handlers.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("abrindo banco em %q: %v", dbPath, err)
	}
	defer store.Close()

	qh := &handlers.QuoteHandler{
		Store:  store,
		Prices: LivePriceProvider{},
	}

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handlers.NewMux(qh)); err != nil {
		log.Fatal(err)
	}
}
