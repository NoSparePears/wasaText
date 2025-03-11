package structs

import "database/sql"

type Conversation struct {
	UserConvoID   int           `json:"userConvoID"`
	UserID        int           `json:"userID"`
	DestUserID    int           `json:"destUserID"`
	GlobalConvoID int           `json:"globalConvoID"`
	LastMsgID     sql.NullInt64 `json:"lastMsgID,omitempty"`
	DelByUser     bool          `json:"delByUser"`
	Visible       bool          `json:"visible"`
	ConvoPropic64 string        `json:"photo"`
}
