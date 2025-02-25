package structs

type Group struct {
	GroupName     string `json:"groupName"`
	GlobalConvoID int    `json:"groupID"`
	ConvoPropic64 string `json:"photo"`
	Members       []User `json:"members"`
}
