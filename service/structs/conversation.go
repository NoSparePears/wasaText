package structs

type Conversation struct {
	UserConvoID   int    `json:"userConvoID"`
	UserID        int    `json:"userID"`
	GlobalConvoID int    `json:"globalConvoID"`
	LastMsgID     int    `json:"lastMsgID,omitempty"`
	DelByUser     bool   `json:"delByUser"`
	Visible       bool   `json:"visible"`
	ConvoPropic64 string `json:"photo"`
}
