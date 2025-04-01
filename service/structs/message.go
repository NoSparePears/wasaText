package structs

type Message struct {
	MsgID         int    `json:"msgID"`
	ConvoID       int    `json:"convoID"`
	SenderID      int    `json:"senderID"`
	Content       string `json:"content"`
	IsPhoto       int    `json:"isPhoto"`
	Timestamp     string `json:"timestamp"` // Formato `RFC3339` (ISO 8601)
	CheckSent     bool   `json:"sent"`
	CheckReceived bool   `json:"received"`
	IsForwarded   int    `json:"isForwarded"`
}
