package structs

type Group struct {
	GroupName     string `json:"groupName"`
	GlobalConvoID int    `json:"groupID"`
	GroupPropic64 string `json:"photo"`
	Members       []User `json:"members"`
	LastMsgID     int    `json:"lastMsgID,omitempty"`
}
