package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	MethodGetUpdates  = "getUpdates"
	MethodSendMessage = "sendMessage"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: basePathFromToken(token),
		client:   new(http.Client),
	}
}

func basePathFromToken(token string) string {
	return "bot" + token
}

// gets telegram user updates for the client
func (cl *Client) Update(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	resp, err := cl.sendRequest(MethodGetUpdates, q.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res UpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

// sends message to telegram user via the client
func (cl *Client) Message(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)
	_, err := cl.sendRequest(MethodSendMessage, q.Encode())
	if err != nil {
		return fmt.Errorf("unable to send message: %w", err)
	}
	return nil
}

// helper function that sends request with the client given method and query
func (cl *Client) sendRequest(method string, query string) (*http.Response, error) {
	apiUrl := url.URL{
		Scheme: "https",
		Host:   cl.host,
		Path:   path.Join(cl.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make new request: %w", err)
	}
	req.URL.RawQuery = query
	resp, err := cl.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("https request error: %w", err)
	}
	return resp, nil
}
