package structs

type Comment struct {
	CommID    int    `json:"commID"`
	MsgID     int    `json:"msgID"`
	SenderID  int    `json:"senderID"`
	Emoji     string `json:"emoji"`
	Timestamp string `json:"timestamp"`
}
