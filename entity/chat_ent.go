package entity

type Chat struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Msg       string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Message struct {
	Type string `json:"type"`
	User string `json:"user,omitempty"`
	Chat Chat   `json:"chat,omitempty"`
}
