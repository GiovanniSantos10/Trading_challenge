package exchanges

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var binanceBaseURL = "https://api.binance.com"

type Data_binance struct {
	Symbol   string
	BidPrice float32
	BidQty   float32
	AskPrice float32
	AskQty   float32
}

func (d *Data_binance) UnmarshalJSON(data []byte) error {
	var raw struct {
		Symbol   string `json:"symbol"`
		BidPrice string `json:"bidPrice"`
		BidQty   string `json:"bidQty"`
		AskPrice string `json:"askPrice"`
		AskQty   string `json:"askQty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	bidPrice, err := strconv.ParseFloat(raw.BidPrice, 32)
	if err != nil {
		return fmt.Errorf("parsing bidPrice %q: %w", raw.BidPrice, err)
	}
	bidQty, err := strconv.ParseFloat(raw.BidQty, 32)
	if err != nil {
		return fmt.Errorf("parsing bidQty %q: %w", raw.BidQty, err)
	}
	askPrice, err := strconv.ParseFloat(raw.AskPrice, 32)
	if err != nil {
		return fmt.Errorf("parsing askPrice %q: %w", raw.AskPrice, err)
	}
	askQty, err := strconv.ParseFloat(raw.AskQty, 32)
	if err != nil {
		return fmt.Errorf("parsing askQty %q: %w", raw.AskQty, err)
	}

	d.Symbol = raw.Symbol
	d.BidPrice = float32(bidPrice)
	d.BidQty = float32(bidQty)
	d.AskPrice = float32(askPrice)
	d.AskQty = float32(askQty)
	return nil
}

func GetBookTicker(token string) (Data_binance, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/bookTicker?symbol=%s", binanceBaseURL, token)

	resp, err := http.Get(url)
	if err != nil {
		return Data_binance{}, fmt.Errorf("binance: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Data_binance{}, fmt.Errorf("binance: reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Data_binance{}, fmt.Errorf("binance: unexpected status %d: %s", resp.StatusCode, body)
	}

	var TokenInfo Data_binance
	if err := json.Unmarshal(body, &TokenInfo); err != nil {
		return Data_binance{}, fmt.Errorf("binance: decoding response: %w", err)
	}

	return TokenInfo, nil
}
