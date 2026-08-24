package handlers

import "testing"

func TestCalculateBRL(t *testing.T) {
	tests := []struct {
		name              string
		mxn               float64
		askBinanceUSDTBRL float64
		askOKXUSDTBRL     float64
		bidBinanceUSDTMXN float64
		want              float64
	}{
		{
			// Binance tem o menor ask -> deve ser usado no lugar do da OKX.
			name:              "usa o menor ask, Binance vence",
			mxn:               1000,
			askBinanceUSDTBRL: 5.50,
			askOKXUSDTBRL:     5.60,
			bidBinanceUSDTMXN: 18.20,
			// 1000/18.20 * 5.50 * 1.006 = 304.0109... -> teto = 304.02
			want: 304.02,
		},
		{
			// OKX tem o menor ask -> deve ser usado no lugar do da Binance.
			name:              "usa o menor ask, OKX vence",
			mxn:               1000,
			askBinanceUSDTBRL: 5.60,
			askOKXUSDTBRL:     5.50,
			bidBinanceUSDTMXN: 18.20,
			want:              304.02,
		},
		{
			// Resultado exato já bate em 2 casas -> teto não deve alterar o valor.
			name:              "resultado exato não sobe com o teto",
			mxn:               1,
			askBinanceUSDTBRL: 1,
			askOKXUSDTBRL:     1,
			bidBinanceUSDTMXN: 1,
			// 1 * 1 * 1.006 = 1.006 -> teto 2 casas = 1.01
			want: 1.01,
		},
		{
			// Fração mínima acima de um centavo já deve subir pro próximo centavo (teto).
			name:              "arredonda pra cima em fração de centavo",
			mxn:               1,
			askBinanceUSDTBRL: 1,
			askOKXUSDTBRL:     1,
			bidBinanceUSDTMXN: 1.006,
			// 1/1.006 * 1 * 1.006 = 1.0 exato -> sem sobra, fica 1.00
			want: 1.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateBRL(tt.mxn, tt.askBinanceUSDTBRL, tt.askOKXUSDTBRL, tt.bidBinanceUSDTMXN)
			if err != nil {
				t.Fatalf("CalculateBRL() error = %v, want nil", err)
			}
			if got != tt.want {
				t.Errorf("CalculateBRL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateBRL_InvalidBid(t *testing.T) {
	tests := []struct {
		name              string
		bidBinanceUSDTMXN float64
	}{
		{"bid zero", 0},
		{"bid negativo", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CalculateBRL(1000, 5.50, 5.60, tt.bidBinanceUSDTMXN)
			if err == nil {
				t.Fatalf("CalculateBRL() error = nil, want non-nil para bidBinanceUSDTMXN = %v", tt.bidBinanceUSDTMXN)
			}
		})
	}
}
