package handlers

import (
	"fmt"
	"math"
	"time"
)

type QuoteStatus string

const (
	StatusPending   QuoteStatus = "pending"
	StatusConfirmed QuoteStatus = "confirmed"
	StatusExpired   QuoteStatus = "expired"
)

type Quote struct {
	ID        string      `json:"id"`
	MXN       float64     `json:"mxn"`
	BRL       float64     `json:"brl"`
	CreatedAt time.Time   `json:"created_at"`
	ExpiresAt time.Time   `json:"expires_at"`
	Status    QuoteStatus `json:"status"`
}

func CalculateBRL(mxn, askBinanceUSDTBRL, askOKXUSDTBRL, bidBinanceUSDTMXN float64) (float64, error) {
	if bidBinanceUSDTMXN <= 0 {
		return 0, fmt.Errorf("bidBinanceUSDTMXN inválido: %v", bidBinanceUSDTMXN)
	}

	askUSDTBRL := min(askBinanceUSDTBRL, askOKXUSDTBRL)
	brl := mxn / bidBinanceUSDTMXN * askUSDTBRL * 1.006
	return math.Ceil(brl*100) / 100, nil
}
