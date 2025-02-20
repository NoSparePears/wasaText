package structs

type Message struct {
	MsgID     int    `json:"msgID"`
	ConvoID   int    `json:"convoID"`
	SenderID  int    `json:"senderID"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"` // Formato `RFC3339` (ISO 8601)
}
