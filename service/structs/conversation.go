package structs

type Conversation struct {
	UserID        int    `json:"userID"`
	DestUserID    int    `json:"destUserID"`
	GlobalConvoID int    `json:"globalConvoID"`
	LastMsgID     int    `json:"lastMsgID,omitempty"`
	DelByUser     bool   `json:"delByUser"`
	Visible       bool   `json:"visible"`
	ConvoPropic64 string `json:"photo"`
}
