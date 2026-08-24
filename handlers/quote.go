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

// CalculateBRL converts an MXN amount into BRL via USDT.
//
// Fórmula: mxn × ask(USDT/BRL, o menor entre Binance e OKX) × 1/bid(USDT/MXN, Binance) × 1.006
//
// Regras:
//   - askBinanceUSDTBRL e askOKXUSDTBRL são o preço de venda de USDT em BRL em cada exchange
//     (BRL por USDT); use o MENOR dos dois.
//   - bidBinanceUSDTMXN é o preço de compra de USDT em MXN na Binance (MXN por USDT); o mxn
//     é convertido para USDT dividindo por esse valor.
//   - 1.006 é a margem fixa aplicada sobre o resultado.
//   - O arredondamento para cima (teto), com 2 casas decimais, acontece UMA ÚNICA VEZ, no
//     final de todo o cálculo — não arredonde valores intermediários.
//
// Retorna erro se bidBinanceUSDTMXN <= 0: um bid inválido (0 ou negativo) geraria divisão
// por zero (+Inf) ou resultado sem sentido (BRL negativo), então é rejeitado explicitamente
// em vez de propagar um valor inválido pro cliente.
func CalculateBRL(mxn, askBinanceUSDTBRL, askOKXUSDTBRL, bidBinanceUSDTMXN float64) (float64, error) {
	if bidBinanceUSDTMXN <= 0 {
		return 0, fmt.Errorf("bidBinanceUSDTMXN inválido: %v", bidBinanceUSDTMXN)
	}

	askUSDTBRL := min(askBinanceUSDTBRL, askOKXUSDTBRL)
	brl := mxn / bidBinanceUSDTMXN * askUSDTBRL * 1.006
	return math.Ceil(brl*100) / 100, nil
}
