package exchanges

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var okxBaseURL = "https://www.okx.com"

type Data_okx struct {
	Symbol   string
	BidPrice float32
	BidQty   float32
	AskPrice float32
	AskQty   float32
}

func (d *Data_okx) UnmarshalJSON(data []byte) error {
	var raw struct {
		InstId string `json:"instId"`
		BidPx  string `json:"bidPx"`
		BidSz  string `json:"bidSz"`
		AskPx  string `json:"askPx"`
		AskSz  string `json:"askSz"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	bidPrice, err := strconv.ParseFloat(raw.BidPx, 32)
	if err != nil {
		return fmt.Errorf("parsing bidPx %q: %w", raw.BidPx, err)
	}
	bidQty, err := strconv.ParseFloat(raw.BidSz, 32)
	if err != nil {
		return fmt.Errorf("parsing bidSz %q: %w", raw.BidSz, err)
	}
	askPrice, err := strconv.ParseFloat(raw.AskPx, 32)
	if err != nil {
		return fmt.Errorf("parsing askPx %q: %w", raw.AskPx, err)
	}
	askQty, err := strconv.ParseFloat(raw.AskSz, 32)
	if err != nil {
		return fmt.Errorf("parsing askSz %q: %w", raw.AskSz, err)
	}

	d.Symbol = raw.InstId
	d.BidPrice = float32(bidPrice)
	d.BidQty = float32(bidQty)
	d.AskPrice = float32(askPrice)
	d.AskQty = float32(askQty)
	return nil
}

type okxTickerResponse struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data []Data_okx `json:"data"`
}

func Get_mxn_okx(token string) (Data_okx, error) {
	url := fmt.Sprintf("%s/api/v5/market/ticker?instId=%s", okxBaseURL, token)

	resp, err := http.Get(url)
	if err != nil {
		return Data_okx{}, fmt.Errorf("okx: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Data_okx{}, fmt.Errorf("okx: reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Data_okx{}, fmt.Errorf("okx: unexpected status %d: %s", resp.StatusCode, body)
	}

	var parsed okxTickerResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return Data_okx{}, fmt.Errorf("okx: decoding response: %w", err)
	}

	if parsed.Code != "0" {
		return Data_okx{}, fmt.Errorf("okx: api error code %s: %s", parsed.Code, parsed.Msg)
	}
	if len(parsed.Data) != 1 {
		return Data_okx{}, fmt.Errorf("okx: expected 1 ticker in data, got %d", len(parsed.Data))
	}

	return parsed.Data[0], nil
}
