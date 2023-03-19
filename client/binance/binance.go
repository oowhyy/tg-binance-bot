package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

const apiv3Path = "/api/v3"
const tickerPath = "/ticker/bookTicker"
const exchangeInfoPath = "/exchangeInfo"

func New(host string) *Client {
	return &Client{
		host:     host,
		basePath: apiv3Path,
		client:   new(http.Client),
	}
}

// send request - GET /api/v3/ticker/bookTicker
func (cl *Client) Ticker() ([]BookTicker, error) {
	apiUrl := url.URL{
		Scheme: "https",
		Host:   cl.host,
		Path:   path.Join(cl.basePath, tickerPath),
	}
	req, err := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make new request: %w", err)
	}
	resp, err := cl.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("https request error: %w", err)
	}
	defer resp.Body.Close()
	var res []BookTicker
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}

// send request - GET /api/v3/exchangeInfo
func (cl *Client) ExchangeInfo() (*ExchangeInfo, error) {
	apiUrl := url.URL{
		Scheme: "https",
		Host:   cl.host,
		Path:   path.Join(cl.basePath, exchangeInfoPath),
	}
	req, err := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make new request: %w", err)
	}
	// req.URL.RawQuery = query
	resp, err := cl.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("https request error: %w", err)
	}
	defer resp.Body.Close()
	var res ExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}
