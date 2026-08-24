package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"KX_trading_challenge/exchanges"
	"KX_trading_challenge/handlers"
)

type LivePriceProvider struct{}

func (LivePriceProvider) GetPrices() (handlers.Prices, error) {
	var (
		wg sync.WaitGroup

		binanceUSDTBRL exchanges.Data_binance
		errBinanceBRL  error

		okxUSDTBRL exchanges.Data_okx
		errOKXBRL  error

		binanceUSDTMXN exchanges.Data_binance
		errBinanceMXN  error
	)

	wg.Add(3)
	go func() {
		defer wg.Done()
		binanceUSDTBRL, errBinanceBRL = exchanges.GetBookTicker("USDTBRL")
	}()
	go func() {
		defer wg.Done()
		okxUSDTBRL, errOKXBRL = exchanges.Get_mxn_okx("USDT-BRL")
	}()
	go func() {
		defer wg.Done()
		binanceUSDTMXN, errBinanceMXN = exchanges.GetBookTicker("USDTMXN")
	}()
	wg.Wait()

	if errBinanceBRL != nil {
		return handlers.Prices{}, fmt.Errorf("binance USDT/BRL: %w", errBinanceBRL)
	}
	if errOKXBRL != nil {
		return handlers.Prices{}, fmt.Errorf("okx USDT/BRL: %w", errOKXBRL)
	}
	if errBinanceMXN != nil {
		return handlers.Prices{}, fmt.Errorf("binance USDT/MXN: %w", errBinanceMXN)
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
