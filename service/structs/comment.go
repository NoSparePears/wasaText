package structs

type Comment struct {
	CommID       int    `json:"commID"`
	MsgID        int    `json:"msgID"`
	SenderID     int    `json:"senderID"`
	SendUsername string `json:"sendUsername"`
	Emoji        string `json:"emoji"`
}
