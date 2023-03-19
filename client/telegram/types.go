package telegram

import "net/http"

type Client struct {
	host     string
	basePath string
	client   *http.Client
}

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	Id      int      `json:"update_id"`
	Message *Message `json:"message"` //optional
}

type Message struct {
	Id   int    `json:"message_id"`
	Text string `json:"text"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type Chat struct {
	Id int `json:"id"`
}
